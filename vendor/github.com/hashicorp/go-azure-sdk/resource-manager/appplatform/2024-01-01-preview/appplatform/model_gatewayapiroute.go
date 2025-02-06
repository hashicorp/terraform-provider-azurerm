package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayApiRoute struct {
	Description *string   `json:"description,omitempty"`
	Filters     *[]string `json:"filters,omitempty"`
	Order       *int64    `json:"order,omitempty"`
	Predicates  *[]string `json:"predicates,omitempty"`
	SsoEnabled  *bool     `json:"ssoEnabled,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
	Title       *string   `json:"title,omitempty"`
	TokenRelay  *bool     `json:"tokenRelay,omitempty"`
	Uri         *string   `json:"uri,omitempty"`
}
