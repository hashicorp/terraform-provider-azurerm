// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
)

var _ sdk.EphemeralResource = &FunctionAppFunctionKeysEphemeralResource{}

func NewFunctionAppFunctionKeysEphemeralResource() ephemeral.EphemeralResource {
	return &FunctionAppFunctionKeysEphemeralResource{}
}

type FunctionAppFunctionKeysEphemeralResource struct {
	sdk.EphemeralResourceMetadata
}

type FunctionAppFunctionKeysEphemeralResourceModel struct {
	FunctionAppID types.String `tfsdk:"function_app_id"`
	FunctionName  types.String `tfsdk:"function_name"`
	DefaultKey    types.String `tfsdk:"default_key"`
	Keys          types.Map    `tfsdk:"keys"`
}

func (e *FunctionAppFunctionKeysEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "azurerm_function_app_function_keys"
}

func (e *FunctionAppFunctionKeysEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	e.Defaults(req, resp)
}

func (e *FunctionAppFunctionKeysEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"function_app_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Function App in which the Function resides.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateFunctionAppID,
					},
				},
			},

			"function_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Function.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validate.FunctionAppFunctionName,
					},
				},
			},

			"default_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The value of the `default` Function key.",
			},

			"keys": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Sensitive:   true,
				Description: "A map of all Function keys, where the map key is the name of the Function key and the value is the key itself.",
			},
		},
	}
}

func (e *FunctionAppFunctionKeysEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	client := e.Client.AppService.WebAppsClient
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	var data FunctionAppFunctionKeysEphemeralResourceModel
	if ok := e.DecodeOpen(ctx, req, resp, &data); !ok {
		return
	}

	appId, err := commonids.ParseFunctionAppID(data.FunctionAppID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing `function_app_id`", err)
		return
	}

	id := webapps.NewFunctionID(appId.SubscriptionId, appId.ResourceGroupName, appId.SiteName, data.FunctionName.ValueString())

	keys, err := listFunctionAppFunctionKeys(ctx, client, id)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "retrieving Function keys", err)
		return
	}

	keysValue, diags := types.MapValueFrom(ctx, types.StringType, keys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Keys = keysValue
	data.DefaultKey = types.StringNull()
	if v, ok := keys[functionAppFunctionDefaultKeyName]; ok {
		data.DefaultKey = types.StringValue(v)
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
