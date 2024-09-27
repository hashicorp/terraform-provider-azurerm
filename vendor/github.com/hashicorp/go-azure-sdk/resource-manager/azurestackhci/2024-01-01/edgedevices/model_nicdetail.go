package edgedevices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NicDetail struct {
	AdapterName          string    `json:"adapterName"`
	ComponentId          *string   `json:"componentId,omitempty"`
	DefaultGateway       *string   `json:"defaultGateway,omitempty"`
	DefaultIsolationId   *string   `json:"defaultIsolationId,omitempty"`
	DnsServers           *[]string `json:"dnsServers,omitempty"`
	DriverVersion        *string   `json:"driverVersion,omitempty"`
	IP4Address           *string   `json:"ip4Address,omitempty"`
	InterfaceDescription *string   `json:"interfaceDescription,omitempty"`
	SubnetMask           *string   `json:"subnetMask,omitempty"`
}
