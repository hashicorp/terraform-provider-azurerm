package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMonitorTestConfiguration struct {
	HTTPConfiguration  *ConnectionMonitorHTTPConfiguration        `json:"httpConfiguration,omitempty"`
	IcmpConfiguration  *ConnectionMonitorIcmpConfiguration        `json:"icmpConfiguration,omitempty"`
	Name               string                                     `json:"name"`
	PreferredIPVersion *PreferredIPVersion                        `json:"preferredIPVersion,omitempty"`
	Protocol           ConnectionMonitorTestConfigurationProtocol `json:"protocol"`
	SuccessThreshold   *ConnectionMonitorSuccessThreshold         `json:"successThreshold,omitempty"`
	TcpConfiguration   *ConnectionMonitorTcpConfiguration         `json:"tcpConfiguration,omitempty"`
	TestFrequencySec   *int64                                     `json:"testFrequencySec,omitempty"`
}
