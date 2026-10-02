package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchSchemaValue struct {
	DisplayName *string   `json:"displayName,omitempty"`
	Facet       bool      `json:"facet"`
	Indexed     bool      `json:"indexed"`
	Name        *string   `json:"name,omitempty"`
	OwnerType   *[]string `json:"ownerType,omitempty"`
	Stored      bool      `json:"stored"`
	Type        *string   `json:"type,omitempty"`
}
