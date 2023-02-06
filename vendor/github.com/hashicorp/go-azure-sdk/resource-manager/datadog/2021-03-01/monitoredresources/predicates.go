package monitoredresources

type MonitoredResourceOperationPredicate struct {
	Id                     *string
	ReasonForLogsStatus    *string
	ReasonForMetricsStatus *string
	SendingLogs            *bool
	SendingMetrics         *bool
}

func (p MonitoredResourceOperationPredicate) Matches(input MonitoredResource) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.ReasonForLogsStatus != nil && (input.ReasonForLogsStatus == nil && *p.ReasonForLogsStatus != *input.ReasonForLogsStatus) {
		return false
	}

	if p.ReasonForMetricsStatus != nil && (input.ReasonForMetricsStatus == nil && *p.ReasonForMetricsStatus != *input.ReasonForMetricsStatus) {
		return false
	}

	if p.SendingLogs != nil && (input.SendingLogs == nil && *p.SendingLogs != *input.SendingLogs) {
		return false
	}

	if p.SendingMetrics != nil && (input.SendingMetrics == nil && *p.SendingMetrics != *input.SendingMetrics) {
		return false
	}

	return true
}
