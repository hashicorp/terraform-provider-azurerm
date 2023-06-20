package replicationvaultsetting

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VaultSettingProperties struct {
	MigrationSolutionId       *string `json:"migrationSolutionId,omitempty"`
	VMwareToAzureProviderType *string `json:"vmwareToAzureProviderType,omitempty"`
}
