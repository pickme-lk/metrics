package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type prometheusReporter struct {
	registry             prometheus.Registerer
	namespace, subSystem string
	availableMetrics     struct {
		counters  map[string]Counter
		histos    map[string]Observer
		summaries map[string]Observer
		gauges    map[string]Gauge
	}
	*sync.Mutex
}

func PrometheusReporter(system string, subsystem string) Reporter {
	r := &prometheusReporter{
		registry:  prometheus.DefaultRegisterer,
		namespace: system,
		subSystem: subsystem,
		Mutex:     new(sync.Mutex),
	}

	r.availableMetrics.counters = make(map[string]Counter)
	r.availableMetrics.gauges = make(map[string]Gauge)
	r.availableMetrics.summaries = make(map[string]Observer)
	r.availableMetrics.histos = make(map[string]Observer)

	return r
}

func (r *prometheusReporter) Counter(path string, labels []string) Counter {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if c, ok := r.availableMetrics.counters[path]; ok {
		return c
	}

	promCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      path,
		Help:      path,
		Namespace: r.namespace,
		Subsystem: r.subSystem,
	}, labels)

	if err := r.registry.Register(promCounter); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	c := &prometheusCounter{
		counter: promCounter,
	}

	r.availableMetrics.counters[path] = c

	return c
}

func (r *prometheusReporter) Gauge(path string, labels []string) Gauge {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if g, ok := r.availableMetrics.gauges[path]; ok {
		return g
	}

	promGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      path,
		Help:      path,
		Namespace: r.namespace,
		Subsystem: r.subSystem,
	}, labels)

	if err := r.registry.Register(promGauge); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	g := &prometheusGauge{
		gauge: promGauge,
	}

	r.availableMetrics.gauges[path] = g

	return g

}

func (r *prometheusReporter) Observer(path string, labels []string) Observer {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics.summaries[path]; ok {
		return o
	}

	promObserver := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      path,
		Help:      path,
		Namespace: r.namespace,
		Subsystem: r.subSystem,
	}, labels)

	if err := r.registry.Register(promObserver); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	h := &prometheusSummary{
		observer: promObserver,
	}

	r.availableMetrics.summaries[path] = h

	return h
}

func (r *prometheusReporter) Summary(path string, labels []string) Observer {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics.summaries[path]; ok {
		return o
	}

	promObserver := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      path,
		Help:      path,
		Namespace: r.namespace,
		Subsystem: r.subSystem,
	}, labels)

	if err := r.registry.Register(promObserver); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	h := &prometheusSummary{
		observer: promObserver,
	}

	r.availableMetrics.summaries[path] = h

	return h
}

func (r *prometheusReporter) Histogram(path string, labels []string) Observer {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics.histos[path]; ok {
		return o
	}

	promObserver := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:      path,
		Help:      path,
		Namespace: r.namespace,
		Subsystem: r.subSystem,
	}, labels)

	if err := r.registry.Register(promObserver); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	h := &prometheusHistogram{
		observer: promObserver,
	}

	r.availableMetrics.histos[path] = h

	return h
}

func (r *prometheusReporter) Info() string { return `` }

type prometheusCounter struct {
	counter *prometheus.CounterVec
}

func (c *prometheusCounter) Count(value float64, lbs map[string]string) {
	c.counter.With(lbs).Add(value)
}

type prometheusGauge struct {
	gauge *prometheus.GaugeVec
}

func (g *prometheusGauge) Count(value float64, lbs map[string]string) {
	g.gauge.With(lbs).Set(value)
}

type prometheusHistogram struct {
	observer *prometheus.HistogramVec
}

func (c *prometheusHistogram) Observe(value float64, lbs map[string]string) {
	c.observer.With(lbs).Observe(value)
}

type prometheusSummary struct {
	observer *prometheus.SummaryVec
}

func (c *prometheusSummary) Observe(value float64, lbs map[string]string) {
	c.observer.With(lbs).Observe(value)
}
