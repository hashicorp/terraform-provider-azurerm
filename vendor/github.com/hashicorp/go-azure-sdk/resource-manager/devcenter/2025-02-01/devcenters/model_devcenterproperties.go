package devcenters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevCenterProperties struct {
	DevBoxProvisioningSettings *DevBoxProvisioningSettings      `json:"devBoxProvisioningSettings,omitempty"`
	DevCenterUri               *string                          `json:"devCenterUri,omitempty"`
	DisplayName                *string                          `json:"displayName,omitempty"`
	Encryption                 *Encryption                      `json:"encryption,omitempty"`
	NetworkSettings            *DevCenterNetworkSettings        `json:"networkSettings,omitempty"`
	ProjectCatalogSettings     *DevCenterProjectCatalogSettings `json:"projectCatalogSettings,omitempty"`
	ProvisioningState          *ProvisioningState               `json:"provisioningState,omitempty"`
}
