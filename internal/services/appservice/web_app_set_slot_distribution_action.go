package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ action.ActionWithConfigure      = &webAppSetSlotDistributionAction{}
	_ action.ActionWithValidateConfig = &webAppSetSlotDistributionAction{}
)

type webAppSetSlotDistributionAction struct {
	sdk.ActionMetadata
}

func newWebAppSetSlotDistributionAction() action.Action {
	return &webAppSetSlotDistributionAction{}
}

func (*webAppSetSlotDistributionAction) Schema(ctx context.Context, _ action.SchemaRequest, response *action.SchemaResponse) {
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
		CustomType: typehelpers.NewListNestedObjectTypeOf[webAppSetSlotDistributionActionRuleModel](ctx),
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
	}
}

func (a *webAppSetSlotDistributionAction) Metadata(_ context.Context, _ action.MetadataRequest, response *action.MetadataResponse) {
	response.TypeName = "azurerm_web_app_set_slot_distribution"
}

func (*webAppSetSlotDistributionAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	data := webAppSetSlotDistributionActionModel{}
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
	// 1. duplicate slot rule names
	// 2. duplicate slot hostnames
	// 3. sum of percentage > 100
	foundHostnameDupe := false
	foundRulenameDupe := false
	uniqueRulenameMap := make(map[string]bool)
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

		if !foundRulenameDupe && !rule.RuleName.IsUnknown() {
			if !uniqueRulenameMap[rule.RuleName.ValueString()] {
				uniqueRulenameMap[rule.RuleName.ValueString()] = true
			} else {
				foundRulenameDupe = true
			}
		}
	}

	if foundHostnameDupe {
		resp.Diagnostics.AddAttributeError(path.Root("slot_rule").AtName("hostname"), "Multiple slot rules have the same target hostname", "The target `hostname` value of each slot rule must be unique across all provided rules.")
	}

	if foundRulenameDupe {
		resp.Diagnostics.AddAttributeError(path.Root("slot_rule").AtName("rule_name"), "Multiple slot rules have the same name", "The `rule_name` value of each slot rule must be unique across all provided rules.")
	}

	if totalRulePercentage > 100 {
		resp.Diagnostics.AddAttributeError(path.Root("slot_rule").AtName("reroute_percentage"), "Total Reroute Percentage greater than 100%", "The total percentage of all slot rules should not be greater than 100%.")
	}
}

func (a *webAppSetSlotDistributionAction) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	client := a.Client.AppService.WebAppsClient

	model := webAppSetSlotDistributionActionModel{}

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

	siteConfigEnvelope := expandWebAppSetSlotDistributionActionModel(ctx, model, &response.Diagnostics)
	if response.Diagnostics.HasError() || siteConfigEnvelope == nil {
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Assigning new Slot Distribution Rules for AppServiceId '%s'", model.AppServiceId),
	})

	if _, err := client.UpdateConfiguration(ctx, *appId, *siteConfigEnvelope); err != nil {
		sdk.SetResponseErrorDiagnostic(response, fmt.Sprintf("error updating configuration for AppServiceId '%s'", model.AppServiceId), err)
		return
	}

	response.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Completed setting of Slot Distribution Rules for AppServiceId '%s'", model.AppServiceId),
	})
}

func (a *webAppSetSlotDistributionAction) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	a.Defaults(ctx, request, response)
}
