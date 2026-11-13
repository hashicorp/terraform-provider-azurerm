package expressroutelags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLagAuthorizationPropertiesFormat struct {
	AuthorizationUseStatus *ExpressRouteLagAuthorizationUseStatus `json:"authorizationUseStatus,omitempty"`
	CircuitResourceUri     *string                                `json:"circuitResourceUri,omitempty"`
	ProvisioningState      *ProvisioningState                     `json:"provisioningState,omitempty"`
}
