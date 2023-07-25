package deviceupdates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GroupConnectivityInformation struct {
	CustomerVisibleFqdns        *[]string `json:"customerVisibleFqdns,omitempty"`
	GroupId                     *string   `json:"groupId,omitempty"`
	InternalFqdn                *string   `json:"internalFqdn,omitempty"`
	MemberName                  *string   `json:"memberName,omitempty"`
	PrivateLinkServiceArmRegion *string   `json:"privateLinkServiceArmRegion,omitempty"`
	RedirectMapId               *string   `json:"redirectMapId,omitempty"`
}
