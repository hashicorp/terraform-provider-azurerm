package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceListKeys struct {
	KeyName                   *string `json:"keyName,omitempty"`
	PrimaryConnectionString   *string `json:"primaryConnectionString,omitempty"`
	PrimaryKey                *string `json:"primaryKey,omitempty"`
	SecondaryConnectionString *string `json:"secondaryConnectionString,omitempty"`
	SecondaryKey              *string `json:"secondaryKey,omitempty"`
}
