package rigis

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type backendHost struct {
	weight int
	url    *url.URL
	rp     *httputil.ReverseProxy
}

func newBackendHost(weight int, beURL *url.URL) backendHost {

	rp := newSingleHostReverseProxy(beURL)
	rp.ModifyResponse = modifyResponse
	rp.ErrorHandler = errorHandler

	return backendHost{
		weight: weight,
		url:    beURL,
		rp:     rp,
	}
}

func newSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

		// set X-Forwarded-Proto header
		if viper.GetBool("TLSEnable") {
			req.Header.Set("X-Forwarded-Proto", "https")
		} else {
			req.Header.Set("X-Forwarded-Proto", "http")
		}

		// set X-Forwarded-Host header
		req.Header.Set("X-Forwarded-Host", req.Host)

		// set X-Real-IP header
		if req.Header.Get("X-Real-IP") == "" {
			if req.Header.Get("X-Forwarded-For") != "" {
				req.Header.Set("X-Real-IP", strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0])
			} else {
				req.Header.Set("X-Real-IP", strings.Split(req.Host, ":")[0])
			}
		}

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func modifyResponse(res *http.Response) error {
	res.Header.Set("Server", viper.GetString("ServerName"))
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
