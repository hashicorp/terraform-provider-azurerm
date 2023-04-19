package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerServiceNetworkProfileKubeProxyConfigIPvsConfig struct {
	Scheduler            *IPvsScheduler `json:"scheduler,omitempty"`
	TcpFinTimeoutSeconds *int64         `json:"tcpFinTimeoutSeconds,omitempty"`
	TcpTimeoutSeconds    *int64         `json:"tcpTimeoutSeconds,omitempty"`
	UdpTimeoutSeconds    *int64         `json:"udpTimeoutSeconds,omitempty"`
}
