package networksecurityperimeterassociations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NspAssociationProperties struct {
	AccessMode            *AssociationAccessMode `json:"accessMode,omitempty"`
	HasProvisioningIssues *string                `json:"hasProvisioningIssues,omitempty"`
	PrivateLinkResource   *SubResource           `json:"privateLinkResource,omitempty"`
	Profile               *SubResource           `json:"profile,omitempty"`
	ProvisioningState     *NspProvisioningState  `json:"provisioningState,omitempty"`
}
