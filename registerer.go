package prom

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Registerer struct {
	registry            *prometheus.Registry
	requestsReceived    *prometheus.CounterVec
	responsesSent       *prometheus.CounterVec
	rpcDurations        *prometheus.SummaryVec
	dataQuality         *prometheus.CounterVec
	gauge               *prometheus.GaugeVec
	apiSignatureCounter *prometheus.CounterVec
}

func (o *Registerer) GetMetricAPIHandler() (h http.Handler) {
	h = promhttp.InstrumentMetricHandler(
		o.registry, promhttp.HandlerFor(o.registry, promhttp.HandlerOpts{}),
	)
	return
}

func (o *Registerer) TrackRequestReceived(method string) {
	o.requestsReceived.WithLabelValues(method).Inc()
}

func (o *Registerer) TrackResponseSent(method, status string) {
	o.responsesSent.WithLabelValues(method, status).Inc()
}

func (o *Registerer) TrackRpcDuration(method, status string, dur time.Duration) {
	o.rpcDurations.WithLabelValues(method, status).Observe(dur.Seconds())
}
