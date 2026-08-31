package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const azureCapacityReservationResourceName = "azurerm_capacity_reservation"

type CapacityReservationListResource struct{}

type CapacityReservationListModel struct {
	CapacityReservationGroupId types.String `tfsdk:"capacity_reservation_group_id"`
}

var _ sdk.FrameworkListWrappedResource = new(CapacityReservationListResource)

func (CapacityReservationListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureCapacityReservationResourceName
}

func (CapacityReservationListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceCapacityReservation()
}

func (CapacityReservationListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"capacity_reservation_group_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: capacityreservationgroups.ValidateCapacityReservationGroupID},
				},
			},
		},
	}
}

func (CapacityReservationListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Compute.CapacityReservationsClient
	groupsClient := metadata.Client.Compute.CapacityReservationGroupsClient

	var data CapacityReservationListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	groupID, err := capacityreservationgroups.ParseCapacityReservationGroupID(data.CapacityReservationGroupId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing parent ID for `%s`", azureCapacityReservationResourceName), err)
		return
	}

	groupResp, err := groupsClient.Get(ctx, *groupID, capacityreservationgroups.DefaultGetOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("retrieving Capacity Reservation Group for `%s`", azureCapacityReservationResourceName), err)
		return
	}
	if groupResp.Model == nil || groupResp.Model.Properties == nil || groupResp.Model.Properties.CapacityReservations == nil {
		stream.Results = func(push func(list.ListResult) bool) {}
		return
	}

	reservationRefs := *groupResp.Model.Properties.CapacityReservations

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", fmt.Errorf("context had no deadline"))
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, ref := range reservationRefs {
			result := request.NewListResult(ctx)

			id, err := capacityreservations.ParseCapacityReservationIDInsensitively(pointer.From(ref.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", azureCapacityReservationResourceName), err)
				return
			}

			result.DisplayName = id.CapacityReservationName

			resp, err := client.Get(ctx, *id, capacityreservations.DefaultGetOperationOptions())
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("retrieving `%s`", azureCapacityReservationResourceName), err)
				return
			}

			rd := resourceCapacityReservation().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceCapacityReservationFlatten(rd, id, resp.Model); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureCapacityReservationResourceName), err)
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
