// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sensitivitylabels"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MssqlDatabaseSensitivityLabelListResource struct{}

type MssqlDatabaseSensitivityLabelListModel struct {
	DatabaseId types.String `tfsdk:"database_id"`
}

var _ sdk.FrameworkListWrappedResource = MssqlDatabaseSensitivityLabelListResource{}

func (MssqlDatabaseSensitivityLabelListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(MssqlDatabaseSensitivityLabelResource{})
}

func (MssqlDatabaseSensitivityLabelListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = MssqlDatabaseSensitivityLabelResource{}.ResourceType()
}

func (MssqlDatabaseSensitivityLabelListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"database_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: commonids.ValidateSqlDatabaseID},
				},
			},
		},
	}
}

func (MssqlDatabaseSensitivityLabelListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.MSSQL.SensitivityLabelsClient
	r := MssqlDatabaseSensitivityLabelResource{}

	var data MssqlDatabaseSensitivityLabelListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	databaseId, err := commonids.ParseSqlDatabaseID(data.DatabaseId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Database ID for `%s`", r.ResourceType()), err)
		return
	}

	resp, err := client.ListCurrentByDatabaseComplete(ctx, *databaseId, sensitivitylabels.DefaultListCurrentByDatabaseOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", r.ResourceType()), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)

			id, err := sensitivitylabels.ParseSensitivityLabelSourceIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Sensitivity Label ID", err)
				return
			}
			result.DisplayName = fmt.Sprintf("%s.%s.%s", id.SchemaName, id.TableName, id.ColumnName)

			rmd := sdk.NewResourceMetaData(metadata.Client, r)
			rmd.SetID(id)
			if err := r.flatten(rmd, id, &item); err != nil {
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
