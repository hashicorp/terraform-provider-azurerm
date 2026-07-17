package connectiontype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionTypeCreateOrUpdateProperties struct {
	FieldDefinitions map[string]FieldDefinition `json:"fieldDefinitions"`
	IsGlobal         *bool                      `json:"isGlobal,omitempty"`
}
