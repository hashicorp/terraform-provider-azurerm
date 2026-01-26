package securitydomains

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityDomainJsonWebKey struct {
	Alg     string   `json:"alg"`
	E       string   `json:"e"`
	KeyOps  []string `json:"key_ops"`
	Kid     string   `json:"kid"`
	Kty     string   `json:"kty"`
	N       string   `json:"n"`
	Use     *string  `json:"use,omitempty"`
	X5c     []string `json:"x5c"`
	X5t     *string  `json:"x5t,omitempty"`
	X5tS256 string   `json:"x5t#S256"`
}
