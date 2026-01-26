package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyBundle struct {
	Attributes    *KeyAttributes     `json:"attributes,omitempty"`
	Key           *JsonWebKey        `json:"key,omitempty"`
	Managed       *bool              `json:"managed,omitempty"`
	ReleasePolicy *KeyReleasePolicy  `json:"release_policy,omitempty"`
	Tags          *map[string]string `json:"tags,omitempty"`
}
