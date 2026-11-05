package servicegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceGatewayAddressLocationResponse struct {
	AddressLocation *string                  `json:"addressLocation,omitempty"`
	Addresses       *[]ServiceGatewayAddress `json:"addresses,omitempty"`
}
