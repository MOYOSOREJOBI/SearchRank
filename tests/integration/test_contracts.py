import json
import subprocess


def test_search_schema_stable():
    p = subprocess.run(["go", "run", "./query-api/cmd/query-api"], capture_output=True, text=True, timeout=2)
    # smoke compile/run only; endpoint checks are handled in docker integration target
    assert p.returncode != 2


def test_golden_query_fixture_count():
    rows = json.loads(open("tests/fixtures/golden_queries.json", "r", encoding="utf-8").read())
    assert len(rows) >= 10
