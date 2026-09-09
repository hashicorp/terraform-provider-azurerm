package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayBackendHealthOnDemand struct {
	BackendAddressPool        *ApplicationGatewayBackendAddressPool        `json:"backendAddressPool,omitempty"`
	BackendHealthHTTPSettings *ApplicationGatewayBackendHealthHTTPSettings `json:"backendHealthHttpSettings,omitempty"`
}
