package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationWebApplicationFirewallMapping struct {
	MigratedFrom *ResourceReference `json:"migratedFrom,omitempty"`
	MigratedTo   *ResourceReference `json:"migratedTo,omitempty"`
}
