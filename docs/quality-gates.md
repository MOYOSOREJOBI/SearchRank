# Quality Gates

These gates define minimum quality for merge and deployment.

## Merge gates

1. **Static checks**
   - Lint and formatting checks for Go, TypeScript, and Python.
2. **Automated tests**
   - Unit tests.
   - Integration tests with containerized OpenSearch + Postgres dependencies.
3. **Contracts**
   - OpenAPI/schema compatibility checks.
   - Generated client freshness checks.
4. **Performance**
   - End-to-end search latency p95 < 250ms.
   - Candidate generation p95 < 120ms.
   - Rerank stage p95 < 60ms.
5. **Relevance**
   - No regression past tolerated threshold for nDCG@10, MRR@10, Recall@50 on the golden set.

## Deployment gates

- Model promotion only if offline metrics pass threshold and schema hash matches online feature schema.
- Rollback target version must remain deployable and warm.
- Error budget burn alerts must be green.

## Failure policy

- Any gate failure blocks merge/promotion.
- Waivers require documented ADR/runbook entry and expiration date.
