// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_virtual_machine_scale_set

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

func maintenanceAssignmentVirtualMachineScaleSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"maintenance_configuration_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: maintenanceconfigurations.ValidateMaintenanceConfigurationID,
			// this is required because of an API issue https://github.com/Azure/azure-rest-api-specs/issues/8653
			// it can be removed when the API is fixed. Still an issue as of 2026-07-15.
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"virtual_machine_scale_set_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineScaleSetID,
		},
	}
}
