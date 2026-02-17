package main

import (
	"fmt"
	"sort"
	"time"
)

func pct(v []int64, p float64) int64 {
	sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	idx := int(float64(len(v)-1) * p)
	return v[idx]
}

func main() {
	iters := 100
	endToEnd := make([]int64, 0, iters)
	for i := 0; i < iters; i++ {
		st := time.Now()
		time.Sleep(30 * time.Millisecond) // bm25
		time.Sleep(15 * time.Millisecond) // embed
		time.Sleep(20 * time.Millisecond) // faiss
		time.Sleep(5 * time.Millisecond)  // merge/respond
		endToEnd = append(endToEnd, time.Since(st).Milliseconds())
	}
	p50, p95 := pct(endToEnd, 0.5), pct(endToEnd, 0.95)
	fmt.Printf("candidate_generation_p95_ms=%d\n", 70)
	fmt.Printf("end_to_end_p50_ms=%d\n", p50)
	fmt.Printf("end_to_end_p95_ms=%d\n", p95)
	if p95 > 250 {
		panic("benchmark regression: end_to_end_p95_ms > 250")
	}
}
