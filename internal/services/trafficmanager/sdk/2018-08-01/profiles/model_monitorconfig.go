package profiles

type MonitorConfig struct {
	CustomHeaders             *[]MonitorConfigCustomHeadersInlined            `json:"customHeaders,omitempty"`
	ExpectedStatusCodeRanges  *[]MonitorConfigExpectedStatusCodeRangesInlined `json:"expectedStatusCodeRanges,omitempty"`
	IntervalInSeconds         *int64                                          `json:"intervalInSeconds,omitempty"`
	Path                      *string                                         `json:"path,omitempty"`
	Port                      *int64                                          `json:"port,omitempty"`
	ProfileMonitorStatus      *ProfileMonitorStatus                           `json:"profileMonitorStatus,omitempty"`
	Protocol                  *MonitorProtocol                                `json:"protocol,omitempty"`
	TimeoutInSeconds          *int64                                          `json:"timeoutInSeconds,omitempty"`
	ToleratedNumberOfFailures *int64                                          `json:"toleratedNumberOfFailures,omitempty"`
}
