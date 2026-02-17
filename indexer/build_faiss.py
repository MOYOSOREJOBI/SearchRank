#!/usr/bin/env python3
import argparse
import hashlib
import json
from pathlib import Path

import numpy as np

try:
    import faiss
except Exception as exc:
    raise SystemExit(f"faiss import failed: {exc}")

EMBEDDER = {
    "name": "deterministic-hash-embedder",
    "version": "1.0.0",
    "model_hash": "sha256:2cb170fbb6d2d7f6f0f7b221b2f5f6b94e9a7b55ecf31fd5fe2f5b295f020e95",
    "dim": 64,
}


def embed(text: str, dim: int) -> np.ndarray:
    out = np.zeros(dim, dtype=np.float32)
    for i in range(dim):
        digest = hashlib.sha256(f"{text}|{i}".encode()).digest()
        out[i] = int.from_bytes(digest[:4], "big") / 2**32
    norm = np.linalg.norm(out)
    return out if norm == 0 else out / norm


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", default="tests/fixtures/chunks.jsonl")
    parser.add_argument("--out", default="artifacts/index/v1")
    args = parser.parse_args()

    rows = [json.loads(x) for x in Path(args.input).read_text().splitlines() if x.strip()]
    out_dir = Path(args.out)
    out_dir.mkdir(parents=True, exist_ok=True)

    mat = np.vstack([embed(r["text"], EMBEDDER["dim"]) for r in rows]).astype(np.float32)
    index = faiss.IndexHNSWFlat(EMBEDDER["dim"], 32)
    index.hnsw.efConstruction = 80
    index.add(mat)
    faiss.write_index(index, str(out_dir / "faiss.index"))

    id_map = [{"faiss_internal_id": i, "chunk_id": r["chunk_id"], "doc_id": r["doc_id"]} for i, r in enumerate(rows)]
    (out_dir / "id_map.json").write_text(json.dumps(id_map, indent=2))
    (out_dir / "embedder_metadata.json").write_text(json.dumps(EMBEDDER, indent=2))

    pg_rows = [
        {
            "chunk_id": r["chunk_id"],
            "doc_id": r["doc_id"],
            "embedding_checksum": hashlib.sha256(mat[i].tobytes()).hexdigest(),
            "embedder_version": EMBEDDER["version"],
        }
        for i, r in enumerate(rows)
    ]
    (out_dir / "pg_vectors.json").write_text(json.dumps(pg_rows, indent=2))


if __name__ == "__main__":
    main()
