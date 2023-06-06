package aggregator

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/metrics"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

type GenerateProof struct {
	pb.UnimplementedAggregatorServiceServer

	cfg Config

	srv    *grpc.Server
	State  stateInterface
	Ethman etherman

	ctx  context.Context
	exit context.CancelFunc

	buildFinalProofBatchNumMutex *sync.Mutex
	buildFinalProofBatchNum      uint64

	startBatchNum uint64

	skippedsMutex *sync.Mutex
	skippeds      SequenceList
}

func newGenerateProof(cfg Config, stateInterface stateInterface, etherman etherman) *GenerateProof {

	return &GenerateProof{
		cfg: cfg,

		State:  stateInterface,
		Ethman: etherman,

		buildFinalProofBatchNumMutex: &sync.Mutex{},
		skippedsMutex:                &sync.Mutex{},
		skippeds:                     make([]state.Sequence, 0),
	}
}

func (g *GenerateProof) start(ctx context.Context) error {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel = context.WithCancel(ctx)
	g.ctx = ctx
	g.exit = cancel

	metrics.Register()

	address := fmt.Sprintf("%s:%d", g.cfg.Host, g.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	g.srv = grpc.NewServer()
	pb.RegisterAggregatorServiceServer(g.srv, g)

	healthService := newHealthChecker()
	grpchealth.RegisterHealthServer(g.srv, healthService)

	go func() {
		log.Infof("Server listening on port %d", g.cfg.Port)
		if err := g.srv.Serve(lis); err != nil {
			g.exit()
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	go g.checkGenerateFinalProof()

	<-ctx.Done()
	return ctx.Err()
}

func (g *GenerateProof) Channel(stream pb.AggregatorService_ChannelServer) error {
	metrics.ConnectedProver()
	defer metrics.DisconnectedProver()

	ctx := stream.Context()
	var proverAddr net.Addr
	p, ok := peer.FromContext(ctx)
	if ok {
		proverAddr = p.Addr
	}
	prover, err := prover.New(stream, proverAddr, g.cfg.ProofStatePollingInterval)
	if err != nil {
		return err
	}

	log := log.WithFields(
		"prover", prover.Name(),
		"proverId", prover.ID(),
		"proverAddr", prover.Addr(),
	)
	log.Info("Establishing stream connection with prover")

	// Check if prover supports the required Fork ID
	if !prover.SupportsForkID(g.cfg.ForkId) {
		err := errors.New("prover does not support required fork ID")
		log.Warn(FirstToUpper(err.Error()))
		return err
	}

	for {
		select {
		case <-g.ctx.Done():
			// server disconnected
			return g.ctx.Err()
		case <-ctx.Done():
			// client disconnected
			return ctx.Err()

		default:
			depoist, err := g.Ethman.JudgeAggregatorDeposit(common.HexToAddress(g.cfg.SenderAddress))
			if !depoist {
				if err != nil {
					log.Errorf("Failed to get deposit acount: %v", err)
				} else {
					log.Debugf("Check deposit amount. senderAddress: %s", g.cfg.SenderAddress)
				}

				time.Sleep(g.cfg.RetryTime.Duration)
				continue
			}

			isIdle, err := prover.IsIdle()
			if err != nil {
				log.Errorf("Failed to check if prover is idle: %v", err)
				time.Sleep(g.cfg.RetryTime.Duration)
				continue
			}
			if !isIdle {
				log.Debug("Prover is not idle")
				time.Sleep(g.cfg.RetryTime.Duration)
				continue
			}
		}
	}
}

func (g *GenerateProof) buildFinalProof(ctx context.Context, prover proverInterface, proof *state.Proof) (*pb.FinalProof, error) {
	log := log.WithFields(
		"prover", prover.Name(),
		"proverId", prover.ID(),
		"proverAddr", prover.Addr(),
		"batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal),
	)

	log.Info("Generating final proof")

	finalProofID, err := prover.FinalProof(proof.Proof, g.cfg.SenderAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get final proof id: %v", err)
	}
	proof.ProofID = finalProofID

	log.Infof("Final proof ID for batches [%d-%d]: %s", proof.BatchNumber, proof.BatchNumberFinal, *proof.ProofID)
	log = log.WithFields("finalProofId", finalProofID)

	finalProof, err := prover.WaitFinalProof(ctx, *proof.ProofID)
	if err != nil {
		return nil, fmt.Errorf("failed to get final proof from prover: %v", err)
	}

	monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, proof.BatchNumber, proof.BatchNumberFinal)

	stateFinalProof := state.FinalProof{
		MonitoredId:  monitoredTxID,
		FinalProof:   finalProof.Proof,
		FinalProofId: *finalProofID,
	}

	if err := g.State.AddFinalProof(g.ctx, &stateFinalProof, nil); err != nil {
		log.Error("failed to add final proof. state-monitoredTxID: %s, err = %v", monitoredTxID, err)
	}

	log.Info("Final proof generated")

	// mock prover sanity check
	if string(finalProof.Public.NewStateRoot) == mockedStateRoot && string(finalProof.Public.NewLocalExitRoot) == mockedLocalExitRoot {
		// This local exit root and state root come from the mock
		// prover, use the one captured by the executor instead
		finalBatch, err := g.State.GetBatchByNumber(ctx, proof.BatchNumberFinal, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve batch with number [%d]", proof.BatchNumberFinal)
		}
		log.Warnf("NewLocalExitRoot and NewStateRoot look like a mock values, using values from executor instead: LER: %v, SR: %v",
			finalBatch.LocalExitRoot.TerminalString(), finalBatch.StateRoot.TerminalString())
		finalProof.Public.NewStateRoot = finalBatch.StateRoot.Bytes()
		finalProof.Public.NewLocalExitRoot = finalBatch.LocalExitRoot.Bytes()
	}

	return finalProof, nil
}

func getFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func (g *GenerateProof) checkGenerateFinalProof() error {
	log := log.WithFields("function", getFunctionName())
	for {
		select {
		case <-g.ctx.Done():
			// server disconnected
			return g.ctx.Err()

		default:
			lastVerifiedBatch, err := g.State.GetLastVerifiedBatch(g.ctx, nil)
			if err != nil && !errors.Is(err, state.ErrNotFound) {
				log.Errorf("failed to get last verified batch, %v", err)
				continue
			}
			buildFinalProofBatchNum := g.buildFinalProofBatchNum
			if lastVerifiedBatch != nil && buildFinalProofBatchNum > lastVerifiedBatch.BatchNumber {
				batchNum := lastVerifiedBatch.BatchNumber
				for {
					sequence, err := g.State.GetSequence(g.ctx, batchNum+1, nil)
					if err != nil {
						log.Debugf("failed to get sequence err: %s. batchNum: %d", err, batchNum+1)
						continue
					}

					monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, sequence.FromBatchNumber, sequence.ToBatchNumber)
					_, errFinalProof := g.State.GetFinalProofByMonitoredId(g.ctx, monitoredTxID, nil)
					if errors.Is(errFinalProof, state.ErrNotFound) {
						g.skippedsMutex.Lock()
						g.skippeds = append(g.skippeds, sequence)
						sort.Sort(g.skippeds)
						g.skippedsMutex.Unlock()
						continue
					}

					if errFinalProof != nil {
						log.Debugf("failed to get final proof. err: %s. monitoredTxID: %s", err, monitoredTxID)
					}
				}
			}
		}
	}
}

func (g *GenerateProof) tryBuildFinalProof(ctx context.Context, prover proverInterface, proof *state.Proof) (bool, error) {
	proverName := prover.Name()
	proverID := prover.ID()

	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", prover.Addr(),
	)
	log.Debug("tryBuildFinalProof start")

	var err error
	log.Debug("Send final proof hash time reached")

	if proof == nil {

		defer func() {
			if err != nil {
				// Set the generating state to false for the proof ("unlock" it)
				proof.GeneratingSince = nil
				err2 := g.State.UpdateGeneratedProof(g.ctx, proof, nil)
				if err2 != nil {
					log.Errorf("Failed to unlock proof: %v", err2)
				}
			}
		}()
	} else {

	}

	return true, nil
}

func (g *GenerateProof) Stop() {
	g.exit()
	g.srv.Stop()
}

// healthChecker will provide an implementation of the HealthCheck interface.
type healthChecker struct{}

// newHealthChecker returns a health checker according to standard package
// grpc.health.v1.
func newHealthChecker() *healthChecker {
	return &healthChecker{}
}

// HealthCheck interface implementation.

// Check returns the current status of the server for unary gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (hc *healthChecker) Check(ctx context.Context, req *grpchealth.HealthCheckRequest) (*grpchealth.HealthCheckResponse, error) {
	log.Info("Serving the Check request for health check")
	return &grpchealth.HealthCheckResponse{
		Status: grpchealth.HealthCheckResponse_SERVING,
	}, nil
}

// Watch returns the current status of the server for stream gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (hc *healthChecker) Watch(req *grpchealth.HealthCheckRequest, server grpchealth.Health_WatchServer) error {
	log.Info("Serving the Watch request for health check")
	return server.Send(&grpchealth.HealthCheckResponse{
		Status: grpchealth.HealthCheckResponse_SERVING,
	})
}
