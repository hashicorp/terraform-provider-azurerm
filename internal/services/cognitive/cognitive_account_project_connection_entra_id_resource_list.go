// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

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

type CognitiveAccountProjectConnectionEntraIDListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CognitiveAccountProjectConnectionEntraIDListResource)

func (CognitiveAccountProjectConnectionEntraIDListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CognitiveAccountProjectConnectionEntraIDResource{})
}

func (CognitiveAccountProjectConnectionEntraIDListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CognitiveAccountProjectConnectionEntraIDResource{}.ResourceType()
}

func (CognitiveAccountProjectConnectionEntraIDListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = cognitiveAccountProjectConnectionEntraIDListResourceConfigSchema()
}

type cognitiveAccountProjectConnectionEntraIDListModel struct {
	CognitiveAccountProjectId types.String `tfsdk:"cognitive_account_project_id"`
}

func cognitiveAccountProjectConnectionEntraIDListResourceConfigSchema() schema.Schema {
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

func (CognitiveAccountProjectConnectionEntraIDListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.ProjectConnectionResourceClient

	var data cognitiveAccountProjectConnectionEntraIDListModel
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

		connProjectId := projectconnectionresource.NewProjectID(projectId.SubscriptionId, projectId.ResourceGroupName, projectId.AccountName, projectId.ProjectName)
		connectionsResp, err := client.ProjectConnectionsListComplete(listCtx, connProjectId, projectconnectionresource.DefaultProjectConnectionsListOperationOptions())
		if err != nil {
			result := request.NewListResult(listCtx)
			sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("listing connections for project `%s`", projectId.ProjectName), err)
			return
		}

		for _, connection := range connectionsResp.Items {
			if connection.Properties == nil {
				continue
			}

			base := connection.Properties.ConnectionPropertiesV2()
			if base.AuthType != projectconnectionresource.ConnectionAuthTypeAAD {
				continue
			}

			connectionId, err := projectconnectionresource.ParseProjectConnectionID(pointer.From(connection.Id))
			if err != nil {
				result := request.NewListResult(listCtx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Account Project Connection ID", err)
				return
			}

			result := request.NewListResult(listCtx)
			result.DisplayName = pointer.From(connection.Name)

			r := CognitiveAccountProjectConnectionEntraIDResource{}
			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(connectionId)

			if err := r.flatten(meta, connectionId, &connection, nil); err != nil {
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
