package prom

import (
	"context"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"regexp"
	"time"
)

// PrometheusHandler  定义了有个RPC server插件
type PrometheusHandler struct {
	registerer *Registerer
}

// NewPrometheusPlugin  实例化一个 InfoHandler
func NewPrometheusPlugin() *PrometheusHandler {
	return &PrometheusHandler{
		registerer: DefaultRegisterer(),
	}
}

func (h *PrometheusHandler) PostReadRequest(ctx context.Context, r *protocol.Message, e error) error {
	path := r.ServicePath
	method := r.ServiceMethod
	method = SanitizeString(path + "-" + method)
	if method == "" {
		return nil
	}

	h.registerer.TrackRequestReceived(method)
	return nil
}

func (h *PrometheusHandler) PostWriteResponse(ctx context.Context, req *protocol.Message, resp *protocol.Message, err error) error {
	path := req.ServicePath
	method := req.ServiceMethod
	method = SanitizeString(path + "-" + method)
	if method == "" {
		return nil
	}

	status := "normal"
	if resp.MessageStatusType() == protocol.Error {
		status = "error"
	}
	h.registerer.TrackResponseSent(method, status)

	t := ctx.Value(server.StartRequestContextKey).(int64)
	if t > 0 {
		t = time.Now().UnixNano() - t
		h.registerer.TrackRpcDuration(method, status, time.Duration(t)*time.Nanosecond)
	}

	return nil
}

// Use a regular expression to match non-alphabetic characters
var reg = regexp.MustCompile("[^a-zA-Z-]")

func SanitizeString(input string) string {
	// Replace non-alphabetic characters with an empty string
	sanitized := reg.ReplaceAllString(input, "_")

	return sanitized
}
