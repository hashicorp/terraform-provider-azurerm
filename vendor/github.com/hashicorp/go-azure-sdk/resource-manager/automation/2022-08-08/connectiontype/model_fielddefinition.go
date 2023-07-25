package connectiontype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FieldDefinition struct {
	IsEncrypted *bool  `json:"isEncrypted,omitempty"`
	IsOptional  *bool  `json:"isOptional,omitempty"`
	Type        string `json:"type"`
}
