package metrics

type Reporter interface {
	Counter(path string, labels []string) Counter
	Observer(path string, labels []string) Observer
	Gauge(path string, labels []string) Gauge
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
