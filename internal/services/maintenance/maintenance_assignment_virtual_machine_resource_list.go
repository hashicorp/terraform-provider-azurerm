// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MaintenanceAssignmentVirtualMachineListResource struct{}

type MaintenanceAssignmentVirtualMachineListModel struct {
	VirtualMachineId types.String `tfsdk:"virtual_machine_id"`
}

var _ sdk.FrameworkListWrappedResource = &MaintenanceAssignmentVirtualMachineListResource{}

func (MaintenanceAssignmentVirtualMachineListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_maintenance_assignment_virtual_machine"
}

func (MaintenanceAssignmentVirtualMachineListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceArmMaintenanceAssignmentVirtualMachine()
}

func (MaintenanceAssignmentVirtualMachineListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"virtual_machine_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: virtualmachines.ValidateVirtualMachineID,
					},
				},
			},
		},
	}
}

func (MaintenanceAssignmentVirtualMachineListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

	var data MaintenanceAssignmentVirtualMachineListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	vmId, err := virtualmachines.ParseVirtualMachineID(data.VirtualMachineId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing `virtual_machine_id`", err)
		return
	}

	scopeId := commonids.NewScopeID(vmId.ID())
	resp, err := client.List(ctx, scopeId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing maintenance configuration assignments for `%s`", vmId), err)
		return
	}

	var items []configurationassignments.ConfigurationAssignment
	if resp.Model != nil && resp.Model.Value != nil {
		items = *resp.Model.Value
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rawId := pointer.From(item.Id)
			if !strings.HasPrefix(rawId, "/") { // List response trims leading / on scoped IDs.
				rawId = "/" + rawId
			}

			id, err := configurationassignments.ParseScopedConfigurationAssignmentIDInsensitively(rawId)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Configuration Assignment ID", err)
				return
			}

			rd := resourceArmMaintenanceAssignmentVirtualMachine().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceArmMaintenanceAssignmentVirtualMachineFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting resource data for `azurerm_maintenance_assignment_virtual_machine`", err)
				return
			}

			sdk.EncodeListResult(ctx, rd, &result)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
