package managedcluster

type LoadBalancingRule struct {
	BackendPort      int64         `json:"backendPort"`
	FrontendPort     int64         `json:"frontendPort"`
	ProbePort        *int64        `json:"probePort,omitempty"`
	ProbeProtocol    ProbeProtocol `json:"probeProtocol"`
	ProbeRequestPath *string       `json:"probeRequestPath,omitempty"`
	Protocol         Protocol      `json:"protocol"`
}
