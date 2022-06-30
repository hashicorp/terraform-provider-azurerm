package applicationgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationRequestProperties struct {
	MigrationPath *string    `json:"migrationPath,omitempty"`
	Operation     *Operation `json:"operation,omitempty"`
}
