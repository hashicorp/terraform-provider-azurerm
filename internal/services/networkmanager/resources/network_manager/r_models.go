// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager

type ManagerModel struct {
	CrossTenantScopes []ManagerCrossTenantScopeModel `tfschema:"cross_tenant_scopes"`
	Scope             []ManagerScopeModel            `tfschema:"scope"`
	ScopeAccesses     []string                       `tfschema:"scope_accesses"`
	Description       string                         `tfschema:"description"`
	Name              string                         `tfschema:"name"`
	Location          string                         `tfschema:"location"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Tags              map[string]any                 `tfschema:"tags"`
}

type ManagerScopeModel struct {
	ManagementGroups []string `tfschema:"management_group_ids"`
	Subscriptions    []string `tfschema:"subscription_ids"`
}

type ManagerCrossTenantScopeModel struct {
	TenantId         string   `tfschema:"tenant_id"`
	ManagementGroups []string `tfschema:"management_groups"`
	Subscriptions    []string `tfschema:"subscriptions"`
}

func (r Resource) ModelObject() any {
	return &ManagerModel{}
}
