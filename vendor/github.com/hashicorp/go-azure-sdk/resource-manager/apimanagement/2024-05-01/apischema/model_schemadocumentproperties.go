package apischema

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaDocumentProperties struct {
	Components  *interface{} `json:"components,omitempty"`
	Definitions *interface{} `json:"definitions,omitempty"`
	Value       *string      `json:"value,omitempty"`
}
