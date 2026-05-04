package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	responsehelper "github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WebAppSlotDistributionAction struct {
	sdk.ActionMetadata
}

func newWebAppSlotDistributionAction() action.Action {
	return &WebAppSlotDistributionAction{}
}

type WebAppSlotDistributionActionModel struct {
	AppServiceId types.String                                                               `tfsdk:"app_service_id"`
	SlotRule     typehelpers.ListNestedObjectValueOf[WebAppSlotDistributionActionRuleModel] `tfsdk:"slot_rule"`
	Timeout      types.String                                                               `tfsdk:"timeout"`
}

type WebAppSlotDistributionActionRuleModel struct {
	Hostname                  types.String  `tfsdk:"hostname"`
	RuleName                  types.String  `tfsdk:"rule_name"`
	ReroutePercentage         types.Float64 `tfsdk:"reroute_percentage"`
	ChangeStep                types.Float64 `tfsdk:"change_step"`
	ChangeIntervalMinutes     types.Int32   `tfsdk:"change_interval_minutes"`
	MinReroutePercentage      types.Float64 `tfsdk:"min_reroute_percentage"`
	MaxReroutePercentage      types.Float64 `tfsdk:"max_reroute_percentage"`
	ChangeDecisionCallbackUrl types.String  `tfsdk:"change_decision_callback_url"`
}

var _ sdk.Action = &WebAppSlotDistributionAction{}

func (*WebAppSlotDistributionAction) Schema(ctx context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"app_service_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateWebAppID,
					},
				},
			},
			"timeout": schema.StringAttribute{
				Optional:            true,
				Description:         "Timeout duration for the Slot Distribution action to complete. Defaults to 20m.",
				MarkdownDescription: "Timeout duration for the Slot Distribution action to complete. Defaults to 20m.",
			},
		},
		Blocks: map[string]schema.Block{
			"slot_rule": webAppSlotDistributionActionRuleSchema(ctx),
		},
	}
}

func webAppSlotDistributionActionRuleSchema(ctx context.Context) schema.Block {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"hostname": schema.StringAttribute{
					Required:    true,
					Description: "The Hostname of the Slot to route traffic to.",
				},
				"rule_name": schema.StringAttribute{
					Required:    true,
					Description: "The name of the unique distribution rule, if not supplied will default to the slot name (recommended).",
				},
				"reroute_percentage": schema.Float64Attribute{
					Required:    true,
					Description: "Percentage of traffic which will be redirected to the supplied slot hostname.",
				},
				"change_step": schema.Float64Attribute{
					Optional:    true,
					Description: "In auto ramp up scenario this is the step to add/remove from ReroutePercentage.",
				},
				"change_interval_minutes": schema.Int32Attribute{
					Optional:    true,
					Description: "Specifies interval in minutes to reevaluate ReroutePercentage.",
				},
				"min_reroute_percentage": schema.Float64Attribute{
					Optional:    true,
					Description: "Specifies lower boundary above which ReroutePercentage will stay.",
				},
				"max_reroute_percentage": schema.Float64Attribute{
					Optional:    true,
					Description: "Specifies upper boundary below which ReroutePercentage will stay.",
				},
				"change_decision_callback_url": schema.StringAttribute{
					Optional:    true,
					Description: "Custom decision algorithm can be provided in TiPCallback site extension which URL can be specified.",
					Validators: []validator.String{
						typehelpers.WrappedStringValidator{
							Func: validation.IsURLWithHTTPorHTTPS,
						},
					},
				},
			},
		},
		CustomType: typehelpers.NewListNestedObjectTypeOf[WebAppSlotDistributionActionRuleModel](ctx),
	}
}

func (a *WebAppSlotDistributionAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_web_app_slot_distribution"
}

func (a *WebAppSlotDistributionAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := a.Client.AppService.WebAppsClient

	model := WebAppSlotDistributionActionModel{}

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

	appId, err := commonids.ParseWebAppID(model.AppServiceId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, "id parsing error", err)
		return
	}

	app, err := client.Get(ctx, *appId)
	if err != nil {
		if responsehelper.WasNotFound(app.HttpResponse) {
			sdk.SetResponseErrorDiagnostic(response, fmt.Sprintf("web app '%s' not found", model.AppServiceId), err)
			return
		}
		sdk.SetResponseErrorDiagnostic(response, fmt.Sprintf("error retrieving web app '%s'", model.AppServiceId), err)
		return
	}

	siteConfigEnvelope := a.expand(ctx, model, &response.Diagnostics)
	if response.Diagnostics.HasError() || siteConfigEnvelope == nil {
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Assigning new Slot Distribution Rules for AppServiceId '%s'", model.AppServiceId),
	})

	_, err = client.UpdateConfiguration(ctx, *appId, *siteConfigEnvelope)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(response, fmt.Sprintf("error updating configuration for AppServiceId '%s'", model.AppServiceId), err)
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Completed setting of Slot Distribution Rules for AppServiceId '%s'", model.AppServiceId),
	})
}

func (a *WebAppSlotDistributionAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	a.Defaults(ctx, request, response)
}

func (a *WebAppSlotDistributionAction) expand(ctx context.Context, model WebAppSlotDistributionActionModel, diags *diag.Diagnostics) *webapps.SiteConfigResource {
	rampupRulesEnvelope := a.expandSlotDistributionRules(ctx, model.SlotRule, diags)
	if diags.HasError() {
		return nil
	}

	siteConfigEnvelope := webapps.SiteConfigResource{
		Properties: &webapps.SiteConfig{
			Experiments: &webapps.Experiments{
				RampUpRules: &rampupRulesEnvelope,
			},
		},
	}

	return &siteConfigEnvelope
}

func (*WebAppSlotDistributionAction) expandSlotDistributionRules(ctx context.Context, input typehelpers.ListNestedObjectValueOf[WebAppSlotDistributionActionRuleModel], diags *diag.Diagnostics) []webapps.RampUpRule {
	rules, d := input.ToSlice(ctx)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	rampupRules := make([]webapps.RampUpRule, 0, len(rules))

	for _, rule := range rules {
		rampUpRule := webapps.RampUpRule{
			Name:              rule.RuleName.ValueStringPointer(),
			ActionHostName:    rule.Hostname.ValueStringPointer(),
			ReroutePercentage: rule.ReroutePercentage.ValueFloat64Pointer(),
		}

		rampupRules = append(rampupRules, rampUpRule)
	}

	return rampupRules
}
