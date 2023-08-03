package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPSecurityRestrictionRule struct {
	Action         Action  `json:"action"`
	Description    *string `json:"description,omitempty"`
	IPAddressRange string  `json:"ipAddressRange"`
	Name           string  `json:"name"`
}
