package services

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchServiceProperties struct {
	AuthOptions                     *DataPlaneAuthOptions                   `json:"authOptions,omitempty"`
	DisableLocalAuth                *bool                                   `json:"disableLocalAuth,omitempty"`
	DisabledDataExfiltrationOptions *[]SearchDisabledDataExfiltrationOption `json:"disabledDataExfiltrationOptions,omitempty"`
	ETag                            *string                                 `json:"eTag,omitempty"`
	EncryptionWithCmk               *EncryptionWithCmk                      `json:"encryptionWithCmk,omitempty"`
	HostingMode                     *HostingMode                            `json:"hostingMode,omitempty"`
	NetworkRuleSet                  *NetworkRuleSet                         `json:"networkRuleSet,omitempty"`
	PartitionCount                  *int64                                  `json:"partitionCount,omitempty"`
	PrivateEndpointConnections      *[]PrivateEndpointConnection            `json:"privateEndpointConnections,omitempty"`
	ProvisioningState               *ProvisioningState                      `json:"provisioningState,omitempty"`
	PublicNetworkAccess             *PublicNetworkAccess                    `json:"publicNetworkAccess,omitempty"`
	ReplicaCount                    *int64                                  `json:"replicaCount,omitempty"`
	SemanticSearch                  *SearchSemanticSearch                   `json:"semanticSearch,omitempty"`
	SharedPrivateLinkResources      *[]SharedPrivateLinkResource            `json:"sharedPrivateLinkResources,omitempty"`
	Status                          *SearchServiceStatus                    `json:"status,omitempty"`
	StatusDetails                   *string                                 `json:"statusDetails,omitempty"`
}
