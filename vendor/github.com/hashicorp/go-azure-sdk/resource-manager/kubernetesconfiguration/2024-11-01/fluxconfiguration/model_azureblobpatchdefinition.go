package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBlobPatchDefinition struct {
	AccountKey            *string                          `json:"accountKey,omitempty"`
	ContainerName         *string                          `json:"containerName,omitempty"`
	LocalAuthRef          *string                          `json:"localAuthRef,omitempty"`
	ManagedIdentity       *ManagedIdentityPatchDefinition  `json:"managedIdentity,omitempty"`
	SasToken              *string                          `json:"sasToken,omitempty"`
	ServicePrincipal      *ServicePrincipalPatchDefinition `json:"servicePrincipal,omitempty"`
	SyncIntervalInSeconds *int64                           `json:"syncIntervalInSeconds,omitempty"`
	TimeoutInSeconds      *int64                           `json:"timeoutInSeconds,omitempty"`
	Url                   *string                          `json:"url,omitempty"`
}
