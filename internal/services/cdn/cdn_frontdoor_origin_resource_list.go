// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorOriginListResource struct{}

type CdnFrontDoorOriginListModel struct {
	OriginGroupId types.String `tfsdk:"cdn_frontdoor_origin_group_id"`
}

var _ sdk.FrameworkListWrappedResource = new(CdnFrontDoorOriginListResource)

func (CdnFrontDoorOriginListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureCdnFrontDoorOriginResourceName
}

func (CdnFrontDoorOriginListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceCdnFrontDoorOrigin()
}

func (CdnFrontDoorOriginListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"cdn_frontdoor_origin_group_id": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: afdorigins.ValidateOriginGroupID},
				},
			},
		},
	}
}

func (CdnFrontDoorOriginListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cdn.FrontDoorOriginsClient

	var data CdnFrontDoorOriginListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := afdorigins.ParseOriginGroupID(data.OriginGroupId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing `cdn_frontdoor_origin_group_id` for `%s`", azureCdnFrontDoorOriginResourceName), err)
		return
	}

	resp, err := client.ListByOriginGroupComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureCdnFrontDoorOriginResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceCdnFrontDoorOrigin().Data(&terraform.InstanceState{})

			id, err := afdorigins.ParseOriginGroupOriginIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", azureCdnFrontDoorOriginResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceCdnFrontDoorOriginFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureCdnFrontDoorOriginResourceName), err)
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
