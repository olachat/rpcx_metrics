package prom

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gogf/gf/net/ghttp"
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

var reqStartTimestampKey = new(int)

func (o *Registerer) Middleware() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		// before request
		path := r.URL.Path
		o.requestsReceived.WithLabelValues(path).Inc()
		r.SetCtxVar(reqStartTimestampKey, time.Now().UnixMicro())

		r.Middleware.Next()

		//after request
		o.responsesSent.WithLabelValues(path, strconv.Itoa(r.Response.Status)).Inc()
		start := r.GetCtxVar(reqStartTimestampKey, 0).Int64()
		if start > 0 {
			o.rpcDurations.WithLabelValues(path, strconv.Itoa(r.Response.Status)).Observe(float64(time.Now().UnixMicro()-start) / 1_000_000)
		}
	}
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
