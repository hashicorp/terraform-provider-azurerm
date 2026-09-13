package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountProperties struct {
	AbusePenalty                  *AbusePenalty                `json:"abusePenalty,omitempty"`
	AllowProjectManagement        *bool                        `json:"allowProjectManagement,omitempty"`
	AllowedFqdnList               *[]string                    `json:"allowedFqdnList,omitempty"`
	AmlWorkspace                  *UserOwnedAmlWorkspace       `json:"amlWorkspace,omitempty"`
	ApiProperties                 *ApiProperties               `json:"apiProperties,omitempty"`
	AssociatedProjects            *[]string                    `json:"associatedProjects,omitempty"`
	CallRateLimit                 *CallRateLimit               `json:"callRateLimit,omitempty"`
	Capabilities                  *[]SkuCapability             `json:"capabilities,omitempty"`
	CommitmentPlanAssociations    *[]CommitmentPlanAssociation `json:"commitmentPlanAssociations,omitempty"`
	CustomSubDomainName           *string                      `json:"customSubDomainName,omitempty"`
	DateCreated                   *string                      `json:"dateCreated,omitempty"`
	DefaultProject                *string                      `json:"defaultProject,omitempty"`
	DeletionDate                  *string                      `json:"deletionDate,omitempty"`
	DisableLocalAuth              *bool                        `json:"disableLocalAuth,omitempty"`
	DynamicThrottlingEnabled      *bool                        `json:"dynamicThrottlingEnabled,omitempty"`
	Encryption                    *Encryption                  `json:"encryption,omitempty"`
	Endpoint                      *string                      `json:"endpoint,omitempty"`
	Endpoints                     *map[string]string           `json:"endpoints,omitempty"`
	InternalId                    *string                      `json:"internalId,omitempty"`
	IsMigrated                    *bool                        `json:"isMigrated,omitempty"`
	Locations                     *MultiRegionSettings         `json:"locations,omitempty"`
	MigrationToken                *string                      `json:"migrationToken,omitempty"`
	NetworkAcls                   *NetworkRuleSet              `json:"networkAcls,omitempty"`
	NetworkInjections             *[]NetworkInjection          `json:"networkInjections,omitempty"`
	PrivateEndpointConnections    *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState             *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicNetworkAccess           *PublicNetworkAccess         `json:"publicNetworkAccess,omitempty"`
	QuotaLimit                    *QuotaLimit                  `json:"quotaLimit,omitempty"`
	RaiMonitorConfig              *RaiMonitorConfig            `json:"raiMonitorConfig,omitempty"`
	Restore                       *bool                        `json:"restore,omitempty"`
	RestrictOutboundNetworkAccess *bool                        `json:"restrictOutboundNetworkAccess,omitempty"`
	ScheduledPurgeDate            *string                      `json:"scheduledPurgeDate,omitempty"`
	SkuChangeInfo                 *SkuChangeInfo               `json:"skuChangeInfo,omitempty"`
	StoredCompletionsDisabled     *bool                        `json:"storedCompletionsDisabled,omitempty"`
	UserOwnedStorage              *[]UserOwnedStorage          `json:"userOwnedStorage,omitempty"`
}
