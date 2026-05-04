package attestationproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JSONWebKey struct {
	Alg *string   `json:"alg,omitempty"`
	Crv *string   `json:"crv,omitempty"`
	D   *string   `json:"d,omitempty"`
	Dp  *string   `json:"dp,omitempty"`
	Dq  *string   `json:"dq,omitempty"`
	E   *string   `json:"e,omitempty"`
	K   *string   `json:"k,omitempty"`
	Kid *string   `json:"kid,omitempty"`
	Kty string    `json:"kty"`
	N   *string   `json:"n,omitempty"`
	P   *string   `json:"p,omitempty"`
	Q   *string   `json:"q,omitempty"`
	Qi  *string   `json:"qi,omitempty"`
	Use *string   `json:"use,omitempty"`
	X   *string   `json:"x,omitempty"`
	X5c *[]string `json:"x5c,omitempty"`
	Y   *string   `json:"y,omitempty"`
}
