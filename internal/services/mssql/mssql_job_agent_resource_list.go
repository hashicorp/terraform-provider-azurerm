// Copyright IBM Corp.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobagents"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MssqlJobAgentListResource struct{}

type MssqlJobAgentListModel struct {
	MssqlServerId types.String `tfsdk:"mssql_server_id"`
}

var _ sdk.FrameworkListWrappedResource = new(MssqlJobAgentListResource)

func (r MssqlJobAgentListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMsSqlJobAgent()
}

func (r MssqlJobAgentListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_mssql_job_agent"
}

func (r MssqlJobAgentListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"mssql_server_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateSqlServerID,
					},
				},
			},
		},
	}
}

func (r MssqlJobAgentListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MSSQL.JobAgentsClient

	var data MssqlJobAgentListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]jobagents.JobAgent, 0)

	if !data.MssqlServerId.IsNull() {
		mssqlServerId, err := commonids.ParseSqlServerID(data.MssqlServerId.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing mssql Server ID for `%s`", "azurerm_mssql_job_agent"), err)
			return
		}

		resp, err := client.ListByServerComplete(ctx, *mssqlServerId)
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", "azurerm_mssql_job_agent"), err)
			return
		}

		results = resp.Items
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, jobagent := range results {
			result := request.NewListResult(ctx)

			result.DisplayName = pointer.From(jobagent.Name)

			rd := resourceMsSqlJobAgent().Data(&terraform.InstanceState{})

			id, err := jobagents.ParseJobAgentID(pointer.From(jobagent.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Mssql JobAgent ID", err)
				return
			}

			rd.SetId(id.ID())

			if err := resourceMssqlJobAgentSetFlatten(rd, id, &jobagent); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", "azurerm_mssql_job_agent"), err)
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
