package services

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchServiceProperties struct {
	AuthOptions                 *DataPlaneAuthOptions               `json:"authOptions,omitempty"`
	ComputeType                 *ComputeType                        `json:"computeType,omitempty"`
	DataExfiltrationProtections *[]SearchDataExfiltrationProtection `json:"dataExfiltrationProtections,omitempty"`
	DisableLocalAuth            *bool                               `json:"disableLocalAuth,omitempty"`
	ETag                        *string                             `json:"eTag,omitempty"`
	EncryptionWithCmk           *EncryptionWithCmk                  `json:"encryptionWithCmk,omitempty"`
	Endpoint                    *string                             `json:"endpoint,omitempty"`
	HostingMode                 *HostingMode                        `json:"hostingMode,omitempty"`
	NetworkRuleSet              *NetworkRuleSet                     `json:"networkRuleSet,omitempty"`
	PartitionCount              *int64                              `json:"partitionCount,omitempty"`
	PrivateEndpointConnections  *[]PrivateEndpointConnection        `json:"privateEndpointConnections,omitempty"`
	ProvisioningState           *ProvisioningState                  `json:"provisioningState,omitempty"`
	PublicNetworkAccess         *PublicNetworkAccess                `json:"publicNetworkAccess,omitempty"`
	ReplicaCount                *int64                              `json:"replicaCount,omitempty"`
	SemanticSearch              *SearchSemanticSearch               `json:"semanticSearch,omitempty"`
	ServiceUpgradedAt           *string                             `json:"serviceUpgradedAt,omitempty"`
	SharedPrivateLinkResources  *[]SharedPrivateLinkResource        `json:"sharedPrivateLinkResources,omitempty"`
	Status                      *SearchServiceStatus                `json:"status,omitempty"`
	StatusDetails               *string                             `json:"statusDetails,omitempty"`
	UpgradeAvailable            *UpgradeAvailable                   `json:"upgradeAvailable,omitempty"`
}

func (o *SearchServiceProperties) GetServiceUpgradedAtAsTime() (*time.Time, error) {
	if o.ServiceUpgradedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ServiceUpgradedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *SearchServiceProperties) SetServiceUpgradedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ServiceUpgradedAt = &formatted
}
