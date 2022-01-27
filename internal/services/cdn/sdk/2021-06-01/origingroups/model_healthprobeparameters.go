package origingroups

type HealthProbeParameters struct {
	ProbeIntervalInSeconds *int64                  `json:"probeIntervalInSeconds,omitempty"`
	ProbePath              *string                 `json:"probePath,omitempty"`
	ProbeProtocol          *ProbeProtocol          `json:"probeProtocol,omitempty"`
	ProbeRequestType       *HealthProbeRequestType `json:"probeRequestType,omitempty"`
}
