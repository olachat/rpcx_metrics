package prom

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gogf/gf/os/glog"
)

// TrackHandler monitor an event handler or business go method like an HTTP Request Handler, and protect it
// from panicking the whole application.
func (o *Registerer) TrackHandler(err *error, methodName string) (onCompleted func()) {
	timeStart := time.Now()

	o.requestsReceived.WithLabelValues(methodName).Inc()

	// onCompleted happen as a deferred func at handler ends
	onCompleted = func() {
		statusCode := 0

		// handle panic & map error to std http status code
		if p := recover(); p != nil {

			statusCode = CustomHTTPCodePanic
			glog.Errorf("TrackHandler: Panic in %s: %#v: stack:\n%s\n", methodName, p, string(debug.Stack()))

		} else if *err == nil {
			statusCode = http.StatusOK
		} else if invalidParam, ok := (*err).(*ErrorInvalidParam); ok {

			statusCode = http.StatusBadRequest
			glog.Errorf("TrackHandler: %s: InvalidArgument err: %v", methodName, invalidParam)

		} else {

			statusCode = http.StatusInternalServerError
			glog.Errorf("TrackHandler: %s: internal err: %v", methodName, *err)

		}

		status := fmt.Sprintf("%d", statusCode)
		o.responsesSent.WithLabelValues(methodName, status).Inc()

		duration := time.Since(timeStart).Seconds()
		o.rpcDurations.WithLabelValues(methodName, status).Observe(duration)
	}

	return
}

const (
	CustomHTTPCodePanic = 999
)
