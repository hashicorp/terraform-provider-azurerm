package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayBackendHealthHTTPSettings struct {
	BackendHTTPSettings *ApplicationGatewayBackendHTTPSettings   `json:"backendHttpSettings,omitempty"`
	Servers             *[]ApplicationGatewayBackendHealthServer `json:"servers,omitempty"`
}
