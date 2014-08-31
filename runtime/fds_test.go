package runtime

import (
	"testing"

	"github.com/codahale/metrics"
)

func TestFdStats(t *testing.T) {
	g := metrics.Gauges()

	expected := []string{
		"FileDescriptors.Max",
		"FileDescriptors.Used",
	}

	for _, name := range expected {
		if _, ok := g[name]; !ok {
			t.Errorf("Missing gauge %q", name)
		}
	}
}
