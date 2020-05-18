package rigis

import (
	"net/http"

	"github.com/rocinax/rigis/pkg/rule"
)

type node struct {
	rule                   rule.Rule
	filter                 filter
	balanceType            string
	backendHosts           []backendHost
	backendHostWeightIndex int
	backendHostIndex       int
}

func newNode(
	nrule rule.Rule,
	nfilter filter,
	balanceType string,
	backendHosts []backendHost,
) node {

	return node{
		rule:                   nrule,
		filter:                 nfilter,
		balanceType:            balanceType,
		backendHosts:           backendHosts,
		backendHostWeightIndex: 0,
		backendHostIndex:       0,
	}
}

func (n node) serveHTTP(rw http.ResponseWriter, req *http.Request) {
	backend := n.getBackendHost()
	backend.serveHTTP(rw, req)
}

func (n node) executeRule(req *http.Request) bool {
	return n.rule.Execute(req)
}

func (n node) executeFilter(req *http.Request) bool {
	return n.filter.Execute(req)
}

func (n node) getBackendHost() *backendHost {
	resultHost := n.backendHosts[n.backendHostIndex]
	n.backendHostWeightIndex++
	if n.backendHostWeightIndex >= n.backendHosts[n.backendHostIndex].weight {
		n.nextBackendHost()
	}
	return &resultHost
}

func (n node) nextBackendHost() {
	n.backendHostWeightIndex = 0
	n.backendHostIndex++
	if n.backendHostIndex >= len(n.backendHosts) {
		n.backendHostIndex = 0
	}
}

func (n node) getNextBackendHost() *backendHost {
	n.nextBackendHost()
	return n.getBackendHost()
}
