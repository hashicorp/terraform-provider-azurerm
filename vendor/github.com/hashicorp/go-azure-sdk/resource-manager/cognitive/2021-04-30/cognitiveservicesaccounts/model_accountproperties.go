package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountProperties struct {
	AllowedFqdnList               *[]string                    `json:"allowedFqdnList,omitempty"`
	ApiProperties                 *ApiProperties               `json:"apiProperties,omitempty"`
	CallRateLimit                 *CallRateLimit               `json:"callRateLimit,omitempty"`
	Capabilities                  *[]SkuCapability             `json:"capabilities,omitempty"`
	CustomSubDomainName           *string                      `json:"customSubDomainName,omitempty"`
	DateCreated                   *string                      `json:"dateCreated,omitempty"`
	DisableLocalAuth              *bool                        `json:"disableLocalAuth,omitempty"`
	Encryption                    *Encryption                  `json:"encryption,omitempty"`
	Endpoint                      *string                      `json:"endpoint,omitempty"`
	Endpoints                     *map[string]string           `json:"endpoints,omitempty"`
	InternalId                    *string                      `json:"internalId,omitempty"`
	IsMigrated                    *bool                        `json:"isMigrated,omitempty"`
	MigrationToken                *string                      `json:"migrationToken,omitempty"`
	NetworkAcls                   *NetworkRuleSet              `json:"networkAcls,omitempty"`
	PrivateEndpointConnections    *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState             *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicNetworkAccess           *PublicNetworkAccess         `json:"publicNetworkAccess,omitempty"`
	QuotaLimit                    *QuotaLimit                  `json:"quotaLimit,omitempty"`
	Restore                       *bool                        `json:"restore,omitempty"`
	RestrictOutboundNetworkAccess *bool                        `json:"restrictOutboundNetworkAccess,omitempty"`
	SkuChangeInfo                 *SkuChangeInfo               `json:"skuChangeInfo,omitempty"`
	UserOwnedStorage              *[]UserOwnedStorage          `json:"userOwnedStorage,omitempty"`
}
