package ssf

import (
	"math/rand"
	"time"
)

// Samples is a batch of SSFSamples, not attached to an SSF span, that
// can be submitted with package metrics's Report function.
type Samples struct {
	Batch []SSFSample
}

// Add appends a sample to the batch of samples.
func (s *Samples) Add(sample ...SSFSample) {
	if s.Batch == nil {
		s.Batch = []SSFSample{}
	}
	s.Batch = append(s.Batch, sample...)
}

// NamePrefix is a string prepended to every SSFSample name generated
// by the constructors in this package. As no separator is added
// between this prefix and the metric name, users must take care to
// attach any separators to the prefix themselves.
var NamePrefix string

// Unit is a functional option for creating an SSFSample. It sets the
// sample's unit name to the name passed.
func Unit(name string) func(SSFSample) {
	return func(s SSFSample) {
		s.Unit = name
	}
}

// SampleRate sets the rate at which a measurement is sampled. The
// rate is a number on the interval (0..1] (1 means that the value is
// not sampled). Any numbers outside this interval result in no change
// to the sample rate (by default, all SSFSamples created with the
// helpers in this package have a SampleRate=1).
func SampleRate(rate float32) func(*SSFSample) {
	return func(s *SSFSample) {
		if rate > 0 && rate <= 1 {
			s.SampleRate = rate
		}
	}
}

var resolutions = map[time.Duration]string{
	time.Nanosecond:  "ns",
	time.Microsecond: "µs",
	time.Millisecond: "ms",
	time.Second:      "s",
	time.Minute:      "min",
	time.Hour:        "h",
}

// TimeUnit sets the unit on a sample to the given resolution's SI
// unit symbol. Valid resolutions are the time duration constants from
// Nanosecond through Hour. The non-SI units "minute" and "hour" are
// represented by "min" and "h" respectively.
//
// If a resolution is passed that does not correspond exactly to the
// duration constants in package time, this option does not affect the
// sample at all.
func TimeUnit(resolution time.Duration) func(*SSFSample) {
	return func(s *SSFSample) {
		if unit, ok := resolutions[resolution]; ok {
			s.Unit = unit
		}
	}
}

func create(base SSFSample) SSFSample {
	return base
}

// RandomlySample takes a rate and a set of measurements, and returns
// a new set of measurements as if sampling had been performed: Each
// original measurement gets rejected/included in the result based on
// a random roll of the RNG according to the rate, and each included
// measurement has its SampleRate field adjusted to be its original
// SampleRate * rate.
func RandomlySample(rate float32, samples ...SSFSample) []SSFSample {
	res := make([]SSFSample, 0, len(samples))

	for _, s := range samples {
		if rand.Float32() <= rate {
			if rate > 0 && rate <= 1 {
				s.SampleRate = s.SampleRate * rate
			}
			res = append(res, s)
		}
	}
	return res
}

// Count returns an SSFSample representing an increment / decrement of
// a counter. It's a convenience wrapper around constructing SSFSample
// objects.
func Count(name string, value float32, tags map[string]string) SSFSample {
	return create(SSFSample{
		Metric:     SSFSample_COUNTER,
		Name:       NamePrefix + name,
		Value:      value,
		Tags:       tags,
		SampleRate: 1.0,
	})
}

// Gauge returns an SSFSample representing a gauge at a certain
// value. It's a convenience wrapper around constructing SSFSample
// objects.
func Gauge(name string, value float32, tags map[string]string) SSFSample {
	return create(SSFSample{
		Metric:     SSFSample_GAUGE,
		Name:       NamePrefix + name,
		Value:      value,
		Tags:       tags,
		SampleRate: 1.0,
	})
}

// Histogram returns an SSFSample representing a value on a histogram,
// like a timer or other range. It's a convenience wrapper around
// constructing SSFSample objects.
func Histogram(name string, value float32, tags map[string]string) SSFSample {
	return create(SSFSample{
		Metric:     SSFSample_HISTOGRAM,
		Name:       NamePrefix + name,
		Value:      value,
		Tags:       tags,
		SampleRate: 1.0,
	})
}

// Set returns an SSFSample representing a value on a set, useful for
// counting the unique values that occur in a certain time bound.
func Set(name string, value string, tags map[string]string) SSFSample {
	return create(SSFSample{
		Metric:     SSFSample_SET,
		Name:       NamePrefix + name,
		Message:    value,
		Tags:       tags,
		SampleRate: 1.0,
	})
}

// TODO FIX THIS
// Timing returns an SSFSample (really a histogram) representing the
// timing in the given resolution.
func Timing(name string, value time.Duration, resolution time.Duration, tags map[string]string) SSFSample {
	time := float32(value / resolution)
	ssfsample := Histogram(name, time, tags)

	TimeUnit(resolution)(&ssfsample)

	return ssfsample

}
