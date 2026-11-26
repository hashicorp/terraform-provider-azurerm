package authenticationpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationPolicyPropertiesFormat struct {
	AssociatedResources      *[]string                        `json:"associatedResources,omitempty"`
	AuthenticationProperties AuthenticationProviderProperties `json:"authenticationProperties"`
	OnUnauthenticatedRequest *OnUnauthenticatedRequest        `json:"onUnauthenticatedRequest,omitempty"`
	ProvisioningState        *ProvisioningState               `json:"provisioningState,omitempty"`
	ResourceGuid             *string                          `json:"resourceGuid,omitempty"`
	UserTrustProviderType    UserTrustProviderType            `json:"userTrustProviderType"`
}
