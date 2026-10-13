package trustedaccess

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrustedAccessRole struct {
	Name               *string                  `json:"name,omitempty"`
	Rules              *[]TrustedAccessRoleRule `json:"rules,omitempty"`
	SourceResourceType *string                  `json:"sourceResourceType,omitempty"`
}
