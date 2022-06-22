package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataSourcesSpec struct {
	Extensions          *[]ExtensionDataSource       `json:"extensions,omitempty"`
	PerformanceCounters *[]PerfCounterDataSource     `json:"performanceCounters,omitempty"`
	Syslog              *[]SyslogDataSource          `json:"syslog,omitempty"`
	WindowsEventLogs    *[]WindowsEventLogDataSource `json:"windowsEventLogs,omitempty"`
}
