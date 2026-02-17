# Architecture Overview

SearchRank is organized as a modular monorepo with strict contract-first boundaries.

## Logical components

- **Web app** (`/apps/web`): query input, result rendering, explain/debug views.
- **Query API** (`/services/query-api`): orchestrates retrieval, reranking, and logging.
- **Vector service** (`/services/vector-search`): loads FAISS indexes and serves ANN search.
- **Indexer job** (`/jobs/indexer`): ingestion, chunking, embedding, lexical/vector index building.
- **Trainer job** (`/jobs/trainer`): label generation, LTR training, evaluation, model packaging.
- **Contracts package** (`/pkg/contracts`): OpenAPI and schemas + generated clients.
- **Infra** (`/infra`): compose stacks, Kubernetes deployment manifests, Terraform IaC.

## Request lifecycle

1. Client sends query and optional filters/session metadata.
2. Query API fans out to BM25 and vector retrieval.
3. Candidates are merged and deduplicated.
4. Feature vectors are built for top candidate set.
5. LTR reranker scores top-K.
6. Top results with explain payload are returned.
7. Impression event is logged; click events are ingested asynchronously.

## Reliability and operability principles

- Contract compatibility before deploy.
- Metrics + tracing required for all critical hops.
- Feature schema hash check on model load.
- Fallback behavior documented for dependency degradation.
