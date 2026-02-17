# ADR 0002: Hybrid normalization and weighting

Per-request min-max normalization for BM25 and vector scores with epsilon guard (`1e-9`), clamped to [0,1].

`hybrid = w_bm25*bm25_norm + w_vec*vec_norm + w_fresh*freshness`

Default weights: `w_bm25=0.6`, `w_vec=0.35`, `w_fresh=0.05`.
