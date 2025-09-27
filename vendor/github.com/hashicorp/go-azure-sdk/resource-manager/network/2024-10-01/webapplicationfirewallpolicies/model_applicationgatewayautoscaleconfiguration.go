package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayAutoscaleConfiguration struct {
	MaxCapacity *int64 `json:"maxCapacity,omitempty"`
	MinCapacity int64  `json:"minCapacity"`
}
