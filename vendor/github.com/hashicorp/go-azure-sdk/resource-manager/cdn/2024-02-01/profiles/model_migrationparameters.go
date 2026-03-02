package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationParameters struct {
	ClassicResourceReference                ResourceReference                         `json:"classicResourceReference"`
	MigrationWebApplicationFirewallMappings *[]MigrationWebApplicationFirewallMapping `json:"migrationWebApplicationFirewallMappings,omitempty"`
	ProfileName                             string                                    `json:"profileName"`
	Sku                                     Sku                                       `json:"sku"`
}
