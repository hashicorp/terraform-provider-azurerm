package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MapsAccountKeys struct {
	PrimaryKey              *string `json:"primaryKey,omitempty"`
	PrimaryKeyLastUpdated   *string `json:"primaryKeyLastUpdated,omitempty"`
	SecondaryKey            *string `json:"secondaryKey,omitempty"`
	SecondaryKeyLastUpdated *string `json:"secondaryKeyLastUpdated,omitempty"`
}
