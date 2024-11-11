package prom

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/olachat/rpcx_metrics/discovery"
	"github.com/olachat/rpcx_metrics/tool"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
)

// PrometheusHandler  定义了有个RPC server插件
type PrometheusHandler struct {
	ctx            context.Context
	cancelFunc     context.CancelFunc
	registerer     *Registerer
	consulAddr     string
	port           int
	wg             sync.WaitGroup
	consulDiscover *discovery.ConsulDiscovery
	registerServer *http.Server
}

// NewPrometheusPlugin  实例化一个 InfoHandler
func NewPrometheusPlugin(consulAddr string) *PrometheusHandler {
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

// Register 注册插件
func (h *PrometheusHandler) Register(name string, _ interface{}, metadata string) error {
	availablePort := tool.GetAvailablePort()
	log.Printf("serviceName: %+v prometheus port: %+v start", name, availablePort)

	h.consulDiscover = discovery.NewConsulDiscovery(h.consulAddr, name, availablePort)

	http.Handle("/metric", DefaultRegisterer().GetMetricAPIHandler())
	// 增加健康检测
	http.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("pong"))
	})

	h.registerServer = &http.Server{Addr: fmt.Sprintf(":%d", availablePort)}

	go func() {
		// 向consul进行注册
		if err := h.consulDiscover.Register([]string{name}, map[string]string{
			"service_name": name,
		}); err != nil {
			panic(err)
		}
		if err := h.registerServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("prometheus ListenAndServer failed: %v", err))
		}
	}()
	return nil
}

// Unregister 解除注册
func (h *PrometheusHandler) Unregister(_ string) error {
	_ = h.consulDiscover.Deregister()
	if err := h.registerServer.Shutdown(context.Background()); err != nil {
		log.Printf("prometheus Shutdown failed err: %+v", err)
	}
	log.Printf("prometheus shutdown finished")
	return nil
}

// Use a regular expression to match non-alphabetic characters
var reg = regexp.MustCompile("[^a-zA-Z-]")

func SanitizeString(input string) string {
	// Replace non-alphabetic characters with an empty string
	sanitized := reg.ReplaceAllString(input, "_")

	return sanitized
}
