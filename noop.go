package metrics

type noopReporter struct{}

func NoopReporter() Reporter {
	return noopReporter{}
}

func (noopReporter) Counter(path string, labels []string) Counter {
	return noopCounter{}
}

func (noopReporter) Observer(path string, labels []string) Observer {
	return noopObserver{}
}

func (noopReporter) Gauge(path string, labels []string) Gauge {
	return noopGauge{}
}

func (noopReporter) Info() string { return `` }

type noopCounter struct {
}

func (noopCounter) Count(value float64, lbs map[string]string) {
	return
}

type noopGauge struct {
}

func (noopGauge) Count(value float64, lbs map[string]string) {
	return
}

type noopObserver struct{}

func (noopObserver) Observe(value float64, lbs map[string]string) {
	return
}
