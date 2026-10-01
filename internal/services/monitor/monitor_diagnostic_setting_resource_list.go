// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MonitorDiagnosticSettingListResource struct{}

type MonitorDiagnosticSettingListModel struct {
	TargetResourceId types.String `tfsdk:"target_resource_id"`
}

var _ sdk.FrameworkListWrappedResource = new(MonitorDiagnosticSettingListResource)

func (MonitorDiagnosticSettingListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_monitor_diagnostic_setting"
}

func (MonitorDiagnosticSettingListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceMonitorDiagnosticSetting()
}

func (MonitorDiagnosticSettingListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target_resource_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateScopeID,
					},
				},
			},
		},
	}
}

func (MonitorDiagnosticSettingListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Monitor.DiagnosticSettingsClient

	var data MonitorDiagnosticSettingListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	targetResourceId := data.TargetResourceId.ValueString()
	scopeId := commonids.NewScopeID(targetResourceId)

	resp, err := client.ListComplete(ctx, scopeId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing diagnostic settings for `%s`", targetResourceId), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rawId := pointer.From(item.Id)
			if !strings.HasPrefix(rawId, "/") { // ListResponse inexplicably trims leading slash in the response?
				rawId = "/" + rawId
			}

			id, err := diagnosticsettings.ParseScopedDiagnosticSettingIDInsensitively(rawId)
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Diagnostic Setting ID", err)
				return
			}

			rd := resourceMonitorDiagnosticSetting().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceMonitorDiagnosticSettingFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "encoding `azurerm_monitor_diagnostic_setting` resource data", err)
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
