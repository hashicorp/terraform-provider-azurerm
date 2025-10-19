package deploymentstacksatresourcegroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DenySettings struct {
	ApplyToChildScopes *bool            `json:"applyToChildScopes,omitempty"`
	ExcludedActions    *[]string        `json:"excludedActions,omitempty"`
	ExcludedPrincipals *[]string        `json:"excludedPrincipals,omitempty"`
	Mode               DenySettingsMode `json:"mode"`
}
