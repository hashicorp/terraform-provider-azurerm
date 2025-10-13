package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServiceReference struct {
	Parameters    *map[string]interface{} `json:"parameters,omitempty"`
	ReferenceName string                  `json:"referenceName"`
	Type          Type                    `json:"type"`
}
