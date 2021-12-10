package servers

type IPv4FirewallSettings struct {
	EnablePowerBIService *bool               `json:"enablePowerBIService,omitempty"`
	FirewallRules        *[]IPv4FirewallRule `json:"firewallRules,omitempty"`
}
