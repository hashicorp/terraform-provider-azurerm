// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynatrace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type TagRulesResource struct{}

var _ sdk.ResourceWithUpdate = TagRulesResource{}

type TagRulesResourceModel struct {
	Name        string       `tfschema:"name"`
	Monitor     string       `tfschema:"monitor_id"`
	LogRules    []LogRule    `tfschema:"log_rule"`
	MetricRules []MetricRule `tfschema:"metric_rule"`
}

type MetricRule struct {
	FilteringTags []FilteringTag `tfschema:"filtering_tag"`
}

type LogRule struct {
	FilteringTags        []FilteringTag `tfschema:"filtering_tag"`
	SendAadLogs          bool           `tfschema:"send_azure_active_directory_logs_enabled"`
	SendActivityLogs     bool           `tfschema:"send_activity_logs_enabled"`
	SendSubscriptionLogs bool           `tfschema:"send_subscription_logs_enabled"`
}

type FilteringTag struct {
	Name   string `tfschema:"name"`
	Value  string `tfschema:"value"`
	Action string `tfschema:"action"`
}

func (r TagRulesResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"monitor_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: monitors.ValidateMonitorID,
		},

		"log_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"send_azure_active_directory_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"send_activity_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"send_subscription_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"filtering_tag": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Include",
										"Exclude",
									}, false),
								},

								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"metric_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"filtering_tag": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Include",
										"Exclude",
									}, false),
								},

								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r TagRulesResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r TagRulesResource) ModelObject() interface{} {
	return &TagRulesResourceModel{}
}

func (r TagRulesResource) ResourceType() string {
	return "azurerm_dynatrace_tag_rules"
}

func (r TagRulesResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.TagRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model TagRulesResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			monitorId, err := monitors.ParseMonitorID(model.Monitor)
			id := tagrules.NewTagRuleID(subscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, model.Name)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			tagRulesProps := tagrules.MonitoringTagRulesProperties{
				LogRules:    ExpandLogRule(model.LogRules),
				MetricRules: ExpandMetricRules(model.MetricRules),
			}
			tagRules := tagrules.TagRule{
				Name:       &model.Name,
				Properties: tagRulesProps,
			}

			if _, err := client.CreateOrUpdate(ctx, id, tagRules); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r TagRulesResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.TagRulesClient
			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if model := resp.Model; model != nil {
				props := model.Properties
				monitorId := monitors.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName)

				state := TagRulesResourceModel{
					Name:        id.TagRuleName,
					Monitor:     monitorId.ID(),
					LogRules:    FlattenLogRules(props.LogRules),
					MetricRules: FlattenMetricRules(props.MetricRules),
				}

				return metadata.Encode(&state)
			}

			return nil
		},
	}
}

func (r TagRulesResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.TagRulesClient
			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r TagRulesResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Dynatrace.TagRulesClient
			id, err := tagrules.ParseTagRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state TagRulesResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := existing.Model

			if metadata.ResourceData.HasChange("metric_rule") {
				model.Properties.MetricRules = ExpandMetricRules(state.MetricRules)
			}

			if metadata.ResourceData.HasChange("log_rule") {
				model.Properties.LogRules = ExpandLogRule(state.LogRules)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r TagRulesResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tagrules.ValidateTagRuleID
}
