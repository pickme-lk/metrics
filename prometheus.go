package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type prometheusReporter struct {
	registry             prometheus.Registerer
	prefix               string
	namespace, subSystem string
	constLabels          map[string]string
	availableMetrics     map[string]Collector
	*sync.Mutex
}

func PrometheusReporter(conf ReporterConf) Reporter {
	constLabels := map[string]string{}
	if conf.ConstLabels != nil {
		for label, val := range conf.ConstLabels {
			constLabels[label] = val
		}
	}

	r := &prometheusReporter{
		registry:    prometheus.DefaultRegisterer,
		namespace:   conf.System,
		subSystem:   conf.Subsystem,
		constLabels: mergeLabels(constLabels, nil),
		Mutex:       new(sync.Mutex),
	}
	r.availableMetrics = make(map[string]Collector)

	return r
}

func (r *prometheusReporter) Reporter(conf ReporterConf) Reporter {
	rConf := ReporterConf{
		System:      r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: r.constLabels,
	}
	if conf.Subsystem != `` {
		rConf.Subsystem = fmt.Sprintf(`%s_%s`, r.subSystem, conf.Subsystem)
	}

	return PrometheusReporter(rConf)
}

func (r *prometheusReporter) Counter(conf MetricConf) Counter {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if c, ok := r.availableMetrics[conf.Path]; ok {
		return c.(Counter)
	}

	promCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: mergeLabels(r.constLabels, conf.ConstLabels),
	}, conf.Labels)

	if err := r.registry.Register(promCounter); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	c := &prometheusCounter{
		counter: promCounter,
	}

	r.availableMetrics[conf.Path] = c

	return c
}

func (r *prometheusReporter) Gauge(conf MetricConf) Gauge {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if g, ok := r.availableMetrics[conf.Path]; ok {
		return g.(Gauge)
	}

	promGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: mergeLabels(r.constLabels, conf.ConstLabels),
	}, conf.Labels)

	if err := r.registry.Register(promGauge); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	g := &prometheusGauge{
		gauge: promGauge,
	}

	r.availableMetrics[conf.Path] = g

	return g
}

func (r *prometheusReporter) GaugeFunc(conf MetricConf, f func() float64) GaugeFunc {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	promGauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: mergeLabels(r.constLabels, conf.ConstLabels),
	}, f)

	if err := r.registry.Register(promGauge); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	g := &prometheusGaugeFunc{
		gauge: promGauge,
	}

	r.availableMetrics[conf.Path] = g

	return g

}

func (r *prometheusReporter) Observer(conf MetricConf) Observer {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics[conf.Path]; ok {
		return o.(Observer)
	}

	promObserver := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: mergeLabels(r.constLabels, conf.ConstLabels),
	}, conf.Labels)

	if err := r.registry.Register(promObserver); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	h := &prometheusSummary{
		observer: promObserver,
	}

	r.availableMetrics[conf.Path] = h

	return h
}

func (r *prometheusReporter) Summary(conf MetricConf) Observer {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics[conf.Path]; ok {
		return o.(Observer)
	}

	promObserver := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: mergeLabels(r.constLabels, conf.ConstLabels),
	}, conf.Labels)

	if err := r.registry.Register(promObserver); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	h := &prometheusSummary{
		observer: promObserver,
	}

	r.availableMetrics[conf.Path] = h

	return h
}

func (r *prometheusReporter) Histogram(conf MetricConf) Observer {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics[conf.Path]; ok {
		return o.(Observer)
	}

	promObserver := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        conf.Path,
		Help:        conf.Path,
		Namespace:   r.namespace,
		Subsystem:   r.subSystem,
		ConstLabels: mergeLabels(r.constLabels, conf.ConstLabels),
	}, conf.Labels)

	if err := r.registry.Register(promObserver); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			panic(err)
		}
	}

	h := &prometheusHistogram{
		observer: promObserver,
	}

	r.availableMetrics[conf.Path] = h

	return h
}

func (r *prometheusReporter) Info() string { return `` }

func (r *prometheusReporter) UnRegister(metrics string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if o, ok := r.availableMetrics[metrics]; ok {
		o.UnRegister()
		delete(r.availableMetrics, metrics)
	}
}

type prometheusCounter struct {
	counter *prometheus.CounterVec
}

func (c *prometheusCounter) Count(value float64, lbs map[string]string) {
	c.counter.With(lbs).Add(value)
}

func (c *prometheusCounter) UnRegister() {
	prometheus.Unregister(c.counter)
}

type prometheusGauge struct {
	gauge *prometheus.GaugeVec
}

func (g *prometheusGauge) Count(value float64, lbs map[string]string) {
	g.gauge.With(lbs).Set(value)
}

func (g *prometheusGauge) UnRegister() {
	prometheus.Unregister(g.gauge)
}

type prometheusGaugeFunc struct {
	gauge prometheus.GaugeFunc
}

func (h *prometheusGaugeFunc) UnRegister() {
	prometheus.Unregister(h.gauge)
}

type prometheusHistogram struct {
	observer *prometheus.HistogramVec
}

func (c *prometheusHistogram) Observe(value float64, lbs map[string]string) {
	c.observer.With(lbs).Observe(value)
}

func (c *prometheusHistogram) UnRegister() {
	prometheus.Unregister(c.observer)
}

type prometheusSummary struct {
	observer *prometheus.SummaryVec
}

func (c *prometheusSummary) Observe(value float64, lbs map[string]string) {
	c.observer.With(lbs).Observe(value)
}

func (c *prometheusSummary) UnRegister() {
	prometheus.Unregister(c.observer)
}

func mergeLabels(from map[string]string, to map[string]string) map[string]string {
	constLabels := map[string]string{}
	// get existing labels
	if from != nil {
		for label, val := range from {
			constLabels[label] = val
		}
	}

	if to != nil {
		for label, val := range to {
			if _, ok := constLabels[label]; ok {
				panic(fmt.Sprintf(`label %s already registred`, label))
			}
			constLabels[label] = val
		}
	}

	return constLabels
}
