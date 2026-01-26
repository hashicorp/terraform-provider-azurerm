package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JsonWebKey struct {
	Crv    *JsonWebKeyCurveName `json:"crv,omitempty"`
	D      *string              `json:"d,omitempty"`
	Dp     *string              `json:"dp,omitempty"`
	Dq     *string              `json:"dq,omitempty"`
	E      *string              `json:"e,omitempty"`
	K      *string              `json:"k,omitempty"`
	KeyHsm *string              `json:"key_hsm,omitempty"`
	KeyOps *[]string            `json:"key_ops,omitempty"`
	Kid    *string              `json:"kid,omitempty"`
	Kty    *JsonWebKeyType      `json:"kty,omitempty"`
	N      *string              `json:"n,omitempty"`
	P      *string              `json:"p,omitempty"`
	Q      *string              `json:"q,omitempty"`
	Qi     *string              `json:"qi,omitempty"`
	X      *string              `json:"x,omitempty"`
	Y      *string              `json:"y,omitempty"`
}
