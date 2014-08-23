// Package metrics provides minimalist instrumentation for your applications in
// the form of counters and gauges.
//
// Measurements from counters and gauges are available as expvars.
package metrics

import (
	"expvar"
	"sync"
	"time"

	"github.com/codahale/hdrhistogram/hdr"
)

// A Counter is a monotonically increasing unsigned integer.
//
// Use a counter to derive rates (e.g., record total number of requests, derive
// requests per second).
type Counter string

// Add increments the counter by one.
func (c Counter) Add() {
	c.AddN(1)
}

// AddN increments the counter by N.
func (c Counter) AddN(delta uint64) {
	cm.Lock()
	counters[string(c)] += delta
	cm.Unlock()
}

// A Gauge is an instantaneous measurement of a value.
//
// Use a gauge to track metrics which increase and decrease (e.g., amount of
// free memory).
type Gauge string

// Set the gauge's value to the given value.
func (g Gauge) Set(value float64) {
	gm.Lock()
	gauges[string(g)] = func() float64 {
		return value
	}
	gm.Unlock()
}

// SetFunc sets the gauge's value to the lazily-called return value of the given
// function.
func (g Gauge) SetFunc(f func() float64) {
	gm.Lock()
	gauges[string(g)] = f
	gm.Unlock()
}

// Reset removes all existing counters and gauges.
func Reset() {
	cm.Lock()
	defer cm.Unlock()

	gm.Lock()
	defer gm.Unlock()

	hm.Lock()
	defer hm.Unlock()

	counters = make(map[string]uint64)
	gauges = make(map[string]func() float64)
	histograms = make(map[string]*Histogram)

}

// Counters returns a snapshot of the current values of all counters.
func Counters() map[string]uint64 {
	cm.Lock()
	defer cm.Unlock()

	c := make(map[string]uint64, len(counters))
	for n, v := range counters {
		c[n] = v
	}
	return c
}

// Gauges returns a snapshot of the current values of all gauges.
func Gauges() map[string]float64 {
	gm.Lock()
	defer gm.Unlock()

	hm.Lock()
	defer hm.Unlock()

	// include room for histogram values
	g := make(map[string]float64, len(gauges)+(len(histograms)*6))

	for n, f := range gauges {
		g[n] = f()
	}

	for n, h := range histograms {
		m := h.merge()
		g[n+".P50"] = float64(m.ValueAtQuantile(50))
		g[n+".P75"] = float64(m.ValueAtQuantile(75))
		g[n+".P90"] = float64(m.ValueAtQuantile(90))
		g[n+".P95"] = float64(m.ValueAtQuantile(95))
		g[n+".P99"] = float64(m.ValueAtQuantile(99))
		g[n+".P999"] = float64(m.ValueAtQuantile(99.9))
	}

	return g
}

// NewHistogram returns a windowed HDR histogram which drops data older than
// five minutes.
//
// Use a histogram to track the distribution of a stream of values (e.g., the
// latency associated with HTTP requests).
func NewHistogram(name string, minValue, maxValue int64, sigfigs int) *Histogram {
	hm.Lock()
	defer hm.Unlock()

	if _, ok := histograms[name]; ok {
		panic(name + " already exists")
	}

	histograms[name] = &Histogram{
		hist: hdr.NewWindowedHistogram(5, minValue, maxValue, sigfigs),
	}
	return histograms[name]
}

// A Histogram measures the distribution of a stream of values.
type Histogram struct {
	hist *hdr.WindowedHistogram
	rw   sync.RWMutex
}

// RecordValue records the given value.
func (h *Histogram) RecordValue(v int64) {
	h.rw.Lock()
	defer h.rw.Unlock()

	h.hist.Current.RecordValue(v)
}

func (h *Histogram) rotate() {
	h.rw.Lock()
	defer h.rw.Unlock()

	h.hist.Rotate()
}

func (h *Histogram) merge() *hdr.Histogram {
	h.rw.RLock()
	defer h.rw.RUnlock()

	return h.hist.Merge()
}

var (
	counters   map[string]uint64
	gauges     map[string]func() float64
	histograms map[string]*Histogram

	cm, gm, hm sync.Mutex
)

func init() {
	Reset()

	expvar.Publish("counters", expvar.Func(func() interface{} {
		return Counters()
	}))

	expvar.Publish("gauges", expvar.Func(func() interface{} {
		return Gauges()
	}))

	go func() {
		for _ = range time.NewTicker(1 * time.Minute).C {
			gm.Lock()
			for _, h := range histograms {
				h.rotate()
			}
			gm.Unlock()
		}
	}()
}
