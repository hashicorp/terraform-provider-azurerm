package cdn

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/convert"
	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-06-01/afdendpoints"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type AzureFrontDoorCachePurgeAction struct {
	sdk.ActionMetadata
}

type AzureFrontDoorCachePurgeActionModel struct {
	AzureFrontDoorId types.String                          `tfsdk:"front_door_id"`
	ContentPaths     typehelpers.ListValueOf[types.String] `tfsdk:"content_paths"`
	Domains          typehelpers.ListValueOf[types.String] `tfsdk:"domains"`
}

var _ sdk.Action = &AzureFrontDoorCachePurgeAction{}

func newAzureFrontDoorCachePurgeAction() action.Action {
	return &AzureFrontDoorCachePurgeAction{}
}

func (a *AzureFrontDoorCachePurgeAction) Schema(ctx context.Context, request action.SchemaRequest, response *action.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"front_door_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: afdendpoints.ValidateAfdEndpointID,
					},
				},
			},

			"content_paths": schema.ListAttribute{
				CustomType:  typehelpers.NewListTypeOf[types.String](ctx),
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.All(
						listvalidator.NoNullValues(),
						listvalidator.ValueStringsAre(
							stringvalidator.LengthAtLeast(1),
						),
					),
				},
			},

			"domains": schema.ListAttribute{
				CustomType:  typehelpers.NewListTypeOf[types.String](ctx),
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.All(
						listvalidator.NoNullValues(),
						listvalidator.ValueStringsAre(
							stringvalidator.LengthAtLeast(1),
						),
					),
				},
			},
		},
	}
}

func (a *AzureFrontDoorCachePurgeAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_azure_front_door_cache_purge"
}

func (a *AzureFrontDoorCachePurgeAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := a.Client.Cdn.AFDEndpointsClient

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	model := AzureFrontDoorCachePurgeActionModel{}

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}

	id, err := afdendpoints.ParseAfdEndpointID(model.AzureFrontDoorId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, "id parsing error", err)
		return
	}

	cp := make([]string, 0)
	convert.Expand(ctx, model.ContentPaths, &cp, &response.Diagnostics)
	if response.Diagnostics.HasError() {
		return
	}

	payload := afdendpoints.AfdPurgeParameters{
		ContentPaths: cp,
	}

	if len(model.Domains.Elements()) > 0 {
		d := make([]string, 0)
		convert.Expand(ctx, model.Domains, &d, &response.Diagnostics)
		if response.Diagnostics.HasError() {
			return
		}

		payload.Domains = pointer.To(d)
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: "Purging Azure Front Door Cache",
	})

	if err := client.PurgeContentThenPoll(ctx, *id, payload); err != nil {
		sdk.SetResponseErrorDiagnostic(response, fmt.Sprintf("purging contents for %s", id), err)
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: "Purged Azure Front Door Cache",
	})
}

func (a *AzureFrontDoorCachePurgeAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	a.Defaults(ctx, request, response)
}
