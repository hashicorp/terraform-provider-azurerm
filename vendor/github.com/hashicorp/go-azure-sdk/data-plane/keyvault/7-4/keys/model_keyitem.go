package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyItem struct {
	Attributes *KeyAttributes     `json:"attributes,omitempty"`
	Kid        *string            `json:"kid,omitempty"`
	Managed    *bool              `json:"managed,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
