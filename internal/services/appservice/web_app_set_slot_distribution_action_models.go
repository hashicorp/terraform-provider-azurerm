package appservice

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type webAppSetSlotDistributionActionModel struct {
	AppServiceId types.String                                                                  `tfsdk:"app_service_id"`
	SlotRule     typehelpers.ListNestedObjectValueOf[webAppSetSlotDistributionActionRuleModel] `tfsdk:"slot_rule"`
	Timeout      types.String                                                                  `tfsdk:"timeout"`
}

type webAppSetSlotDistributionActionRuleModel struct {
	Hostname                  types.String  `tfsdk:"hostname"`
	RuleName                  types.String  `tfsdk:"rule_name"`
	ReroutePercentage         types.Float64 `tfsdk:"reroute_percentage"`
	ChangeStep                types.Float64 `tfsdk:"change_step"`
	ChangeIntervalMinutes     types.Int64   `tfsdk:"change_interval_minutes"`
	MinReroutePercentage      types.Float64 `tfsdk:"min_reroute_percentage"`
	MaxReroutePercentage      types.Float64 `tfsdk:"max_reroute_percentage"`
	ChangeDecisionCallbackUrl types.String  `tfsdk:"change_decision_callback_url"`
}

func expandWebAppSetSlotDistributionActionModel(ctx context.Context, model webAppSetSlotDistributionActionModel, diags *diag.Diagnostics) *webapps.SiteConfigResource {
	rules, d := model.SlotRule.ToSlice(ctx)
	if d.HasError() {
		diags.Append(d...)
		return nil
	}

	rampupRulesEnvelope := make([]webapps.RampUpRule, 0, len(rules))
	for _, rule := range rules {
		rampupRulesEnvelope = append(rampupRulesEnvelope, expandWebAppSetSlotDistributionActionRuleModel(rule))
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

func expandWebAppSetSlotDistributionActionRuleModel(rule *webAppSetSlotDistributionActionRuleModel) webapps.RampUpRule {
	return webapps.RampUpRule{
		Name:                      rule.RuleName.ValueStringPointer(),
		ActionHostName:            rule.Hostname.ValueStringPointer(),
		ReroutePercentage:         rule.ReroutePercentage.ValueFloat64Pointer(),
		ChangeStep:                rule.ChangeStep.ValueFloat64Pointer(),
		ChangeIntervalInMinutes:   rule.ChangeIntervalMinutes.ValueInt64Pointer(),
		MinReroutePercentage:      rule.MinReroutePercentage.ValueFloat64Pointer(),
		MaxReroutePercentage:      rule.MaxReroutePercentage.ValueFloat64Pointer(),
		ChangeDecisionCallbackURL: rule.ChangeDecisionCallbackUrl.ValueStringPointer(),
	}
}
