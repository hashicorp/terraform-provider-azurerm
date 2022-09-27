package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationProperties struct {
	MigrationProgress *MigrationProgress `json:"migrationProgress,omitempty"`
	OldSubnetId       *string            `json:"oldSubnetId,omitempty"`
	OldVnetSiteId     *string            `json:"oldVnetSiteId,omitempty"`
}
