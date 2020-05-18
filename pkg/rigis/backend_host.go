package rigis

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type backendHost struct {
	weight int
	url    *url.URL
	rp     *httputil.ReverseProxy
}

func newBackendHost(weight int, beURL *url.URL) backendHost {

	bh := backendHost{
		weight: weight,
		url:    beURL,
		rp:     &httputil.ReverseProxy{},
	}

	bh.rp.Director = bh.director
	bh.rp.ModifyResponse = bh.modifyResponse
	bh.rp.ErrorHandler = bh.errorHandler

	return bh
}

func (bh backendHost) serveHTTP(rw http.ResponseWriter, req *http.Request) {
	bh.rp.ServeHTTP(rw, req)
}

func (bh backendHost) director(req *http.Request) {

	// Set Reverse Proxy Params
	req.URL.Scheme = bh.url.Scheme
	req.URL.Host = bh.url.Host
	req.URL.Path = normalizeRpPath(bh.url.Path, req.URL.Path)
	req.URL.RawQuery = normalizeRawQuery(bh.url.RawQuery, req.URL.RawQuery)

	// Set Optional HTTP Request Headers
	req.Header.Set("User-Agent", getUserAgent(req.Header.Get("User-Agent")))
	req.Header.Set("X-Forwarded-Proto", getXForwardedProto(req.Header.Get("X-Forwarded-Proto"), viper.GetBool("TLSEnable")))
	req.Header.Set("X-Forwarded-Host", getXForwardedHost(req.Header.Get("X-Forwarded-Host"), req.Host))
	req.Header.Set("X-Real-IP", getRemoteAddr(req))
}

func (bh backendHost) modifyResponse(res *http.Response) error {
	res.Header.Set("Server", viper.GetString("ServerName"))
	logrus.WithFields(formatAccessLog(res)).Info("rigis access log")
	return nil
}

func (bh backendHost) errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	logrus.WithFields(formatErrorLog(req)).Error("backend server is dead")
	responseError(rw, req, getErrorEntity(http.StatusServiceUnavailable))
}
