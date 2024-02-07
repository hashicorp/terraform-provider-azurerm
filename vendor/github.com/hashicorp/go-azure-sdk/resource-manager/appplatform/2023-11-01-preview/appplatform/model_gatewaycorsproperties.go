package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayCorsProperties struct {
	AllowCredentials      *bool     `json:"allowCredentials,omitempty"`
	AllowedHeaders        *[]string `json:"allowedHeaders,omitempty"`
	AllowedMethods        *[]string `json:"allowedMethods,omitempty"`
	AllowedOriginPatterns *[]string `json:"allowedOriginPatterns,omitempty"`
	AllowedOrigins        *[]string `json:"allowedOrigins,omitempty"`
	ExposedHeaders        *[]string `json:"exposedHeaders,omitempty"`
	MaxAge                *int64    `json:"maxAge,omitempty"`
}
