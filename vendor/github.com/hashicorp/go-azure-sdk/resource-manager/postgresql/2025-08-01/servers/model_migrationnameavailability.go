package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationNameAvailability struct {
	Message       *string                          `json:"message,omitempty"`
	Name          string                           `json:"name"`
	NameAvailable *bool                            `json:"nameAvailable,omitempty"`
	Reason        *MigrationNameAvailabilityReason `json:"reason,omitempty"`
	Type          string                           `json:"type"`
}
