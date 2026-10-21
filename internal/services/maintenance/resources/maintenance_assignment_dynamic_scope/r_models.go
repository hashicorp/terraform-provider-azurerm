// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_dynamic_scope

type MaintenanceDynamicScopeModel struct {
	Name                       string   `tfschema:"name"`
	MaintenanceConfigurationId string   `tfschema:"maintenance_configuration_id"`
	Filter                     []Filter `tfschema:"filter"`
}

func (Resource) ModelObject() any {
	return &MaintenanceDynamicScopeModel{}
}
