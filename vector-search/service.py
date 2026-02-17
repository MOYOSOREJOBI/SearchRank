#!/usr/bin/env python3
"""Vector retrieval service skeleton.
Loads FAISS index/versioned artifacts and exposes Search/Health handlers.
"""
import json
from pathlib import Path
from typing import Dict, List

import numpy as np
from google.protobuf.json_format import MessageToDict

from vector_search_pb2 import HealthResponse, SearchResponse, SearchResult

try:
    import faiss
except Exception as exc:
    raise SystemExit(f"faiss import failed: {exc}")


class VectorSearchService:
    def __init__(self, artifact_dir: str = "artifacts/index/v1"):
        self.artifact_dir = Path(artifact_dir)
        self.index = faiss.read_index(str(self.artifact_dir / "faiss.index"))
        self.id_map = json.loads((self.artifact_dir / "id_map.json").read_text())
        self.meta = json.loads((self.artifact_dir / "embedder_metadata.json").read_text())

    def health(self) -> HealthResponse:
        return HealthResponse(status="ok", loaded_index_version=self.meta["version"])

    def search(self, query_embedding: List[float], top_n: int, filters: Dict[str, str]) -> SearchResponse:
        q = np.array([query_embedding], dtype=np.float32)
        scores, ids = self.index.search(q, top_n)
        results = []
        for score, idx in zip(scores[0], ids[0]):
            if idx < 0:
                continue
            row = self.id_map[int(idx)]
            results.append(SearchResult(chunk_id=row["chunk_id"], doc_id=row["doc_id"], vec_score=float(score)))
        return SearchResponse(results=results, index_version=self.meta["version"])


if __name__ == "__main__":
    svc = VectorSearchService()
    print(json.dumps(MessageToDict(svc.health()), indent=2))
