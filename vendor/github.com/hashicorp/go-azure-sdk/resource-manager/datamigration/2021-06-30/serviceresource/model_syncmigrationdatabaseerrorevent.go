package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SyncMigrationDatabaseErrorEvent struct {
	EventText       *string `json:"eventText,omitempty"`
	EventTypeString *string `json:"eventTypeString,omitempty"`
	TimestampString *string `json:"timestampString,omitempty"`
}
