package blobservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CorsRule struct {
	AllowedHeaders  []string         `json:"allowedHeaders"`
	AllowedMethods  []AllowedMethods `json:"allowedMethods"`
	AllowedOrigins  []string         `json:"allowedOrigins"`
	ExposedHeaders  []string         `json:"exposedHeaders"`
	MaxAgeInSeconds int64            `json:"maxAgeInSeconds"`
}
