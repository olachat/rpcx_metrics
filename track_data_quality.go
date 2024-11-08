package prom

func (o *Registerer) TrackDataQuality(queryName string, inputGroup string, outputGroup string) {
	o.dataQuality.WithLabelValues(queryName, inputGroup, outputGroup).Inc()
	return
}

func (o *Registerer) TrackDataQualityScore(queryName string, inputGroup string, outputGroup string, floatIncrement float64) {
	o.dataQuality.WithLabelValues(queryName, inputGroup, outputGroup).Add(floatIncrement)
	return
}

func (o *Registerer) TrackCount(countName, inputGroup string, count float64) {
	o.gauge.WithLabelValues(countName, inputGroup).Set(count)
}

func (o *Registerer) IncCount(gaugeType string, name string, key string) {
	o.gauge.WithLabelValues(gaugeType, name, key).Inc()
}

func (o *Registerer) DecCount(gaugeType string, name string, key string) {
	o.gauge.WithLabelValues(gaugeType, name, key).Dec()
}

func (o *Registerer) SetCount(gaugeType string, name string, key string, count int) {
	o.gauge.WithLabelValues(gaugeType, name, key).Set(float64(count))
}
