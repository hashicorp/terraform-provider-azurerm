package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FollowerDatabaseDefinition struct {
	AttachedDatabaseConfigurationName string                       `json:"attachedDatabaseConfigurationName"`
	ClusterResourceId                 string                       `json:"clusterResourceId"`
	DatabaseName                      *string                      `json:"databaseName,omitempty"`
	DatabaseShareOrigin               *DatabaseShareOrigin         `json:"databaseShareOrigin,omitempty"`
	TableLevelSharingProperties       *TableLevelSharingProperties `json:"tableLevelSharingProperties,omitempty"`
}
