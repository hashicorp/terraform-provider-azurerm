package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataSourcesSpec struct {
	DataImports         *DataImportSources               `json:"dataImports,omitempty"`
	Extensions          *[]ExtensionDataSource           `json:"extensions,omitempty"`
	IisLogs             *[]IisLogsDataSource             `json:"iisLogs,omitempty"`
	LogFiles            *[]LogFilesDataSource            `json:"logFiles,omitempty"`
	PerformanceCounters *[]PerfCounterDataSource         `json:"performanceCounters,omitempty"`
	PlatformTelemetry   *[]PlatformTelemetryDataSource   `json:"platformTelemetry,omitempty"`
	PrometheusForwarder *[]PrometheusForwarderDataSource `json:"prometheusForwarder,omitempty"`
	Syslog              *[]SyslogDataSource              `json:"syslog,omitempty"`
	WindowsEventLogs    *[]WindowsEventLogDataSource     `json:"windowsEventLogs,omitempty"`
	WindowsFirewallLogs *[]WindowsFirewallLogsDataSource `json:"windowsFirewallLogs,omitempty"`
}
