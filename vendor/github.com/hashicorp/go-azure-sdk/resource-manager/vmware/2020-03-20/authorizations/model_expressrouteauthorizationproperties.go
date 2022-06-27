package authorizations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteAuthorizationProperties struct {
	ExpressRouteAuthorizationId  *string                                     `json:"expressRouteAuthorizationId,omitempty"`
	ExpressRouteAuthorizationKey *string                                     `json:"expressRouteAuthorizationKey,omitempty"`
	ProvisioningState            *ExpressRouteAuthorizationProvisioningState `json:"provisioningState,omitempty"`
}
