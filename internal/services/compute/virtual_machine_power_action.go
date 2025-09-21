package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type VirtualMachinePowerAction struct {
	sdk.ActionMetadata
}

var _ sdk.Action = &VirtualMachinePowerAction{}

func newVirtualMachinePowerAction() action.Action {
	return &VirtualMachinePowerAction{}
}

type VirtualMachinePowerActionModel struct {
	VirtualMachineId types.String `tfsdk:"virtual_machine_id"`
	Action           types.String `tfsdk:"power_action"`
	// TODO - Custom timeouts?
}

func (v *VirtualMachinePowerAction) Schema(_ context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"virtual_machine_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the virtual machine on which to perform the action.",
				MarkdownDescription: "The ID of the virtual machine on which to perform the action.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: virtualmachines.ValidateVirtualMachineID,
					},
				},
			},

			"power_action": schema.StringAttribute{
				Required:            true,
				Description:         "The power state action to take on this virtual machine. Possible values include `restart`, `power_on`, and `power_off`.",
				MarkdownDescription: "The power state action to take on this virtual machine. Possible values include `restart`, `power_on`, and `power_off`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"power_on",
						"power_off",
						"restart",
					),
				},
			},
		},
	}
}

func (v *VirtualMachinePowerAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_virtual_machine_power"
}

func (v *VirtualMachinePowerAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := v.Client.Compute.VirtualMachinesClient

	ctx, cancel := context.WithTimeout(ctx, time.Minute*15)
	defer cancel()

	model := VirtualMachinePowerActionModel{}

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}

	id, err := virtualmachines.ParseVirtualMachineID(model.VirtualMachineId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, "parsing id", err)
		return
	}

	powerAction := model.Action.ValueString()

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("invoking %s on %s", powerAction, id.VirtualMachineName),
	})

	switch powerAction {
	case "restart":
		if err := client.RestartThenPoll(ctx, *id); err != nil {
			sdk.SetResponseErrorDiagnostic(response, "running action", fmt.Sprintf("restarting %s: %+v", id, err))
		}

	case "power_on":
		if err := client.StartThenPoll(ctx, *id); err != nil {
			sdk.SetResponseErrorDiagnostic(response, "running action", fmt.Sprintf("starting %s: %+v", id, err))
		}

	case "power_off":
		if err := client.PowerOffThenPoll(ctx, *id, virtualmachines.DefaultPowerOffOperationOptions()); err != nil {
			sdk.SetResponseErrorDiagnostic(response, "running action", fmt.Sprintf("stopping %s: %+v", id, err))
		}
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("action %s on %s completed", powerAction, id.VirtualMachineName),
	})
}

func (v *VirtualMachinePowerAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	v.Defaults(ctx, request, response)
}
