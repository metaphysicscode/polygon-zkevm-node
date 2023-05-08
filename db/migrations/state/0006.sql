-- +migrate Down
DROP table state.proof_hash;
DROP table state.prover_proof;

-- +migrate Up
CREATE TABLE IF NOT EXISTS state.proof_hash
(
    id              SERIAL PRIMARY KEY,
    block_num       BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE,
    sender          VARCHAR NOT NULL,
    init_num_batch  BIGINT NOT NULL,
    final_new_batch BIGINT NOT NULL,
    proof_hash      VARCHAR NOT NULL
);


CREATE TABLE IF NOT EXISTS state.prover_proof
(
    id              SERIAL PRIMARY KEY,
    init_num_batch  BIGINT NOT NULL,
    final_new_batch BIGINT NOT NULL,
    local_exit_root VARCHAR,
    new_state_root  VARCHAR,
    proof           VARCHAR,
    proof_hash      VARCHAR NOT NULL
);