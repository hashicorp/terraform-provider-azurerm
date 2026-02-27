package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyCreateParameters struct {
	Attributes     *KeyAttributes         `json:"attributes,omitempty"`
	Crv            *JsonWebKeyCurveName   `json:"crv,omitempty"`
	KeyOps         *[]JsonWebKeyOperation `json:"key_ops,omitempty"`
	KeySize        *int64                 `json:"key_size,omitempty"`
	Kty            JsonWebKeyType         `json:"kty"`
	PublicExponent *int64                 `json:"public_exponent,omitempty"`
	ReleasePolicy  *KeyReleasePolicy      `json:"release_policy,omitempty"`
	Tags           *map[string]string     `json:"tags,omitempty"`
}
