package metrics

type MetricConf struct {
	Path        string
	Labels      []string
	ConstLabels map[string]string
}

type ReporterConf struct {
	System      string
	Subsystem   string
	ConstLabels map[string]string
}

type Reporter interface {
	Reporter(ReporterConf) Reporter
	Counter(MetricConf) Counter
	Observer(MetricConf) Observer
	Gauge(MetricConf) Gauge
	GaugeFunc(MetricConf, func() float64) GaugeFunc
	Info() string
	UnRegister(metrics string)
}

type Collector interface {
	UnRegister()
}

type Counter interface {
	Collector
	Count(value float64, lbs map[string]string)
}

type Gauge interface {
	Collector
	Count(value float64, lbs map[string]string)
}

type GaugeFunc interface {
	Collector
}

type Observer interface {
	Collector
	Observe(value float64, lbs map[string]string)
}
