package compute

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type OrchestratedVirtualMachineScaleSetDataSource struct{}

var _ sdk.DataSource = OrchestratedVirtualMachineScaleSetDataSource{}

type OrchestratedVirtualMachineScaleSetDataSourceModel struct {
	Name          string `tfschema:"name"`
	ResourceGroup string `tfschema:"resource_group_name"`
	Location      string `tfschema:"location"`
}

func (r OrchestratedVirtualMachineScaleSetDataSource) ModelObject() interface{} {
	return &OrchestratedVirtualMachineScaleSetDataSourceModel{}
}

func (r OrchestratedVirtualMachineScaleSetDataSource) ResourceType() string {
	return "azurerm_orchestrated_virtual_machine_scale_set"
}

func (r OrchestratedVirtualMachineScaleSetDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: computeValidate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r OrchestratedVirtualMachineScaleSetDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VMScaleSetClient

			var orchestratedVMSS OrchestratedVirtualMachineScaleSetDataSource
			if err := metadata.Decode(&servicePlan); err != nil {
				return err
			}
		},
	}
}
