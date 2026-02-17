package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"searchrank/query-api/hybrid"
)

type SearchResult struct {
	DocID   string  `json:"doc_id"`
	Score   float64 `json:"score"`
	Explain struct {
		BM25      float64 `json:"bm25"`
		Vec       float64 `json:"vec"`
		Freshness float64 `json:"freshness"`
		LTR       float64 `json:"ltr"`
	} `json:"explain"`
}

type SearchResponse struct {
	RequestID string         `json:"request_id"`
	Results   []SearchResult `json:"results"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = time.Now().UTC().Format("20060102150405.000000")
	}
	bm25 := []hybrid.BM25Hit{{DocID: "doc-1", Score: 3.2}, {DocID: "doc-2", Score: 1.4}}
	vec := []hybrid.VecHit{{DocID: "doc-2", Score: 0.8}, {DocID: "doc-3", Score: 0.9}}
	merged := hybrid.MergeAndScore(bm25, vec, map[string]float64{"doc-1": 0.1, "doc-2": 0.2, "doc-3": 0.3}, hybrid.Weights{BM25: 0.6, Vec: 0.35, Freshness: 0.05})

	resp := SearchResponse{RequestID: requestID, Results: make([]SearchResult, 0, 10)}
	for i, d := range merged {
		if i >= 10 {
			break
		}
		var sr SearchResult
		sr.DocID = d.DocID
		sr.Score = d.Hybrid
		sr.Explain.BM25 = d.BM25
		sr.Explain.Vec = d.Vec
		sr.Explain.Freshness = d.Freshness
		sr.Explain.LTR = 0.0
		resp.Results = append(resp.Results, sr)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
	log.Printf("stage=score_finalize request_id=%s latency_ms=%d", requestID, time.Since(start).Milliseconds())
}

func main() {
	http.HandleFunc("/search", handler)
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	log.Fatal(http.ListenAndServe(":8080", nil))
}
