package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VaultPropertiesRedundancySettings struct {
	CrossRegionRestore            *CrossRegionRestore            `json:"crossRegionRestore,omitempty"`
	StandardTierStorageRedundancy *StandardTierStorageRedundancy `json:"standardTierStorageRedundancy,omitempty"`
}
