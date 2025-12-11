// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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

type VirtualMachinePowerAction struct{}

func (v VirtualMachinePowerAction) ModelObject() any {
	return &VirtualMachinePowerActionModel{}
}

func (v VirtualMachinePowerAction) Timeout() time.Duration {
	return 30 * time.Minute
}

var _ sdk.WrappedAction = VirtualMachinePowerAction{}

type VirtualMachinePowerActionModel struct {
	sdk.BaseActionModel

	VirtualMachineId types.String `tfsdk:"virtual_machine_id"`
	Action           types.String `tfsdk:"power_action"`
}

func (v VirtualMachinePowerAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
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

			//"timeout": schema.StringAttribute{
			//	Optional:            true,
			//	Description:         "Timeout duration for the action to complete. Defaults to `30m`.",
			//	MarkdownDescription: "Timeout duration for the action to complete. Defaults to `30m`.",
			//},
		},
	}
}

func (v VirtualMachinePowerAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "azurerm_virtual_machine_power"
}

func (v VirtualMachinePowerAction) Invoke(ctx context.Context, _ action.InvokeRequest, resp *action.InvokeResponse, config any, metadata sdk.ActionMetadata) {
	client := metadata.Client.Compute.VirtualMachinesClient

	model := sdk.AssertActionModelType[VirtualMachinePowerActionModel](config, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := virtualmachines.ParseVirtualMachineID(model.VirtualMachineId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing id", err)
		return
	}

	powerAction := model.Action.ValueString()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("invoking %s on %s", powerAction, id.VirtualMachineName),
	})

	switch powerAction {
	case "restart":
		if err := client.RestartThenPoll(ctx, *id); err != nil {
			sdk.SetResponseErrorDiagnostic(resp, "running action", fmt.Sprintf("restarting %s: %+v", id, err))
			return
		}

	case "power_on":
		if err := client.StartThenPoll(ctx, *id); err != nil {
			sdk.SetResponseErrorDiagnostic(resp, "running action", fmt.Sprintf("starting %s: %+v", id, err))
			return
		}

	case "power_off":
		if err := client.PowerOffThenPoll(ctx, *id, virtualmachines.DefaultPowerOffOperationOptions()); err != nil {
			sdk.SetResponseErrorDiagnostic(resp, "running action", fmt.Sprintf("stopping %s: %+v", id, err))
			return
		}
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("action %s on %s completed", powerAction, id.VirtualMachineName),
	})
}
