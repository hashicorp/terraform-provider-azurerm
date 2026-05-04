package projects

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectProperties struct {
	CatalogSettings    *ProjectCatalogSettings `json:"catalogSettings,omitempty"`
	Description        *string                 `json:"description,omitempty"`
	DevCenterId        *string                 `json:"devCenterId,omitempty"`
	DevCenterUri       *string                 `json:"devCenterUri,omitempty"`
	DisplayName        *string                 `json:"displayName,omitempty"`
	MaxDevBoxesPerUser *int64                  `json:"maxDevBoxesPerUser,omitempty"`
	ProvisioningState  *ProvisioningState      `json:"provisioningState,omitempty"`
}
