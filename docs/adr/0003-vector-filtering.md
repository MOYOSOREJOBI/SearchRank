# ADR 0003: Filtering strategy for vector retrieval

Filters (`source`, `date_range`, `tags`) must match BM25 semantics.

Current approach:
1. Apply filters in BM25 query.
2. Retrieve vector topN.
3. Apply same filters post-vector retrieval when ANN backend cannot filter.

Tradeoff: possible wasted vector candidates and slight recall loss after filtering.
