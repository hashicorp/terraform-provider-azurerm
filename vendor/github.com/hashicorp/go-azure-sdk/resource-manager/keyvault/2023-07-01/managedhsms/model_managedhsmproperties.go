package managedhsms

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHsmProperties struct {
	CreateMode                 *CreateMode                          `json:"createMode,omitempty"`
	EnablePurgeProtection      *bool                                `json:"enablePurgeProtection,omitempty"`
	EnableSoftDelete           *bool                                `json:"enableSoftDelete,omitempty"`
	HsmUri                     *string                              `json:"hsmUri,omitempty"`
	InitialAdminObjectIds      *[]string                            `json:"initialAdminObjectIds,omitempty"`
	NetworkAcls                *MHSMNetworkRuleSet                  `json:"networkAcls,omitempty"`
	PrivateEndpointConnections *[]MHSMPrivateEndpointConnectionItem `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState                   `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess                 `json:"publicNetworkAccess,omitempty"`
	Regions                    *[]MHSMGeoReplicatedRegion           `json:"regions,omitempty"`
	ScheduledPurgeDate         *string                              `json:"scheduledPurgeDate,omitempty"`
	SecurityDomainProperties   *ManagedHSMSecurityDomainProperties  `json:"securityDomainProperties,omitempty"`
	SoftDeleteRetentionInDays  *int64                               `json:"softDeleteRetentionInDays,omitempty"`
	StatusMessage              *string                              `json:"statusMessage,omitempty"`
	TenantId                   *string                              `json:"tenantId,omitempty"`
}

func (o *ManagedHsmProperties) GetScheduledPurgeDateAsTime() (*time.Time, error) {
	if o.ScheduledPurgeDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ScheduledPurgeDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ManagedHsmProperties) SetScheduledPurgeDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ScheduledPurgeDate = &formatted
}
