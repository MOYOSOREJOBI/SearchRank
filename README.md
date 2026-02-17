# SearchRank Milestone 2: Vector Retrieval + Hybrid Candidate Generation

## Architecture / Data flow
1. `indexer/build_faiss.py` reads chunk fixtures and generates deterministic embeddings.
2. Indexer writes versioned artifacts in `artifacts/index/v1/`:
   - `faiss.index`
   - `id_map.json`
   - `embedder_metadata.json`
   - `pg_vectors.json`
3. `vector-search/service.py` loads FAISS + metadata on startup and serves typed search contracts from `proto/vector_search.proto`.
4. `query-api` does:
   - BM25 retrieve (`bm25_retrieve` span)
   - query embedding (`embed_query` span placeholder)
   - vector retrieve with timeout + graceful fallback (`faiss_retrieve`)
   - merge/dedupe (`merge_dedupe`)
   - normalized weighted scoring (`score_finalize`)

## Determinism controls
- Pinned Python deps in `indexer/requirements.txt` and `vector-search/requirements.txt`.
- Explicit embedder metadata with name/version/hash.
- Versioned artifacts under `artifacts/index/v1`.

## API schema stability
`/search` always returns:
- `explain.bm25`
- `explain.vec`
- `explain.freshness`
- `explain.ltr` (fixed `0.0` until LTR milestone)

## Commands
- `make dev`
- `make test`
- `make benchmark`
- `make reindex`

## Integration tests
Run:
- `make reindex`
- `docker compose up -d`
- `curl -s localhost:8080/search | jq .`

## Benchmark
`make benchmark` prints p50/p95 and fails if p95 > 250ms.

## Example response
```json
{
  "request_id": "20260212010101.123456",
  "results": [
    {
      "doc_id": "doc-1",
      "score": 0.71,
      "explain": {
        "bm25": 1.0,
        "vec": 0.2,
        "freshness": 0.1,
        "ltr": 0.0
      }
    }
  ]
}
```

## Acceptance checklist
- [ ] `docker compose up` brings up all services cleanly
- [ ] Indexer builds OpenSearch index + FAISS artifact deterministically
- [ ] `/search` returns stable JSON schema with explain fields
- [ ] Vector down => BM25 fallback + valid schema
- [ ] Integration tests pass in CI
- [ ] Benchmark prints p50/p95 and meets targets
- [ ] ADRs exist and match implementation
