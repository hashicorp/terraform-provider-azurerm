package apioperation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RepresentationContract struct {
	ContentType    string                               `json:"contentType"`
	Examples       *map[string]ParameterExampleContract `json:"examples,omitempty"`
	FormParameters *[]ParameterContract                 `json:"formParameters,omitempty"`
	SchemaId       *string                              `json:"schemaId,omitempty"`
	TypeName       *string                              `json:"typeName,omitempty"`
}
