package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPConfigurationPublicIPAddressConfiguration struct {
	IPTags                 *[]IPTag                `json:"ipTags,omitempty"`
	Name                   string                  `json:"name"`
	PublicIPAddressVersion *PublicIPAddressVersion `json:"publicIPAddressVersion,omitempty"`
}
