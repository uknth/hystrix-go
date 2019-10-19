package metric

import (
	"sync"

	"github.com/unbxd/hystrix-go/hystrix/rolling"
)

// DefaultCollector holds information about the circuit state.
// This implementation of Collector is the canonical source of information about the circuit.
// It is used for for all internal hystrix operations
// including circuit health checks and metrics sent to the hystrix dashboard.
//
// Metric Collectors do not need Mutexes as they are updated by circuits within a locked context.
type DefaultCollector struct {
	mutex *sync.RWMutex

	numRequests *rolling.Number
	errors      *rolling.Number

	successes               *rolling.Number
	failures                *rolling.Number
	rejects                 *rolling.Number
	shortCircuits           *rolling.Number
	timeouts                *rolling.Number
	contextCanceled         *rolling.Number
	contextDeadlineExceeded *rolling.Number

	fallbackSuccesses *rolling.Number
	fallbackFailures  *rolling.Number
	totalDuration     *rolling.Timing
	runDuration       *rolling.Timing
}

func newDefaultCollector(name string) Collector {
	m := &DefaultCollector{}
	m.mutex = &sync.RWMutex{}
	m.Reset()
	return m
}

// NumRequests returns the rolling number of requests
func (d *DefaultCollector) NumRequests() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.numRequests
}

// Errors returns the rolling number of errors
func (d *DefaultCollector) Errors() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.errors
}

// Successes returns the rolling number of successes
func (d *DefaultCollector) Successes() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.successes
}

// Failures returns the rolling number of failures
func (d *DefaultCollector) Failures() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.failures
}

// Rejects returns the rolling number of rejects
func (d *DefaultCollector) Rejects() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.rejects
}

// ShortCircuits returns the rolling number of short circuits
func (d *DefaultCollector) ShortCircuits() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.shortCircuits
}

// Timeouts returns the rolling number of timeouts
func (d *DefaultCollector) Timeouts() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.timeouts
}

// FallbackSuccesses returns the rolling number of fallback successes
func (d *DefaultCollector) FallbackSuccesses() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.fallbackSuccesses
}

func (d *DefaultCollector) ContextCanceled() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.contextCanceled
}

func (d *DefaultCollector) ContextDeadlineExceeded() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.contextDeadlineExceeded
}

// FallbackFailures returns the rolling number of fallback failures
func (d *DefaultCollector) FallbackFailures() *rolling.Number {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.fallbackFailures
}

// TotalDuration returns the rolling total duration
func (d *DefaultCollector) TotalDuration() *rolling.Timing {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.totalDuration
}

// RunDuration returns the rolling run duration
func (d *DefaultCollector) RunDuration() *rolling.Timing {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.runDuration
}

func (d *DefaultCollector) Update(r Result) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	d.numRequests.Increment(r.Attempts)
	d.errors.Increment(r.Errors)
	d.successes.Increment(r.Successes)
	d.failures.Increment(r.Failures)
	d.rejects.Increment(r.Rejects)
	d.shortCircuits.Increment(r.ShortCircuits)
	d.timeouts.Increment(r.Timeouts)
	d.fallbackSuccesses.Increment(r.FallbackSuccesses)
	d.fallbackFailures.Increment(r.FallbackFailures)
	d.contextCanceled.Increment(r.ContextCanceled)
	d.contextDeadlineExceeded.Increment(r.ContextDeadlineExceeded)

	d.totalDuration.Add(r.TotalDuration)
	d.runDuration.Add(r.RunDuration)
}

// Reset resets all metrics in this collector to 0.
func (d *DefaultCollector) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.numRequests = rolling.NewNumber()
	d.errors = rolling.NewNumber()
	d.successes = rolling.NewNumber()
	d.rejects = rolling.NewNumber()
	d.shortCircuits = rolling.NewNumber()
	d.failures = rolling.NewNumber()
	d.timeouts = rolling.NewNumber()
	d.fallbackSuccesses = rolling.NewNumber()
	d.fallbackFailures = rolling.NewNumber()
	d.contextCanceled = rolling.NewNumber()
	d.contextDeadlineExceeded = rolling.NewNumber()
	d.totalDuration = rolling.NewTiming()
	d.runDuration = rolling.NewTiming()
}
