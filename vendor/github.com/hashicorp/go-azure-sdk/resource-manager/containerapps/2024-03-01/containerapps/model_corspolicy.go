package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CorsPolicy struct {
	AllowCredentials *bool     `json:"allowCredentials,omitempty"`
	AllowedHeaders   *[]string `json:"allowedHeaders,omitempty"`
	AllowedMethods   *[]string `json:"allowedMethods,omitempty"`
	AllowedOrigins   []string  `json:"allowedOrigins"`
	ExposeHeaders    *[]string `json:"exposeHeaders,omitempty"`
	MaxAge           *int64    `json:"maxAge,omitempty"`
}
