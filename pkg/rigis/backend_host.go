package rigis

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sirupsen/logrus"
)

type backendHost struct {
	weight int
	url    *url.URL
	rp     *httputil.ReverseProxy
}

func newBackendHost(weight int, beURL *url.URL) backendHost {

	rp := httputil.NewSingleHostReverseProxy(beURL)
	rp.ModifyResponse = outputAccessLog
	rp.ErrorHandler = errorHandler

	return backendHost{
		weight: weight,
		url:    beURL,
		rp:     rp,
	}
}

func outputAccessLog(res *http.Response) error {
	logrus.WithFields(FormatAccessLog(res)).Info("write access log at information level")
	return nil
}

func errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	logrus.WithFields(FormatErrorLog(req)).Error("backend server is dead")
	ResponseError(rw, req, map[string]interface{}{
		"Status":           http.StatusServiceUnavailable,
		"Error":            "Service Unavailable",
		"ErrorDescription": "HTTP Service Unavailable.",
	})
}
