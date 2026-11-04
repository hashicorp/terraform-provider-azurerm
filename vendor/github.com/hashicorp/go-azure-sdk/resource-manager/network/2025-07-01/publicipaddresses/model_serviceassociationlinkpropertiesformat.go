package publicipaddresses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceAssociationLinkPropertiesFormat struct {
	AllowDelete        *bool              `json:"allowDelete,omitempty"`
	Link               *string            `json:"link,omitempty"`
	LinkedResourceType *string            `json:"linkedResourceType,omitempty"`
	Locations          *[]string          `json:"locations,omitempty"`
	ProvisioningState  *ProvisioningState `json:"provisioningState,omitempty"`
}
