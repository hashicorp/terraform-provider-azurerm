package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBlobDefinition struct {
	AccountKey            *string                     `json:"accountKey,omitempty"`
	ContainerName         *string                     `json:"containerName,omitempty"`
	LocalAuthRef          *string                     `json:"localAuthRef,omitempty"`
	ManagedIdentity       *ManagedIdentityDefinition  `json:"managedIdentity,omitempty"`
	SasToken              *string                     `json:"sasToken,omitempty"`
	ServicePrincipal      *ServicePrincipalDefinition `json:"servicePrincipal,omitempty"`
	SyncIntervalInSeconds *int64                      `json:"syncIntervalInSeconds,omitempty"`
	TimeoutInSeconds      *int64                      `json:"timeoutInSeconds,omitempty"`
	Url                   *string                     `json:"url,omitempty"`
}
