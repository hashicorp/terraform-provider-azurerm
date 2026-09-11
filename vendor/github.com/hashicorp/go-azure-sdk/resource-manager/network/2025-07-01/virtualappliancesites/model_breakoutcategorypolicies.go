package virtualappliancesites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BreakOutCategoryPolicies struct {
	Allow    *bool `json:"allow,omitempty"`
	Default  *bool `json:"default,omitempty"`
	Optimize *bool `json:"optimize,omitempty"`
}
