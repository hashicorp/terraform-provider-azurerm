package indexers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FieldMappingFunction struct {
	Name       string       `json:"name"`
	Parameters *interface{} `json:"parameters,omitempty"`
}
