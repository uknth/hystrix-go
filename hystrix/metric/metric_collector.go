package metric

import (
	"sync"
	"time"
)

// Registry is the default metricCollectorRegistry that circuits will use to
// collect statistics about the health of the circuit.
var Registry = metricCollectorRegistry{
	lock: &sync.RWMutex{},
	registry: []func(name string) Collector{
		newDefaultCollector,
	},
}

type metricCollectorRegistry struct {
	lock     *sync.RWMutex
	registry []func(name string) Collector
}

// InitializeCollectors runs the registried Collector Initializers to create an array of Collectors.
func (m *metricCollectorRegistry) InitializeCollectors(name string) []Collector {
	m.lock.RLock()
	defer m.lock.RUnlock()

	metrics := make([]Collector, len(m.registry))
	for i, metricCollectorInitializer := range m.registry {
		metrics[i] = metricCollectorInitializer(name)
	}
	return metrics
}

// Register places a Collector Initializer in the registry maintained by this metricCollectorRegistry.
func (m *metricCollectorRegistry) Register(initCollector func(string) Collector) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.registry = append(m.registry, initCollector)
}

type Result struct {
	Attempts                float64
	Errors                  float64
	Successes               float64
	Failures                float64
	Rejects                 float64
	ShortCircuits           float64
	Timeouts                float64
	FallbackSuccesses       float64
	FallbackFailures        float64
	ContextCanceled         float64
	ContextDeadlineExceeded float64
	TotalDuration           time.Duration
	RunDuration             time.Duration
	ConcurrencyInUse        float64
}

// Collector represents the contract that all collectors must fulfill to gather circuit statistics.
// Implementations of this interface do not have to maintain locking around thier data stores so long as
// they are not modified outside of the hystrix context.
type Collector interface {
	// Update accepts a set of metrics from a command execution for remote instrumentation
	Update(Result)
	// Reset resets the internal counters and timers.
	Reset()
}
