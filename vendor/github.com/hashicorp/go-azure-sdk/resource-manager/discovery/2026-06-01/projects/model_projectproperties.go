package projects

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectProperties struct {
	FoundryProjectEndpoint *string            `json:"foundryProjectEndpoint,omitempty"`
	ProvisioningState      *ProvisioningState `json:"provisioningState,omitempty"`
	Settings               *ProjectSettings   `json:"settings,omitempty"`
	StorageContainerIds    *[]string          `json:"storageContainerIds,omitempty"`
}
