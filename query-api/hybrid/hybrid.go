package hybrid

import (
	"context"
	"errors"
	"sort"
	"time"
)

type BM25Hit struct {
	DocID string
	Score float64
}

type VecHit struct {
	DocID string
	Score float64
}

type DocFeatures struct {
	DocID      string
	BM25       float64
	Vec        float64
	Freshness  float64
	Hybrid     float64
	VectorUsed bool
}

type Weights struct {
	BM25      float64
	Vec       float64
	Freshness float64
}

type VectorClient interface {
	Search(ctx context.Context, queryEmbedding []float32, topN int, filters map[string]string) ([]VecHit, error)
}

var ErrVectorUnavailable = errors.New("vector unavailable")

func MinMaxNormalize(values []float64) []float64 {
	out := make([]float64, len(values))
	if len(values) == 0 {
		return out
	}
	mn, mx := values[0], values[0]
	for _, v := range values {
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	eps := 1e-9
	den := mx - mn
	if den < eps {
		for i := range out {
			out[i] = 0.0
		}
		return out
	}
	for i, v := range values {
		n := (v - mn) / den
		if n < 0 {
			n = 0
		}
		if n > 1 {
			n = 1
		}
		out[i] = n
	}
	return out
}

func MergeAndScore(bm25 []BM25Hit, vec []VecHit, freshness map[string]float64, w Weights) []DocFeatures {
	bm25Map := map[string]float64{}
	vecMap := map[string]float64{}
	for _, h := range bm25 {
		if h.Score > bm25Map[h.DocID] {
			bm25Map[h.DocID] = h.Score
		}
	}
	for _, h := range vec {
		if h.Score > vecMap[h.DocID] {
			vecMap[h.DocID] = h.Score
		}
	}

	docIDs := map[string]struct{}{}
	for id := range bm25Map {
		docIDs[id] = struct{}{}
	}
	for id := range vecMap {
		docIDs[id] = struct{}{}
	}

	ids := make([]string, 0, len(docIDs))
	bmVals := make([]float64, 0, len(docIDs))
	vecVals := make([]float64, 0, len(docIDs))
	for id := range docIDs {
		ids = append(ids, id)
		bmVals = append(bmVals, bm25Map[id])
		vecVals = append(vecVals, vecMap[id])
	}
	bmNorm := MinMaxNormalize(bmVals)
	vecNorm := MinMaxNormalize(vecVals)

	res := make([]DocFeatures, 0, len(ids))
	for i, id := range ids {
		f := freshness[id]
		h := w.BM25*bmNorm[i] + w.Vec*vecNorm[i] + w.Freshness*f
		res = append(res, DocFeatures{DocID: id, BM25: bmNorm[i], Vec: vecNorm[i], Freshness: f, Hybrid: h, VectorUsed: vecMap[id] > 0})
	}
	sort.SliceStable(res, func(i, j int) bool {
		if res[i].Hybrid == res[j].Hybrid {
			return res[i].DocID < res[j].DocID
		}
		return res[i].Hybrid > res[j].Hybrid
	})
	return res
}

func RetrieveVecWithTimeout(ctx context.Context, client VectorClient, embedding []float32, topN int, filters map[string]string, timeout time.Duration) ([]VecHit, error) {
	vctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	res, err := client.Search(vctx, embedding, topN, filters)
	if err != nil {
		return nil, ErrVectorUnavailable
	}
	return res, nil
}
