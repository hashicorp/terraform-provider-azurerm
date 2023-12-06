package schema

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalSchemaContractProperties struct {
	Description *string      `json:"description,omitempty"`
	Document    *interface{} `json:"document,omitempty"`
	SchemaType  SchemaType   `json:"schemaType"`
	Value       *interface{} `json:"value,omitempty"`
}
