package azuremonitorworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMonitorWorkspace struct {
	AccountId                  *string                      `json:"accountId,omitempty"`
	DefaultIngestionSettings   *IngestionSettings           `json:"defaultIngestionSettings,omitempty"`
	Metrics                    *Metrics                     `json:"metrics,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess         `json:"publicNetworkAccess,omitempty"`
}
