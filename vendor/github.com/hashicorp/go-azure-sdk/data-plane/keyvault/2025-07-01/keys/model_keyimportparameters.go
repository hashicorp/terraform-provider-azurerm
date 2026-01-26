package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyImportParameters struct {
	Attributes    *KeyAttributes     `json:"attributes,omitempty"`
	Hsm           *bool              `json:"Hsm,omitempty"`
	Key           JsonWebKey         `json:"key"`
	ReleasePolicy *KeyReleasePolicy  `json:"release_policy,omitempty"`
	Tags          *map[string]string `json:"tags,omitempty"`
}
