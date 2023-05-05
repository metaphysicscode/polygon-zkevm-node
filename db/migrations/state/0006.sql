CREATE TABLE state.proof_hash
(
    id              SERIAL PRIMARY KEY,
    block_num       BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE,
    sender          VARCHAR NOT NULL,
    init_num_batch  BIGINT NOT NULL,
    final_new_batch BIGINT NOT NULL,
    proof_hash      VARCHAR NOT NULL
);