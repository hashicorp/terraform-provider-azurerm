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

type FrontDoorCachePurgeAction struct {
	sdk.ActionMetadata
}

type FrontDoorCachePurgeActionModel struct {
	AzureFrontDoorId types.String                          `tfsdk:"front_door_endpoint_id"`
	ContentPaths     typehelpers.ListValueOf[types.String] `tfsdk:"content_paths"`
	Domains          typehelpers.ListValueOf[types.String] `tfsdk:"domains"`
	Timeout          types.String                          `tfsdk:"timeout"`
}

// var _ sdk.Action = &FrontDoorCachePurgeAction{}
func newCDNFrontDoorCachePurgeAction() action.Action {
	return &FrontDoorCachePurgeAction{}
}

func (a *FrontDoorCachePurgeAction) Schema(ctx context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"front_door_endpoint_id": schema.StringAttribute{
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

			"timeout": schema.StringAttribute{
				Optional:            true,
				Description:         "Timeout duration for the Front Door Purge action to complete. Defaults to 30m.",
				MarkdownDescription: "Timeout duration for the Front Door Purge action to complete. Defaults to 30m.",
			},
		},
	}
}

func (a *FrontDoorCachePurgeAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_cdn_front_door_cache_purge"
}

func (a *FrontDoorCachePurgeAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := a.Client.Cdn.AFDEndpointsClient

	model := FrontDoorCachePurgeActionModel{}
	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)
	if response.Diagnostics.HasError() {
		return
	}

	ctxTimeout := 30 * time.Minute
	if t := model.Timeout; !t.IsNull() {
		timeout, err := time.ParseDuration(t.ValueString())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(response, "parsing `timeout`", err)
			return
		}
		ctxTimeout = timeout
	}

	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

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
		Message: fmt.Sprintf("Purging Azure Front Door Cache, paths (%s) / domains (%s)", model.ContentPaths, model.Domains),
	})

	if err := client.PurgeContentThenPoll(ctx, *id, payload); err != nil {
		sdk.SetResponseErrorDiagnostic(response, fmt.Sprintf("purging contents for %s", id), err)
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Purged Azure Front Door Cache, paths (%s) / domains (%s)", model.ContentPaths, model.Domains),
	})
}

func (a *FrontDoorCachePurgeAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	a.Defaults(ctx, request, response)
}
