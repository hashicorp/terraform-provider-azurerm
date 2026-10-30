// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2024-07-10/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MaintenanceAssignmentArcMachineListResource struct{}

type MaintenanceAssignmentArcMachineListModel struct {
	ArcMachineId types.String `tfsdk:"arc_machine_id"`
}

var _ sdk.FrameworkListWrappedResource = MaintenanceAssignmentArcMachineListResource{}

func (MaintenanceAssignmentArcMachineListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(MaintenanceAssignmentArcMachineResource{})
}

func (MaintenanceAssignmentArcMachineListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = MaintenanceAssignmentArcMachineResource{}.ResourceType()
}

func (MaintenanceAssignmentArcMachineListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"arc_machine_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: machines.ValidateMachineID},
				},
			},
		},
	}
}

func (MaintenanceAssignmentArcMachineListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Maintenance.ConfigurationAssignmentsClient

	var data MaintenanceAssignmentArcMachineListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	arcMachineId, err := machines.ParseMachineID(data.ArcMachineId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Arc Machine ID", err)
		return
	}

	r := MaintenanceAssignmentArcMachineResource{}
	resp, err := client.List(ctx, commonids.NewScopeID(arcMachineId.ID()))
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}
	if resp.Model == nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), errors.New("model was nil"))
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, assignment := range pointer.From(resp.Model.Value) {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(assignment.Name)

			id, err := configurationassignments.ParseScopedConfigurationAssignmentIDInsensitively(pointer.From(assignment.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Maintenance Assignment ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			if err := r.flatten(rmd, id, &assignment); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(ctx, rmd.ResourceData, &result)
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
