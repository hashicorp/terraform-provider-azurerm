package datadogmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerBillingEntity struct {
	Id               *string `json:"id,omitempty"`
	Name             *string `json:"name,omitempty"`
	PartnerEntityUri *string `json:"partnerEntityUri,omitempty"`
}
