package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMonitorParameters struct {
	AutoStart                   *bool                                 `json:"autoStart,omitempty"`
	Destination                 *ConnectionMonitorDestination         `json:"destination,omitempty"`
	Endpoints                   *[]ConnectionMonitorEndpoint          `json:"endpoints,omitempty"`
	MonitoringIntervalInSeconds *int64                                `json:"monitoringIntervalInSeconds,omitempty"`
	Notes                       *string                               `json:"notes,omitempty"`
	Outputs                     *[]ConnectionMonitorOutput            `json:"outputs,omitempty"`
	Source                      *ConnectionMonitorSource              `json:"source,omitempty"`
	TestConfigurations          *[]ConnectionMonitorTestConfiguration `json:"testConfigurations,omitempty"`
	TestGroups                  *[]ConnectionMonitorTestGroup         `json:"testGroups,omitempty"`
}
