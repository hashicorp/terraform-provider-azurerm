// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesprojects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	response.Schema = cognitiveAccountProjectConnectionListResourceConfigSchema()
}

func (CognitiveAccountProjectConnectionEntraIDListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.ProjectConnectionResourceClient

	var data cognitiveAccountProjectConnectionListModel
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

	projects, err := cognitiveAccountProjectConnectionListProjects(ctx, metadata, data)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", CognitiveAccountProjectConnectionEntraIDResource{}.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		listCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, project := range projects {
			projectId, err := cognitiveservicesprojects.ParseProjectID(pointer.From(project.Id))
			if err != nil {
				result := request.NewListResult(listCtx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Account Project ID", err)
				return
			}

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
				if string(base.AuthType) != string(projectconnectionresource.ConnectionAuthTypeAAD) {
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

				rd := sdk.WrappedResource(CognitiveAccountProjectConnectionEntraIDResource{}).Data(&terraform.InstanceState{})
				rd.SetId(connectionId.ID())
				if err := pluginsdk.SetResourceIdentityData(rd, connectionId); err != nil {
					sdk.SetErrorDiagnosticAndPushListResult(result, push, "setting Cognitive Account Project Connection identity", err)
					return
				}
				_ = rd.Set("name", connectionId.ConnectionName)
				_ = rd.Set("cognitive_account_project_id", projectconnectionresource.NewProjectID(connectionId.SubscriptionId, connectionId.ResourceGroupName, connectionId.AccountName, connectionId.ProjectName).ID())
				_ = rd.Set("category", pointer.FromEnum(base.Category))
				_ = rd.Set("target", pointer.From(base.Target))
				_ = rd.Set("metadata", pointer.From(base.Metadata))

				sdk.EncodeListResult(listCtx, rd, &result)
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
}
