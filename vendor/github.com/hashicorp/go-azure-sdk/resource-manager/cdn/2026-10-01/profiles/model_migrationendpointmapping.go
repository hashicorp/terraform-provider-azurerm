package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationEndpointMapping struct {
	MigratedFrom *string `json:"migratedFrom,omitempty"`
	MigratedTo   *string `json:"migratedTo,omitempty"`
}
