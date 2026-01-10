package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessCreateRoleBindingRequestModel struct {
	CrnPattern *string `json:"crn_pattern,omitempty"`
	Principal  *string `json:"principal,omitempty"`
	RoleName   *string `json:"role_name,omitempty"`
}
