package netappaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LdapSearchScopeOpt struct {
	GroupDN               *string `json:"groupDN,omitempty"`
	GroupMembershipFilter *string `json:"groupMembershipFilter,omitempty"`
	UserDN                *string `json:"userDN,omitempty"`
}
