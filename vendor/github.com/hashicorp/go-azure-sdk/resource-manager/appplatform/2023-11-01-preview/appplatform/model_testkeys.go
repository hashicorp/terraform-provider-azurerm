package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TestKeys struct {
	Enabled               *bool   `json:"enabled,omitempty"`
	PrimaryKey            *string `json:"primaryKey,omitempty"`
	PrimaryTestEndpoint   *string `json:"primaryTestEndpoint,omitempty"`
	SecondaryKey          *string `json:"secondaryKey,omitempty"`
	SecondaryTestEndpoint *string `json:"secondaryTestEndpoint,omitempty"`
}
