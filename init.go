package prom

import "github.com/prometheus/client_golang/prometheus"

var (
	defaultPrometheusRegistry *prometheus.Registry
	defaultRegisterer         *Registerer
)

func DefaultRegistry() *prometheus.Registry {
	return defaultPrometheusRegistry
}

func DefaultRegisterer() *Registerer {
	return defaultRegisterer
}

func init() {
	requestsReceived := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "req_total",
			Help: "Number of RPC requests received.",
		},
		[]string{"method"},
	)

	responsesSent := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "resp_total",
			Help: "Number of RPC responses sent.",
		},
		[]string{"method", "status"},
	)

	rpcDurations := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "resp_time_sec",
			Help:       "RPC latency distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "status"},
	)

	dataQuality := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "data_quality",
			Help: "Count of query result grouped by its quality level.",
		},
		[]string{"query", "input_group", "output_group"},
	)

	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gauge",
			Help: "record count at current time.",
		},
		[]string{"count", "input_group"},
	)

	apiSignatureCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "signature_version",
			Help: "Count check signature results grouped by version, validity and path.",
		},
		[]string{"version", "validity", "path"},
	)

	// don't registry := prometheus.NewRegistry()
	// use prometheus default registry to inherit existing go process metrics
	registry := prometheus.DefaultRegisterer.(*prometheus.Registry)
	registry.MustRegister(requestsReceived)
	registry.MustRegister(responsesSent)
	registry.MustRegister(rpcDurations)
	registry.MustRegister(dataQuality)
	registry.MustRegister(gauge)
	registry.MustRegister(apiSignatureCounter)

	defaultPrometheusRegistry = registry

	defaultRegisterer = &Registerer{
		registry:            registry,
		requestsReceived:    requestsReceived,
		responsesSent:       responsesSent,
		rpcDurations:        rpcDurations,
		dataQuality:         dataQuality,
		gauge:               gauge,
		apiSignatureCounter: apiSignatureCounter,
	}
}

func TrackHandler(err *error, methodName string) (onCompleted func()) {
	return defaultRegisterer.TrackHandler(err, methodName)
}

func TrackDataQuality(queryName string, inputGroup string, outputGroup string) {
	defaultRegisterer.TrackDataQuality(queryName, inputGroup, outputGroup)
}

func TrackDataQualityScore(queryName string, inputGroup string, outputGroup string, floatIncrement float64) {
	defaultRegisterer.TrackDataQualityScore(queryName, inputGroup, outputGroup, floatIncrement)
}

func TrackCount(countName, inputGroup string, floatIncrement float64) {
	defaultRegisterer.TrackCount(countName, inputGroup, floatIncrement)
}

func TrackAPISignature(version uint64, validType APISignatureType, path string) {
	defaultRegisterer.TrackAPISignature(version, validType, path)
}
