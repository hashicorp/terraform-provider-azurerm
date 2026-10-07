package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CdnMigrationToAfdParameters struct {
	MigrationEndpointMappings *[]MigrationEndpointMapping `json:"migrationEndpointMappings,omitempty"`
	Sku                       Sku                         `json:"sku"`
}
