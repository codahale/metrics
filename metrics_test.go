package metrics_test

import (
	"testing"

	"github.com/codahale/metrics"
)

func TestCounter(t *testing.T) {
	metrics.Reset()

	metrics.Counter("whee").Add()
	metrics.Counter("whee").AddN(10)

	counters := metrics.Counters()
	if v, want := counters["whee"], uint64(11); v != want {
		t.Errorf("Counter was %v, but expected %v", v, want)
	}
}

func TestGaugeValue(t *testing.T) {
	metrics.Reset()

	metrics.Gauge("whee").Set(100.01)

	gauges := metrics.Gauges()
	if v, want := gauges["whee"], 100.01; v != want {
		t.Errorf("Gauge was %v, but expected %v", v, want)
	}
}

func TestGaugeFunc(t *testing.T) {
	metrics.Reset()

	metrics.Gauge("whee").SetFunc(func() float64 {
		return 100.03
	})

	gauges := metrics.Gauges()
	if v, want := gauges["whee"], 100.03; v != want {
		t.Errorf("Gauge was %v, but expected %v", v, want)
	}
}

func TestHistogram(t *testing.T) {
	metrics.Reset()

	h := metrics.NewHistogram("heyo", 1, 1000, 3)
	for i := 100; i > 0; i-- {
		for j := 0; j < i; j++ {
			h.RecordValue(int64(i))
		}
	}

	gauges := metrics.Gauges()

	if v, want := gauges["heyo.P50"], 71.0; v != want {
		t.Errorf("P50 was %v, but expected %v", v, want)
	}

	if v, want := gauges["heyo.P75"], 87.0; v != want {
		t.Errorf("P75 was %v, but expected %v", v, want)
	}

	if v, want := gauges["heyo.P90"], 95.0; v != want {
		t.Errorf("P90 was %v, but expected %v", v, want)
	}

	if v, want := gauges["heyo.P95"], 98.0; v != want {
		t.Errorf("P95 was %v, but expected %v", v, want)
	}

	if v, want := gauges["heyo.P99"], 100.0; v != want {
		t.Errorf("P99 was %v, but expected %v", v, want)
	}

	if v, want := gauges["heyo.P999"], 100.0; v != want {
		t.Errorf("P999 was %v, but expected %v", v, want)
	}
}

func BenchmarkCounterAdd(b *testing.B) {
	metrics.Reset()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.Counter("test1").Add()
		}
	})
}

func BenchmarkCounterAddN(b *testing.B) {
	metrics.Reset()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.Counter("test2").AddN(100)
		}
	})
}

func BenchmarkGaugeSet(b *testing.B) {
	metrics.Reset()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.Gauge("test2").Set(100)
		}
	})
}

func BenchmarkHistogramRecordValue(b *testing.B) {
	metrics.Reset()
	h := metrics.NewHistogram("hist", 1, 1000, 3)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.RecordValue(100)
		}
	})
}
