// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorOriginGroupListResource struct{}

type CdnFrontDoorOriginGroupListModel struct {
	ProfileId types.String `tfsdk:"cdn_frontdoor_profile_id"`
}

var _ sdk.FrameworkListWrappedResource = new(CdnFrontDoorOriginGroupListResource)

func (CdnFrontDoorOriginGroupListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azureCdnFrontDoorOriginGroupResourceName
}

func (CdnFrontDoorOriginGroupListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceCdnFrontDoorOriginGroup()
}

func (CdnFrontDoorOriginGroupListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"cdn_frontdoor_profile_id": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: afdorigingroups.ValidateProfileID,
					},
				},
			},
		},
	}
}

func (CdnFrontDoorOriginGroupListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cdn.FrontDoorOriginGroupsClient

	var data CdnFrontDoorOriginGroupListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := afdorigingroups.ParseProfileID(data.ProfileId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing CDN FrontDoor Profile ID for `%s`", azureCdnFrontDoorOriginGroupResourceName), err)
		return
	}

	resp, err := client.ListByProfileComplete(ctx, *parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", azureCdnFrontDoorOriginGroupResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceCdnFrontDoorOriginGroup().Data(&terraform.InstanceState{})

			id, err := afdorigingroups.ParseOriginGroupIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", azureCdnFrontDoorOriginGroupResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceCdnFrontDoorOriginGroupFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", azureCdnFrontDoorOriginGroupResourceName), err)
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
