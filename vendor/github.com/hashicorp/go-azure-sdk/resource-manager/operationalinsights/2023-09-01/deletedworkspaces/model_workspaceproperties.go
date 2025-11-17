package deletedworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProperties struct {
	CreatedDate                         *string                      `json:"createdDate,omitempty"`
	CustomerId                          *string                      `json:"customerId,omitempty"`
	DefaultDataCollectionRuleResourceId *string                      `json:"defaultDataCollectionRuleResourceId,omitempty"`
	Features                            *WorkspaceFeatures           `json:"features,omitempty"`
	ForceCmkForQuery                    *bool                        `json:"forceCmkForQuery,omitempty"`
	ModifiedDate                        *string                      `json:"modifiedDate,omitempty"`
	PrivateLinkScopedResources          *[]PrivateLinkScopedResource `json:"privateLinkScopedResources,omitempty"`
	ProvisioningState                   *WorkspaceEntityStatus       `json:"provisioningState,omitempty"`
	PublicNetworkAccessForIngestion     *PublicNetworkAccessType     `json:"publicNetworkAccessForIngestion,omitempty"`
	PublicNetworkAccessForQuery         *PublicNetworkAccessType     `json:"publicNetworkAccessForQuery,omitempty"`
	RetentionInDays                     *int64                       `json:"retentionInDays,omitempty"`
	Sku                                 *WorkspaceSku                `json:"sku,omitempty"`
	WorkspaceCapping                    *WorkspaceCapping            `json:"workspaceCapping,omitempty"`
}
