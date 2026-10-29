// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_project_pool

func (r Resource) ModelObject() any {
	return &DevCenterProjectPoolResourceModel{}
}

type DevCenterProjectPoolResourceModel struct {
	Name                               string            `tfschema:"name"`
	Location                           string            `tfschema:"location"`
	DevCenterProjectId                 string            `tfschema:"dev_center_project_id"`
	DevBoxDefinitionName               string            `tfschema:"dev_box_definition_name"`
	LocalAdministratorEnabled          bool              `tfschema:"local_administrator_enabled"`
	DevCenterAttachedNetworkName       string            `tfschema:"dev_center_attached_network_name"`
	ManagedVirtualNetworkRegions       []string          `tfschema:"managed_virtual_network_regions"`
	SingleSignOnEnabled                bool              `tfschema:"single_sign_on_enabled"`
	StopOnDisconnectGracePeriodMinutes int64             `tfschema:"stop_on_disconnect_grace_period_minutes"`
	Tags                               map[string]string `tfschema:"tags"`
}
