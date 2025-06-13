package trustedaccess

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrustedAccessRoleRule struct {
	ApiGroups       *[]string `json:"apiGroups,omitempty"`
	NonResourceURLs *[]string `json:"nonResourceURLs,omitempty"`
	ResourceNames   *[]string `json:"resourceNames,omitempty"`
	Resources       *[]string `json:"resources,omitempty"`
	Verbs           *[]string `json:"verbs,omitempty"`
}
