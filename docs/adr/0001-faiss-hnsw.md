# ADR 0001: FAISS index choice (HNSW)

## Decision
Use `faiss.IndexHNSWFlat` as the first production ANN index for semantic retrieval.

## Why
- No training phase required.
- Deterministic build from fixed input embeddings.
- Good recall/latency tradeoff for medium corpora.

## Consequences
- Higher memory footprint than IVF-PQ.
- Future milestone can shard or quantize.
