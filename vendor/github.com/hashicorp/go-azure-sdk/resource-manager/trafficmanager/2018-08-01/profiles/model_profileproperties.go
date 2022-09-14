package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProfileProperties struct {
	AllowedEndpointRecordTypes  *[]AllowedEndpointRecordType `json:"allowedEndpointRecordTypes,omitempty"`
	DnsConfig                   *DnsConfig                   `json:"dnsConfig,omitempty"`
	Endpoints                   *[]Endpoint                  `json:"endpoints,omitempty"`
	MaxReturn                   *int64                       `json:"maxReturn,omitempty"`
	MonitorConfig               *MonitorConfig               `json:"monitorConfig,omitempty"`
	ProfileStatus               *ProfileStatus               `json:"profileStatus,omitempty"`
	TrafficRoutingMethod        *TrafficRoutingMethod        `json:"trafficRoutingMethod,omitempty"`
	TrafficViewEnrollmentStatus *TrafficViewEnrollmentStatus `json:"trafficViewEnrollmentStatus,omitempty"`
}
