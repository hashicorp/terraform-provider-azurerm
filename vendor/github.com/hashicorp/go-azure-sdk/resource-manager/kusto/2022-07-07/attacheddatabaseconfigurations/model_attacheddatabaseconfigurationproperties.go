package attacheddatabaseconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttachedDatabaseConfigurationProperties struct {
	AttachedDatabaseNames             *[]string                         `json:"attachedDatabaseNames,omitempty"`
	ClusterResourceId                 string                            `json:"clusterResourceId"`
	DatabaseName                      string                            `json:"databaseName"`
	DatabaseNameOverride              *string                           `json:"databaseNameOverride,omitempty"`
	DatabaseNamePrefix                *string                           `json:"databaseNamePrefix,omitempty"`
	DefaultPrincipalsModificationKind DefaultPrincipalsModificationKind `json:"defaultPrincipalsModificationKind"`
	ProvisioningState                 *ProvisioningState                `json:"provisioningState,omitempty"`
	TableLevelSharingProperties       *TableLevelSharingProperties      `json:"tableLevelSharingProperties,omitempty"`
}
