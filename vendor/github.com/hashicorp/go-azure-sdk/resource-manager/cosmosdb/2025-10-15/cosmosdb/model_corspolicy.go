package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CorsPolicy struct {
	AllowedHeaders  *string `json:"allowedHeaders,omitempty"`
	AllowedMethods  *string `json:"allowedMethods,omitempty"`
	AllowedOrigins  string  `json:"allowedOrigins"`
	ExposedHeaders  *string `json:"exposedHeaders,omitempty"`
	MaxAgeInSeconds *int64  `json:"maxAgeInSeconds,omitempty"`
}
