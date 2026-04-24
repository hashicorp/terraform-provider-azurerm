package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGateway struct {
	BackendAddressPools *[]ApplicationGatewayBackendAddressPool `json:"backendAddressPools,omitempty"`
	Resource            *string                                 `json:"resource,omitempty"`
}
