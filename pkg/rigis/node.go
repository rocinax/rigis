package rigis

import (
	"net/http"
	"net/http/httputil"

	"github.com/rocinax/rigis/pkg/rule"
)

type node struct {
	rule             rule.Rule
	filter           filter
	balanceType      string
	backendHosts     []backendHost
	backendProxyList []*httputil.ReverseProxy
	backendIndex     int
}

func newNode(
	nrule rule.Rule,
	nfilter filter,
	balanceType string,
	backendHosts []backendHost,
) node {

	return node{
		rule:             nrule,
		filter:           nfilter,
		balanceType:      balanceType,
		backendHosts:     backendHosts,
		backendProxyList: newBackendProxyList(backendHosts),
		backendIndex:     0,
	}
}

func newBackendProxyList(backendHosts []backendHost) []*httputil.ReverseProxy {
	var backendProxyList []*httputil.ReverseProxy

	for i := 0; i < len(backendHosts); i++ {
		for j := 0; j < backendHosts[i].weight; j++ {
			backendProxyList = append(backendProxyList, backendHosts[i].rp)
		}
	}
	return backendProxyList
}

func (n node) executeRule(req *http.Request) bool {
	return n.rule.Execute(req)
}

func (n node) executeFilter(req *http.Request) bool {
	return n.filter.Execute(req)
}

func (n node) getBackend() *httputil.ReverseProxy {

	if n.backendIndex >= len(n.backendProxyList) {
		n.backendIndex = 0
	}

	resultRP := n.backendProxyList[n.backendIndex]
	n.backendIndex++

	return resultRP
}
