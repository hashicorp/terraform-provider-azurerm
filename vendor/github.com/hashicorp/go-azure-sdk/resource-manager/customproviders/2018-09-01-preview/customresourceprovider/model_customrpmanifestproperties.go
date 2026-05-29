package customresourceprovider

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomRPManifestProperties struct {
	Actions           *[]CustomRPActionRouteDefinition       `json:"actions,omitempty"`
	ProvisioningState *ProvisioningState                     `json:"provisioningState,omitempty"`
	ResourceTypes     *[]CustomRPResourceTypeRouteDefinition `json:"resourceTypes,omitempty"`
	Validations       *[]CustomRPValidations                 `json:"validations,omitempty"`
}
