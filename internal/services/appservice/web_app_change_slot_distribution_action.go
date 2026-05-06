package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ action.ActionWithValidateConfig = &WebAppSetSlotDistributionAction{}

type WebAppSetSlotDistributionAction struct {
	sdk.ActionMetadata
}

func newWebAppSetSlotDistributionAction() action.Action {
	return &WebAppSetSlotDistributionAction{}
}

type WebAppSetSlotDistributionActionModel struct {
	AppServiceId types.String                                                                  `tfsdk:"app_service_id"`
	SlotRule     typehelpers.ListNestedObjectValueOf[WebAppSetSlotDistributionActionRuleModel] `tfsdk:"slot_rule"`
	Timeout      types.String                                                                  `tfsdk:"timeout"`
}

type WebAppSetSlotDistributionActionRuleModel struct {
	Hostname                  types.String  `tfsdk:"hostname"`
	RuleName                  types.String  `tfsdk:"rule_name"`
	ReroutePercentage         types.Float64 `tfsdk:"reroute_percentage"`
	ChangeStep                types.Float64 `tfsdk:"change_step"`
	ChangeIntervalMinutes     types.Int64   `tfsdk:"change_interval_minutes"`
	MinReroutePercentage      types.Float64 `tfsdk:"min_reroute_percentage"`
	MaxReroutePercentage      types.Float64 `tfsdk:"max_reroute_percentage"`
	ChangeDecisionCallbackUrl types.String  `tfsdk:"change_decision_callback_url"`
}

func (*WebAppSetSlotDistributionAction) Schema(ctx context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
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
			"slot_rule": webAppSetSlotDistributionActionRuleSchema(ctx),
		},
	}
}

func webAppSetSlotDistributionActionRuleSchema(ctx context.Context) schema.Block {
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
					Validators: []validator.Float64{
						float64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("change_interval_minutes")),
					},
				},
				"change_interval_minutes": schema.Int64Attribute{
					Optional:    true,
					Description: "Specifies interval in minutes to reevaluate ReroutePercentage.",
					Validators: []validator.Int64{
						int64validator.AlsoRequires(path.MatchRelative().AtParent().AtName("change_step")),
					},
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
		CustomType: typehelpers.NewListNestedObjectTypeOf[WebAppSetSlotDistributionActionRuleModel](ctx),
	}
}

func (a *WebAppSetSlotDistributionAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_web_app_set_slot_distribution"
}

func (*WebAppSetSlotDistributionAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	data := WebAppSetSlotDistributionActionModel{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rules, d := data.SlotRule.ToSlice(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Checking two different validations by looping through all rules:
	// 1. duplicate slot hostnames
	// 2. sum of percentage > 100
	foundHostnameDupe := false
	uniqueHostnameMap := make(map[string]bool)
	totalRulePercentage := float64(0)
	for _, rule := range rules {
		// NOTE: there are scenarios where values are Unknown when validation is run.
		// If any values are Unknown they will not be used for both of these validations.
		if !rule.ReroutePercentage.IsUnknown() {
			totalRulePercentage += rule.ReroutePercentage.ValueFloat64()
		}

		if !foundHostnameDupe && !rule.Hostname.IsUnknown() {
			if !uniqueHostnameMap[rule.Hostname.ValueString()] {
				uniqueHostnameMap[rule.Hostname.ValueString()] = true
			} else {
				foundHostnameDupe = true
			}
		}
	}

	if foundHostnameDupe {
		resp.Diagnostics.AddAttributeError(path.Root("slot_rule").AtName("hostname"), "Multiple slot rules have the same target hostname", "The target `hostname` value of each slot rule must be unique across all provided rules.")
	}

	if totalRulePercentage > 100 {
		resp.Diagnostics.AddAttributeError(path.Root("slot_rule").AtName("reroute_percentage"), "Total Reroute Percentage greater than 100%", "The total percentage of all slot rules should not be greater than 100%.")
	}
}

func (a *WebAppSetSlotDistributionAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := a.Client.AppService.WebAppsClient

	model := WebAppSetSlotDistributionActionModel{}

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

func (a *WebAppSetSlotDistributionAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	a.Defaults(ctx, request, response)
}

func (a *WebAppSetSlotDistributionAction) expand(ctx context.Context, model WebAppSetSlotDistributionActionModel, diags *diag.Diagnostics) *webapps.SiteConfigResource {
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

func (*WebAppSetSlotDistributionAction) expandSlotDistributionRules(ctx context.Context, input typehelpers.ListNestedObjectValueOf[WebAppSetSlotDistributionActionRuleModel], diags *diag.Diagnostics) []webapps.RampUpRule {
	rules, d := input.ToSlice(ctx)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	rampupRules := make([]webapps.RampUpRule, 0, len(rules))

	for _, rule := range rules {
		rampUpRule := webapps.RampUpRule{
			Name:                      rule.RuleName.ValueStringPointer(),
			ActionHostName:            rule.Hostname.ValueStringPointer(),
			ReroutePercentage:         rule.ReroutePercentage.ValueFloat64Pointer(),
			ChangeStep:                rule.ChangeStep.ValueFloat64Pointer(),
			ChangeIntervalInMinutes:   rule.ChangeIntervalMinutes.ValueInt64Pointer(),
			MinReroutePercentage:      rule.MinReroutePercentage.ValueFloat64Pointer(),
			MaxReroutePercentage:      rule.MaxReroutePercentage.ValueFloat64Pointer(),
			ChangeDecisionCallbackURL: rule.ChangeDecisionCallbackUrl.ValueStringPointer(),
		}

		rampupRules = append(rampupRules, rampUpRule)
	}

	return rampupRules
}
