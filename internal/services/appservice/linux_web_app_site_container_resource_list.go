// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinuxWebAppSiteContainerListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(LinuxWebAppSiteContainerListResource)

func (LinuxWebAppSiteContainerListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(LinuxWebAppSiteContainerResource{})
}

func (LinuxWebAppSiteContainerListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = LinuxWebAppSiteContainerResource{}.ResourceType()
}

type LinuxWebAppSiteContainerListModel struct {
	LinuxWebAppId types.String `tfsdk:"linux_web_app_id"`
}

func (LinuxWebAppSiteContainerListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"linux_web_app_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateWebAppID,
					},
				},
			},
		},
	}
}

func (LinuxWebAppSiteContainerListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.AppService.WebAppsClient

	var data LinuxWebAppSiteContainerListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	appId, err := commonids.ParseWebAppID(data.LinuxWebAppId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing `linux_web_app_id`", err)
		return
	}

	r := LinuxWebAppSiteContainerResource{}

	resp, err := client.ListSiteContainersComplete(ctx, *appId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, container := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(container.Name)

			id, err := webapps.ParseSitecontainerIDInsensitively(pointer.From(container.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Site Container ID", err)
				return
			}

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)

			state := flattenLinuxWebAppSiteContainer(*id, container, "")
			if err := rmd.Encode(&state); err != nil {
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
