package metrics

type MetricConf struct {
	Path        string
	Labels      []string
	ConstLabels map[string]string
}

type Reporter interface {
	//Reporter(labels []string) Reporter
	Counter(MetricConf) Counter
	Observer(MetricConf) Observer
	Gauge(MetricConf) Gauge
	Info() string
}

type Counter interface {
	Count(value float64, lbs map[string]string)
}

type Gauge interface {
	Count(value float64, lbs map[string]string)
}

type Observer interface {
	Observe(value float64, lbs map[string]string)
}
