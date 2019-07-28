package metrics

type noopReporter struct{}

func NoopReporter() Reporter {
	return noopReporter{}
}

func (noopReporter) Reporter(ReporterConf) Reporter { return noopReporter{} }

func (noopReporter) Counter(MetricConf) Counter { return noopCounter{} }

func (noopReporter) Observer(MetricConf) Observer { return noopObserver{} }

func (noopReporter) Gauge(MetricConf) Gauge { return noopGauge{} }

func (noopReporter) GaugeFunc(MetricConf, func() float64) GaugeFunc { return noopGaugeFunc{} }

func (noopReporter) Info() string { return `` }

func (noopReporter) UnRegister(metrics string) {}

type noopCounter struct{}

func (noopCounter) Count(value float64, lbs map[string]string) {}
func (noopCounter) UnRegister()                                {}

type noopGauge struct{}

func (noopGauge) Count(value float64, lbs map[string]string) {}
func (noopGauge) UnRegister()                                {}

type noopGaugeFunc struct{}

func (noopGaugeFunc) UnRegister() {}

type noopObserver struct{}

func (noopObserver) Observe(value float64, lbs map[string]string) {}
func (noopObserver) UnRegister()                                  {}
