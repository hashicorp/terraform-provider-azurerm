package deletedworkspaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProperties struct {
	CreatedDate                         *string                         `json:"createdDate,omitempty"`
	CustomerId                          *string                         `json:"customerId,omitempty"`
	DefaultDataCollectionRuleResourceId *string                         `json:"defaultDataCollectionRuleResourceId,omitempty"`
	Failover                            *WorkspaceFailoverProperties    `json:"failover,omitempty"`
	Features                            *WorkspaceFeatures              `json:"features,omitempty"`
	ForceCmkForQuery                    *bool                           `json:"forceCmkForQuery,omitempty"`
	ModifiedDate                        *string                         `json:"modifiedDate,omitempty"`
	PrivateLinkScopedResources          *[]PrivateLinkScopedResource    `json:"privateLinkScopedResources,omitempty"`
	ProvisioningState                   *WorkspaceEntityStatus          `json:"provisioningState,omitempty"`
	PublicNetworkAccessForIngestion     *PublicNetworkAccessType        `json:"publicNetworkAccessForIngestion,omitempty"`
	PublicNetworkAccessForQuery         *PublicNetworkAccessType        `json:"publicNetworkAccessForQuery,omitempty"`
	Replication                         *WorkspaceReplicationProperties `json:"replication,omitempty"`
	RetentionInDays                     *int64                          `json:"retentionInDays,omitempty"`
	Sku                                 *WorkspaceSku                   `json:"sku,omitempty"`
	WorkspaceCapping                    *WorkspaceCapping               `json:"workspaceCapping,omitempty"`
}

func (o *WorkspaceProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkspaceProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}

func (o *WorkspaceProperties) GetModifiedDateAsTime() (*time.Time, error) {
	if o.ModifiedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ModifiedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkspaceProperties) SetModifiedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ModifiedDate = &formatted
}
