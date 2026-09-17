package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CorsSettings struct {
	AllowedOrigins     *[]string `json:"allowedOrigins,omitempty"`
	SupportCredentials *bool     `json:"supportCredentials,omitempty"`
}
