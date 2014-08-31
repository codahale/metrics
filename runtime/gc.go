package runtime

import (
	"runtime"

	"github.com/codahale/metrics"
)

func init() {
	msg := &memStatGauges{}

	metrics.Gauge("Mem.GC.Runs").SetBatchFunc("Mem", msg.init, msg.numGC)
	metrics.Gauge("Mem.GC.TotalPauseTime").SetBatchFunc("Mem", msg.init, msg.totalPause)
	metrics.Gauge("Mem.GC.LastPauseTime").SetBatchFunc("Mem", msg.init, msg.lastPause)
	metrics.Gauge("Mem.Alloc").SetBatchFunc("Mem", msg.init, msg.alloc)
	metrics.Gauge("Mem.Objects").SetBatchFunc("Mem", msg.init, msg.objects)
}

type memStatGauges struct {
	stats runtime.MemStats
}

func (msg *memStatGauges) init() {
	runtime.ReadMemStats(&msg.stats)
}

func (msg *memStatGauges) numGC() float64 {
	return float64(msg.stats.NumGC)
}

func (msg *memStatGauges) totalPause() float64 {
	return float64(msg.stats.PauseTotalNs)
}

func (msg *memStatGauges) lastPause() float64 {
	return float64(msg.stats.LastGC)
}

func (msg *memStatGauges) alloc() float64 {
	return float64(msg.stats.Alloc)
}

func (msg *memStatGauges) objects() float64 {
	return float64(msg.stats.HeapObjects)
}
