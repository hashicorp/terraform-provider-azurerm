package servicegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceGatewayUpdateAddressLocationsRequest struct {
	Action           *UpdateAction                    `json:"action,omitempty"`
	AddressLocations *[]ServiceGatewayAddressLocation `json:"addressLocations,omitempty"`
}
