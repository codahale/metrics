package runtime

import (
	"testing"

	"github.com/codahale/metrics"
)

func TestMemStats(t *testing.T) {
	_, gauges := metrics.Snapshot()

	expected := []string{
		"Mem.NumGC",
		"Mem.PauseTotalNs",
		"Mem.LastGC",
		"Mem.Alloc",
		"Mem.HeapObjects",
	}

	for _, name := range expected {
		if _, ok := gauges[name]; !ok {
			t.Errorf("Missing gauge %q", name)
		}
	}
}
