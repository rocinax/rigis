package rigis

// Config Rocinax Castrum Configuration
type Config struct {
	Filter ConfigFilter
	Nodes  []ConfigNode
}

// ConfigRule :
type ConfigRule struct {
	Type    string
	Setting map[string]interface{}
	Rules   []ConfigRule
}

// ConfigFilter :
type ConfigFilter struct {
	Rule   ConfigRule
	Accept bool
}

// ConfigNode :
type ConfigNode struct {
	Rule         ConfigRule
	Filter       ConfigFilter
	BalanceType  string
	BackendHosts []ConfigBackendHost
}

// ConfigBackendHost :
type ConfigBackendHost struct {
	Weight int
	URL    string
}

// Check Config Format Check
func (config Config) Check() bool {

	return true
}
