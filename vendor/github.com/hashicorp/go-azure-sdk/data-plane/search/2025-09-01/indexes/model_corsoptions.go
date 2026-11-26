package indexes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CorsOptions struct {
	AllowedOrigins  []string `json:"allowedOrigins"`
	MaxAgeInSeconds *int64   `json:"maxAgeInSeconds,omitempty"`
}
