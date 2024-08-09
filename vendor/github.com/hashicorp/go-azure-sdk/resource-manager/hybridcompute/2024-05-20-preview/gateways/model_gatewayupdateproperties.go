package gateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayUpdateProperties struct {
	AllowedFeatures *[]string `json:"allowedFeatures,omitempty"`
}
