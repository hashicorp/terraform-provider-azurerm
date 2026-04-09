package fhirservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FhirServiceCorsConfiguration struct {
	AllowCredentials *bool     `json:"allowCredentials,omitempty"`
	Headers          *[]string `json:"headers,omitempty"`
	MaxAge           *int64    `json:"maxAge,omitempty"`
	Methods          *[]string `json:"methods,omitempty"`
	Origins          *[]string `json:"origins,omitempty"`
}
