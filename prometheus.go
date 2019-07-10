package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type prometheusReporter struct {
	registry             prometheus.Registerer
	prefix               string
	namespace, subSystem string
	constLabels          map[string]string
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
		registry:    prometheus.DefaultRegisterer,
		namespace:   system,
		subSystem:   subsystem,
		constLabels: map[string]string{},
		Mutex:       new(sync.Mutex),
	}

	r.availableMetrics.counters = make(map[string]Counter)
	r.availableMetrics.gauges = make(map[string]Gauge)
	r.availableMetrics.summaries = make(map[string]Observer)
	r.availableMetrics.histos = make(map[string]Observer)

	return r
}

func (r *prometheusReporter) Reporter(labels []string) Reporter {
	return nil
}

func (r *prometheusReporter) Counter(conf MetricConf) Counter {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if c, ok := r.availableMetrics.counters[conf.Path]; ok {
		return c
	}

	promCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: r.constLabels,
	}, conf.Labels)

	if err := r.registry.Register(promCounter); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	c := &prometheusCounter{
		counter: promCounter,
	}

	r.availableMetrics.counters[conf.Path] = c

	return c
}

func (r *prometheusReporter) Gauge(conf MetricConf) Gauge {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if g, ok := r.availableMetrics.gauges[conf.Path]; ok {
		return g
	}

	promGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: r.constLabels,
	}, conf.Labels)

	if err := r.registry.Register(promGauge); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	g := &prometheusGauge{
		gauge: promGauge,
	}

	r.availableMetrics.gauges[conf.Path] = g

	return g

}

func (r *prometheusReporter) Observer(conf MetricConf) Observer {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics.summaries[conf.Path]; ok {
		return o
	}

	promObserver := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: r.constLabels,
	}, conf.Labels)

	if err := r.registry.Register(promObserver); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	h := &prometheusSummary{
		observer: promObserver,
	}

	r.availableMetrics.summaries[conf.Path] = h

	return h
}

func (r *prometheusReporter) Summary(conf MetricConf) Observer {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics.summaries[conf.Path]; ok {
		return o
	}

	promObserver := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: r.constLabels,
	}, conf.Labels)

	if err := r.registry.Register(promObserver); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	h := &prometheusSummary{
		observer: promObserver,
	}

	r.availableMetrics.summaries[conf.Path] = h

	return h
}

func (r *prometheusReporter) Histogram(conf MetricConf) Observer {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics.histos[conf.Path]; ok {
		return o
	}

	promObserver := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: r.constLabels,
	}, conf.Labels)

	if err := r.registry.Register(promObserver); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	h := &prometheusHistogram{
		observer: promObserver,
	}

	r.availableMetrics.histos[conf.Path] = h

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
