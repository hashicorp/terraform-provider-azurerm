package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyProperties struct {
	Crv        *JsonWebKeyCurveName `json:"crv,omitempty"`
	Exportable *bool                `json:"exportable,omitempty"`
	KeySize    *int64               `json:"key_size,omitempty"`
	Kty        *JsonWebKeyType      `json:"kty,omitempty"`
	ReuseKey   *bool                `json:"reuse_key,omitempty"`
}
