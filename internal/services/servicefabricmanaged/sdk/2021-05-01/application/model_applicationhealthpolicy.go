package application

type ApplicationHealthPolicy struct {
	ConsiderWarningAsError                  bool                                `json:"considerWarningAsError"`
	DefaultServiceTypeHealthPolicy          *ServiceTypeHealthPolicy            `json:"defaultServiceTypeHealthPolicy,omitempty"`
	MaxPercentUnhealthyDeployedApplications int64                               `json:"maxPercentUnhealthyDeployedApplications"`
	ServiceTypeHealthPolicyMap              *map[string]ServiceTypeHealthPolicy `json:"serviceTypeHealthPolicyMap,omitempty"`
}
