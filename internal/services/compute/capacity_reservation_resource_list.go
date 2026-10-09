// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservation"
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

func (r CapacityReservationListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Compute.CapacityReservationClient

	var data CapacityReservationListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	groupID, err := capacityreservation.ParseCapacityReservationGroupID(data.CapacityReservationGroupId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing parent ID for `%s`", azureCapacityReservationResourceName), err)
		return
	}
	resp, err := client.ListByCapacityReservationGroupComplete(ctx, *groupID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("retrieving Capacity Reservation Group for `%s`", azureCapacityReservationResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, item := range resp.Items {
			result := request.NewListResult(ctx)

			id, err := capacityreservations.ParseCapacityReservationIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", azureCapacityReservationResourceName), err)
				return
			}

			rd := resourceCapacityReservation().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceCapacityReservationFlatten(rd, id, r.capacityReservationToCapacityReservations(item)); err != nil {
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

// capacityReservationToCapacityReservations converts a returned capacity resevervation response to the type expected by the resource's flatten function
func (r CapacityReservationListResource) capacityReservationToCapacityReservations(item capacityreservation.CapacityReservation) *capacityreservations.CapacityReservation {
	result := &capacityreservations.CapacityReservation{
		Name:  item.Name,
		Type:  item.Type,
		Zones: item.Zones,
	}
	result.Sku.Name = item.Sku.Name
	result.Sku.Capacity = item.Sku.Capacity

	return result
}
