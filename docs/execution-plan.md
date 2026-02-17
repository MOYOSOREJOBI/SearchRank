# SearchRank Execution Plan

This plan defines the product and engineering bar for SearchRank and the build sequence from baseline retrieval to continuously improving LTR.

## 0) North-star product bar

### Non-negotiable promises

1. **Trustworthy ranking:** explainable and auditable ranking behavior.
2. **Fast:** p95 end-to-end latency < 250ms; rerank top-K (K=50) < 60ms.
3. **Measurable:** `nDCG@10`, `MRR@10`, `Recall@50` tracked and used as release gates.
4. **Reproducible:** deterministic builds, pinned dependencies, versioned artifacts.
5. **Operable:** SLOs + error budgets govern rollouts.

### Differentiator: SearchRank Studio

- Side-by-side ranker diff for identical query input.
- Candidate set and feature vector comparisons.
- LTR score delta inspection.
- Per-request trace drilldown across retrieval, feature generation, and reranking.

## 1) Architecture boundaries

### Services

1. **Query API (Go)**
   - Validation, bucketing, retrieval orchestration, merge/dedupe, feature generation, rerank, response shaping, event logging.
2. **BM25 retrieval (OpenSearch/Elasticsearch)**
   - Lexical retrieval, filtering, snippets/highlights.
3. **Vector retrieval (FAISS microservice)**
   - ANN index load and query serving.
4. **Metadata/events store (Postgres)**
   - Metadata, chunk map, impressions/clicks, experiments, model registry.
5. **Offline jobs (Python)**
   - Indexing pipeline and training/evaluation pipeline.

### Online flow

1. BM25 top-N retrieval.
2. Vector top-N retrieval.
3. Merge + dedupe.
4. Feature build.
5. Rerank top-K = 50.
6. Return top-10 with explain fields.
7. Log impression/click events with trace/request IDs.

## 2) Engineering standards

- Go, TypeScript, and Python style/lint discipline enforced in CI.
- OpenAPI and JSON schema are source-of-truth contracts.
- No opaque scoring: all score components exposed in logs and explain payload.
- Containerized integration tests to prevent machine-specific drift.

## 3) Canonical repository layout

```text
/apps/web
/services/query-api
/services/vector-search
/jobs/indexer
/jobs/trainer
/pkg/contracts
/infra/docker
/infra/k8s
/infra/terraform
/docs/adr
/docs/runbooks
```

## 4) Data contracts and model safety

### Public Search API

- Request: `query`, `filters`, `client_id`, optional `session_id`.
- Response: `query`, `took_ms`, `experiment`, `results[]` (including explain fields).

### Event APIs

- Impression event: request_id, query, ranked doc IDs/positions, bucket, ranker_version.
- Click event: request_id, doc_id, position, dwell_ms, timestamp.

### Model registry constraints

- Persist model version, training window, metrics, feature schema hash, and artifact URI.
- Online serving must reject model loads when feature schema hash mismatches.

## 5) Feature design v1

Per `(query, candidate)` features:

- `bm25_score`
- `vec_similarity`
- `freshness`
- `field_match`
- `length_norm`
- `query_len`
- `term_overlap`

All features are logged and eligible for explain surfaces.

## 6) Labeling and bias correction

### Initial labels

- Positive: clicked with dwell >= configured threshold.
- Negative: shown above clicked item but not clicked (guarded heuristics).

### Bias mitigation path

- Run a controlled randomization slice to estimate propensities.
- Move to propensity-weighted LTR once data volume is sufficient.

## 7) Milestones and acceptance criteria

### M0 — Foundation

- Contracts-first API specs, generated clients, CI skeleton, tracing scaffold.
- Acceptance: stable hello-search schema + end-to-end trace propagation.

### M1 — BM25 baseline

- Indexer v1, BM25 query path, filters/snippets, baseline UI.
- Acceptance: rank-eval pipeline wired + p95 baseline measured.

### M2 — Vector retrieval

- Embeddings + FAISS artifact + ANN serving.
- Acceptance: recall gains on semantic sets + index choice benchmark record.

### M3 — Merge and feature builder

- Deterministic merge/dedupe and versioned feature vectors.
- Acceptance: stable feature schema hash + replay determinism tests.

### M4 — Offline LTR pipeline

- LightGBM LambdaMART training/eval + model card + artifact registration.
- Acceptance: reproducibility and metric non-regression gates.

### M5 — Online reranking

- Top-K (50) online reranker with explain score composition.
- Acceptance: rerank p95 < 60ms; API schema unchanged.

### M6 — Experiments and logging

- Sticky bucketing, impression/click logging, UI experiment debugger.
- Acceptance: response includes experiment metadata; SQL analytics samples verified.

### M7 — Retraining loop

- Scheduled retraining + promote-only-if-better + immediate rollback path.
- Acceptance: model lineage and gated promotion documented and enforced.

### M8 — Hardening

- Load/perf tests, k8s probes, Terraform backend discipline, SLO alerting.
- Acceptance: CI blocks latency/relevance regressions + incident runbooks complete.

## 8) CI quality gates

Required for merge:

- Lint + unit tests + integration tests.
- Contract compatibility tests.
- Latency SLO checks.
- Relevance non-regression checks.

## 9) Observability baseline

Every request trace includes spans for:

- `opensearch_query`
- `faiss_search`
- `merge_dedupe`
- `feature_build`
- `rerank`
- `postgres_log`

Metrics include stage latencies, candidate counts, dedupe rates, rerank failures, model load timings, CTR, and dwell distributions by experiment bucket.

## 10) Documentation expectations

README and docs should always include:

- Architecture and dataflow diagram.
- SLOs and measurement method.
- Relevance metrics and golden set methodology.
- Latency benchmark context.
- Reproducible run/eval commands.
- ADRs for key relevance and infrastructure decisions.
