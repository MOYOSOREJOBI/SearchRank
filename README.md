# SearchRank

SearchRank is a search relevance platform built with production-grade engineering standards: strict API contracts, explainable ranking, measurable quality and latency gates, and a disciplined online/offline iteration loop.

## Product goals

- Trustworthy and explainable ranking decisions.
- Fast user experience (`p95 < 250ms` end-to-end, rerank top-50 `p95 < 60ms`).
- Relevance metrics tracked continuously (`nDCG@10`, `MRR@10`, `Recall@50`).
- Reproducible training and deployment with pinned dependencies and versioned artifacts.
- Operable via explicit SLOs and error budgets.

## Execution plan

The full milestone-by-milestone plan is in [`docs/execution-plan.md`](docs/execution-plan.md).

Highlights:

1. **Foundation:** Contracts-first APIs, CI skeleton, tracing/metrics scaffold.
2. **BM25 baseline:** OpenSearch retrieval and measurable baseline relevance.
3. **Hybrid retrieval:** Add FAISS ANN and deterministic candidate merge.
4. **LTR pipeline:** Train/evaluate LightGBM LambdaMART with model registry.
5. **Online reranking:** Top-K reranker with explainability and latency budgets.
6. **Experimentation:** Bucketed experiments and impression/click logging.
7. **Retraining loop:** Promote-only-if-better automation.
8. **Hardening:** Load/perf tests, k8s probes, Terraform state discipline, SLO-driven operations.

## SearchRank Studio (differentiator)

SearchRank Studio is the built-in relevance debugger:

- Side-by-side ranker comparison for the same query.
- Candidate/feature/score diff views.
- Request trace linkage to retrieval/rerank timings.
- Human-readable “Why this ranked” breakdown backed by logged feature contributions.

## Target architecture

```text
/apps/web                 Next.js UI
/services/query-api       Go orchestration API
/services/vector-search   FAISS ANN service
/jobs/indexer             Python ingestion/index builds
/jobs/trainer             Python LTR training/evaluation
/pkg/contracts            OpenAPI + JSON schemas + generated clients
/infra                    docker/k8s/terraform
/docs                     ADRs, runbooks, planning
```

See details in [`docs/architecture.md`](docs/architecture.md).

## Quality gates

Merge is blocked unless all required checks pass:

- Unit + integration tests
- Contract compatibility checks
- Latency gates (`p95` thresholds)
- Relevance non-regression gates on golden query sets

See [`docs/quality-gates.md`](docs/quality-gates.md) for exact thresholds and enforcement policy.
