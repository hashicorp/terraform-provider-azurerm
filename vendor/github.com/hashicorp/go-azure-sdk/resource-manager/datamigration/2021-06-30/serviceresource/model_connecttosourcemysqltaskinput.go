package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToSourceMySqlTaskInput struct {
	CheckPermissionsGroup *ServerLevelPermissionsGroup `json:"checkPermissionsGroup,omitempty"`
	IsOfflineMigration    *bool                        `json:"isOfflineMigration,omitempty"`
	SourceConnectionInfo  MySqlConnectionInfo          `json:"sourceConnectionInfo"`
	TargetPlatform        *MySqlTargetPlatformType     `json:"targetPlatform,omitempty"`
}
