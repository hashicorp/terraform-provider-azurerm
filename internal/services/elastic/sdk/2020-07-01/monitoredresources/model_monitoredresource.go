package monitoredresources

type MonitoredResource struct {
	Id                  *string      `json:"id,omitempty"`
	ReasonForLogsStatus *string      `json:"reasonForLogsStatus,omitempty"`
	SendingLogs         *SendingLogs `json:"sendingLogs,omitempty"`
}
