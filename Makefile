.PHONY: dev test benchmark reindex integration

dev:
	docker compose up --build

test:
	go test ./...
	python3 -m pytest tests/integration -q

benchmark:
	go run ./benchmarks/harness.go

reindex:
	python3 indexer/build_faiss.py --input tests/fixtures/chunks.jsonl --out artifacts/index/v1

integration:
	docker compose up -d
	curl -s localhost:8080/search | jq .
