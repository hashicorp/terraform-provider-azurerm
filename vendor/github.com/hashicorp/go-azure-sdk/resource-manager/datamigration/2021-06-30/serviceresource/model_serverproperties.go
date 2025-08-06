package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerProperties struct {
	ServerDatabaseCount          *int64  `json:"serverDatabaseCount,omitempty"`
	ServerEdition                *string `json:"serverEdition,omitempty"`
	ServerName                   *string `json:"serverName,omitempty"`
	ServerOperatingSystemVersion *string `json:"serverOperatingSystemVersion,omitempty"`
	ServerPlatform               *string `json:"serverPlatform,omitempty"`
	ServerVersion                *string `json:"serverVersion,omitempty"`
}
