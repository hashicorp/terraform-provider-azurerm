package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogSettings struct {
	ApplicationInsights   *ApplicationInsights `json:"applicationInsights,omitempty"`
	CommonDestination     *LogDestination      `json:"commonDestination,omitempty"`
	DecryptLogDestination *LogDestination      `json:"decryptLogDestination,omitempty"`
	LogOption             *LogOption           `json:"logOption,omitempty"`
	LogType               *LogType             `json:"logType,omitempty"`
	ThreatLogDestination  *LogDestination      `json:"threatLogDestination,omitempty"`
	TrafficLogDestination *LogDestination      `json:"trafficLogDestination,omitempty"`
}
