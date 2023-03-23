package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicaSet struct {
	DomainControllerIPAddress *[]string        `json:"domainControllerIpAddress,omitempty"`
	ExternalAccessIPAddress   *string          `json:"externalAccessIpAddress,omitempty"`
	HealthAlerts              *[]HealthAlert   `json:"healthAlerts,omitempty"`
	HealthLastEvaluated       *string          `json:"healthLastEvaluated,omitempty"`
	HealthMonitors            *[]HealthMonitor `json:"healthMonitors,omitempty"`
	Location                  *string          `json:"location,omitempty"`
	ReplicaSetId              *string          `json:"replicaSetId,omitempty"`
	ServiceStatus             *string          `json:"serviceStatus,omitempty"`
	SubnetId                  *string          `json:"subnetId,omitempty"`
	VnetSiteId                *string          `json:"vnetSiteId,omitempty"`
}
