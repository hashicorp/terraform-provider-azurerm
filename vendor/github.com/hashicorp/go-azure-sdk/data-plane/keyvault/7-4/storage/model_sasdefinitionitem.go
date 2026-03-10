package storage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SasDefinitionItem struct {
	Attributes *SasDefinitionAttributes `json:"attributes,omitempty"`
	Id         *string                  `json:"id,omitempty"`
	Sid        *string                  `json:"sid,omitempty"`
	Tags       *map[string]string       `json:"tags,omitempty"`
}
