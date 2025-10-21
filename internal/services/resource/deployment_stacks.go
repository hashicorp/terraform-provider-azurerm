// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

type ActionOnUnmanageModel struct {
	Resources        string `tfschema:"resources"`
	ResourceGroups   string `tfschema:"resource_groups"`
	ManagementGroups string `tfschema:"management_groups"`
}

type DenySettingsModel struct {
	Mode               string    `tfschema:"mode"`
	ApplyToChildScopes bool      `tfschema:"apply_to_child_scopes"`
	ExcludedActions    *[]string `tfschema:"excluded_actions"`
	ExcludedPrincipals *[]string `tfschema:"excluded_principals"`
}
