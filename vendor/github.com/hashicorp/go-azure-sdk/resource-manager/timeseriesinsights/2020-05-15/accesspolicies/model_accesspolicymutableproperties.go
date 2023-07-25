package accesspolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyMutableProperties struct {
	Description *string             `json:"description,omitempty"`
	Roles       *[]AccessPolicyRole `json:"roles,omitempty"`
}
