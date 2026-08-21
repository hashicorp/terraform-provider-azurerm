package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityInformation struct {
	AvgLatencyInMs   *int64             `json:"avgLatencyInMs,omitempty"`
	ConnectionStatus *ConnectionStatus  `json:"connectionStatus,omitempty"`
	Hops             *[]ConnectivityHop `json:"hops,omitempty"`
	MaxLatencyInMs   *int64             `json:"maxLatencyInMs,omitempty"`
	MinLatencyInMs   *int64             `json:"minLatencyInMs,omitempty"`
	ProbesFailed     *int64             `json:"probesFailed,omitempty"`
	ProbesSent       *int64             `json:"probesSent,omitempty"`
}
