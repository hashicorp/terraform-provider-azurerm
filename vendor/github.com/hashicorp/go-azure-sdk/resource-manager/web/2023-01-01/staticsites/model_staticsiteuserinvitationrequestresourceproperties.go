package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteUserInvitationRequestResourceProperties struct {
	Domain               *string `json:"domain,omitempty"`
	NumHoursToExpiration *int64  `json:"numHoursToExpiration,omitempty"`
	Provider             *string `json:"provider,omitempty"`
	Roles                *string `json:"roles,omitempty"`
	UserDetails          *string `json:"userDetails,omitempty"`
}
