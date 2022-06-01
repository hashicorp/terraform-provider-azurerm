package rules

type MonitoringTagRulesProperties struct {
	LogRules          *LogRules          `json:"logRules,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
