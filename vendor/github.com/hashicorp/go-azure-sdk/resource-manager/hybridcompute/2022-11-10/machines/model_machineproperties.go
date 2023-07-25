package machines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachineProperties struct {
	AdFqdn                     *string                         `json:"adFqdn,omitempty"`
	AgentConfiguration         *AgentConfiguration             `json:"agentConfiguration,omitempty"`
	AgentVersion               *string                         `json:"agentVersion,omitempty"`
	ClientPublicKey            *string                         `json:"clientPublicKey,omitempty"`
	CloudMetadata              *CloudMetadata                  `json:"cloudMetadata,omitempty"`
	DetectedProperties         *map[string]string              `json:"detectedProperties,omitempty"`
	DisplayName                *string                         `json:"displayName,omitempty"`
	DnsFqdn                    *string                         `json:"dnsFqdn,omitempty"`
	DomainName                 *string                         `json:"domainName,omitempty"`
	ErrorDetails               *[]ErrorDetail                  `json:"errorDetails,omitempty"`
	Extensions                 *[]MachineExtensionInstanceView `json:"extensions,omitempty"`
	LastStatusChange           *string                         `json:"lastStatusChange,omitempty"`
	LocationData               *LocationData                   `json:"locationData,omitempty"`
	MachineFqdn                *string                         `json:"machineFqdn,omitempty"`
	MssqlDiscovered            *string                         `json:"mssqlDiscovered,omitempty"`
	OsName                     *string                         `json:"osName,omitempty"`
	OsProfile                  *OSProfile                      `json:"osProfile,omitempty"`
	OsSku                      *string                         `json:"osSku,omitempty"`
	OsType                     *string                         `json:"osType,omitempty"`
	OsVersion                  *string                         `json:"osVersion,omitempty"`
	ParentClusterResourceId    *string                         `json:"parentClusterResourceId,omitempty"`
	PrivateLinkScopeResourceId *string                         `json:"privateLinkScopeResourceId,omitempty"`
	ProvisioningState          *string                         `json:"provisioningState,omitempty"`
	ServiceStatuses            *ServiceStatuses                `json:"serviceStatuses,omitempty"`
	Status                     *StatusTypes                    `json:"status,omitempty"`
	VMId                       *string                         `json:"vmId,omitempty"`
	VMUuid                     *string                         `json:"vmUuid,omitempty"`
}

func (o *MachineProperties) GetLastStatusChangeAsTime() (*time.Time, error) {
	if o.LastStatusChange == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStatusChange, "2006-01-02T15:04:05Z07:00")
}

func (o *MachineProperties) SetLastStatusChangeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStatusChange = &formatted
}
