package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = AssignmentVirtualMachineScaleSetV0ToV1{}

type AssignmentVirtualMachineScaleSetV0ToV1 struct {
}

func (AssignmentVirtualMachineScaleSetV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"maintenance_configuration_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"virtual_machine_scale_set_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (AssignmentVirtualMachineScaleSetV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := configurationassignments.ParseScopedConfigurationAssignmentIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, fmt.Errorf("parsing old id %q: %+v", oldIdRaw, err)
		}

		virtualMachineScaleSetId, err := commonids.ParseVirtualMachineScaleSetIDInsensitively(oldId.Scope)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a virtual machine scale set id: %+v", oldId.Scope, err)
		}

		newId := configurationassignments.NewScopedConfigurationAssignmentID(virtualMachineScaleSetId.ID(), oldId.ConfigurationAssignmentName)
		newIdRaw := newId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newIdRaw)
		rawState["id"] = newIdRaw

		return rawState, nil
	}
}
