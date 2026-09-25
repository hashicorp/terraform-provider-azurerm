package privatelinkservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceNavigationLinkFormat struct {
	Link               *string            `json:"link,omitempty"`
	LinkedResourceType *string            `json:"linkedResourceType,omitempty"`
	ProvisioningState  *ProvisioningState `json:"provisioningState,omitempty"`
}
