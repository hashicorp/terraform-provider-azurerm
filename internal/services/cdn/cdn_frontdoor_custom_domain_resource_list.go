// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CdnFrontDoorCustomDomainListResource struct{}

type CdnFrontDoorCustomDomainListModel struct {
	ProfileId types.String `tfsdk:"cdn_frontdoor_profile_id"`
}

var _ sdk.FrameworkListWrappedResource = new(CdnFrontDoorCustomDomainListResource)

func (CdnFrontDoorCustomDomainListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = azurermCdnFrontDoorCustomDomainResourceName
}

func (CdnFrontDoorCustomDomainListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceCdnFrontDoorCustomDomain()
}

func (CdnFrontDoorCustomDomainListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"cdn_frontdoor_profile_id": listschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: profiles.ValidateProfileID},
				},
			},
		},
	}
}

func (CdnFrontDoorCustomDomainListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cdn.AFDCustomDomainsClient

	var data CdnFrontDoorCustomDomainListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parsedParentID, err := profiles.ParseProfileID(data.ProfileId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing `cdn_frontdoor_profile_id`", err)
		return
	}

	parentID := afddomains.NewProfileID(parsedParentID.SubscriptionId, parsedParentID.ResourceGroupName, parsedParentID.ProfileName)

	resp, err := client.AFDCustomDomainsListByProfileComplete(ctx, parentID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "listing "+azurermCdnFrontDoorCustomDomainResourceName, err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceCdnFrontDoorCustomDomain().Data(&terraform.InstanceState{})
			id, err := afddomains.ParseCustomDomainIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing "+azurermCdnFrontDoorCustomDomainResourceName+" ID", err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceCdnFrontDoorCustomDomainFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "encoding "+azurermCdnFrontDoorCustomDomainResourceName+" resource data", err)
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
