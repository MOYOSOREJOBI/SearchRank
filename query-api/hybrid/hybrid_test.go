package hybrid

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeVector struct {
	out   []VecHit
	err   error
	delay time.Duration
}

func (f fakeVector) Search(ctx context.Context, _ []float32, _ int, _ map[string]string) ([]VecHit, error) {
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if f.err != nil {
		return nil, f.err
	}
	return f.out, nil
}

func TestMergeAndDedupe(t *testing.T) {
	bm := []BM25Hit{{DocID: "a", Score: 5}, {DocID: "a", Score: 3}, {DocID: "b", Score: 2}}
	vec := []VecHit{{DocID: "b", Score: 0.2}, {DocID: "c", Score: 0.9}}
	res := MergeAndScore(bm, vec, map[string]float64{"a": 0.2, "b": 0.1, "c": 0.4}, Weights{BM25: 0.5, Vec: 0.4, Freshness: 0.1})
	if len(res) != 3 {
		t.Fatalf("want 3 docs got %d", len(res))
	}
	if res[0].DocID != "a" && res[0].DocID != "c" {
		t.Fatalf("unexpected top doc %s", res[0].DocID)
	}
}

func TestNormalizationDeterministic(t *testing.T) {
	in := []float64{1, 2, 3}
	a := MinMaxNormalize(in)
	b := MinMaxNormalize(in)
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("non deterministic at %d", i)
		}
	}
}

func TestFallbackWhenVectorSlowOrDown(t *testing.T) {
	ctx := context.Background()
	_, err := RetrieveVecWithTimeout(ctx, fakeVector{delay: 50 * time.Millisecond}, []float32{1}, 10, nil, 10*time.Millisecond)
	if !errors.Is(err, ErrVectorUnavailable) {
		t.Fatalf("expected fallback error got %v", err)
	}

	_, err = RetrieveVecWithTimeout(ctx, fakeVector{err: errors.New("down")}, []float32{1}, 10, nil, 10*time.Millisecond)
	if !errors.Is(err, ErrVectorUnavailable) {
		t.Fatalf("expected fallback error got %v", err)
	}
}
