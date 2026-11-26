// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/externalauths"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type RedHatHcpOpenShiftExternalAuthListResource struct{}

type RedHatHcpOpenShiftExternalAuthListModel struct {
	ClusterId types.String `tfsdk:"hcp_openshift_cluster_id"`
}

var _ sdk.FrameworkListWrappedResourceWithConfig = new(RedHatHcpOpenShiftExternalAuthListResource)

func (RedHatHcpOpenShiftExternalAuthListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = RedHatHcpOpenShiftExternalAuth{}.ResourceType()
}

func (RedHatHcpOpenShiftExternalAuthListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(RedHatHcpOpenShiftExternalAuth{})
}

func (RedHatHcpOpenShiftExternalAuthListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hcp_openshift_cluster_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: externalauths.ValidateHcpOpenShiftClusterID,
					},
				},
			},
		},
	}
}

func (RedHatHcpOpenShiftExternalAuthListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	var data RedHatHcpOpenShiftExternalAuthListModel
	if diags := request.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	clusterId, err := externalauths.ParseHcpOpenShiftClusterID(data.ClusterId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing HCP Open Shift Cluster ID", err)
		return
	}
	r := RedHatHcpOpenShiftExternalAuth{}
	resp, err := metadata.Client.RedHatOpenShift.ExternalAuthsClient.ListByParentComplete(ctx, *clusterId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}
	stream.Results = func(push func(list.ListResult) bool) {
		for _, externalAuth := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(externalAuth.Name)
			id, err := externalauths.ParseExternalAuthID(pointer.From(externalAuth.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing HCP Open Shift External Auth ID", err)
				return
			}
			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(id)
			if err := r.flatten(meta, *id, externalAuth); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}
			sdk.EncodeListResult(ctx, meta.ResourceData, &result)
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
