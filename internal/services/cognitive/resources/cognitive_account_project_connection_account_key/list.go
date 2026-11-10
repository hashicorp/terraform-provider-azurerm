// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_account_key

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(ListResource)

func (ListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(Resource{})
}

func (ListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = Resource{}.ResourceType()
}

func (ListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = cognitiveAccountProjectConnectionAccountKeyListResourceConfigSchema()
}

type cognitiveAccountProjectConnectionAccountKeyListModel struct {
	CognitiveAccountProjectId types.String `tfsdk:"cognitive_account_project_id"`
}

func cognitiveAccountProjectConnectionAccountKeyListResourceConfigSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cognitive_account_project_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: projectconnectionresource.ValidateProjectID,
					},
				},
			},
		},
	}
}

func (ListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.ProjectConnectionResourceClient

	var data cognitiveAccountProjectConnectionAccountKeyListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	projectId, err := projectconnectionresource.ParseProjectID(data.CognitiveAccountProjectId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Cognitive Account Project ID", err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		listCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		connectionsResp, err := client.ProjectConnectionsListComplete(listCtx, *projectId, projectconnectionresource.DefaultProjectConnectionsListOperationOptions())
		if err != nil {
			result := request.NewListResult(listCtx)
			sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("listing connections for `%s`", projectId.ProjectName), err)
			return
		}

		for _, connection := range connectionsResp.Items {
			if connection.Properties == nil {
				continue
			}

			base := connection.Properties.ConnectionPropertiesV2()
			if base.AuthType != projectconnectionresource.ConnectionAuthTypeAccountKey {
				continue
			}

			connectionId, err := projectconnectionresource.ParseProjectConnectionIDInsensitively(pointer.From(connection.Id))
			if err != nil {
				result := request.NewListResult(listCtx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Account Project Connection ID", err)
				return
			}

			result := request.NewListResult(listCtx)
			result.DisplayName = pointer.From(connection.Name)

			r := Resource{}
			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(connectionId)

			if err := r.flatten(meta, connectionId, &connection, nil, ""); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", r.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(listCtx, meta.ResourceData, &result)
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
