// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AlertPrometheusRuleGroupResourceModel struct {
	Name              string                `tfschema:"name"`
	Location          string                `tfschema:"location"`
	ResourceGroupName string                `tfschema:"resource_group_name"`
	Rule              []PrometheusRuleModel `tfschema:"rule"`
	Scopes            []string              `tfschema:"scopes"`
	ClusterName       string                `tfschema:"cluster_name"`
	Description       string                `tfschema:"description"`
	Interval          string                `tfschema:"interval"`
	RuleGroupEnabled  bool                  `tfschema:"rule_group_enabled"`
	Tags              map[string]string     `tfschema:"tags"`
}

type PrometheusRuleModel struct {
	Action          []PrometheusRuleGroupActionModel     `tfschema:"action"`
	Alert           string                               `tfschema:"alert"`
	Annotations     map[string]string                    `tfschema:"annotations"`
	Enabled         bool                                 `tfschema:"enabled"`
	Expression      string                               `tfschema:"expression"`
	For             string                               `tfschema:"for"`
	Labels          map[string]string                    `tfschema:"labels"`
	Record          string                               `tfschema:"record"`
	AlertResolution []PrometheusRuleAlertResolutionModel `tfschema:"alert_resolution"`
	Severity        int64                                `tfschema:"severity"`
}

type PrometheusRuleGroupActionModel struct {
	ActionGroupId    string            `tfschema:"action_group_id"`
	ActionProperties map[string]string `tfschema:"action_properties"`
}

type PrometheusRuleAlertResolutionModel struct {
	AutoResolved  bool   `tfschema:"auto_resolved"`
	TimeToResolve string `tfschema:"time_to_resolve"`
}

type AlertPrometheusRuleGroupResource struct{}

var _ sdk.ResourceWithUpdate = AlertPrometheusRuleGroupResource{}

func (r AlertPrometheusRuleGroupResource) ResourceType() string {
	return "azurerm_monitor_alert_prometheus_rule_group"
}

func (r AlertPrometheusRuleGroupResource) ModelObject() interface{} {
	return &AlertPrometheusRuleGroupResourceModel{}
}

func (r AlertPrometheusRuleGroupResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AlertPrometheusRuleGroupResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			for i, r := range model.Rule {
				if (r.Alert != "" && r.Record != "") || (r.Alert == "" && r.Record == "") {
					return fmt.Errorf("one and only one of [rule.%d.record, rule.%d.alert] for %s must be set", i, i, model.Name)
				}

				// actions, severity, annotations, for, alert_resolution must be empty when type is recording rule
				if r.Record != "" {
					_, actionOk := metadata.ResourceDiff.GetOk("rule." + strconv.Itoa(i) + ".action")
					_, severityOk := metadata.ResourceDiff.GetOk("rule." + strconv.Itoa(i) + ".severity")
					_, annotationsOk := metadata.ResourceDiff.GetOk("rule." + strconv.Itoa(i) + ".annotations")
					_, forOk := metadata.ResourceDiff.GetOk("rule." + strconv.Itoa(i) + ".for")
					_, resolveConfigurationOk := metadata.ResourceDiff.GetOk("rule." + strconv.Itoa(i) + ".alert_resolution")

					if actionOk || severityOk || annotationsOk || forOk || resolveConfigurationOk {
						return fmt.Errorf("the rule.[%d].action, rule.[%d].severity, rule.[%d].annotations, rule.[%d].for and rule.[%d].alert_resolution must be empty when the rule type of Alert Prometheus Rule Group (%s) is record", i, i, i, i, i, model.Name)
					}
				}
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r AlertPrometheusRuleGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return prometheusrulegroups.ValidatePrometheusRuleGroupID
}

func (r AlertPrometheusRuleGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 260),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"rule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 5,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action_group_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"action_properties": {
									Type:     pluginsdk.TypeMap,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"alert": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"annotations": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"expression": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"for": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"labels": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: schema.TypeString,
						},
					},

					"record": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"alert_resolution": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"auto_resolved": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"time_to_resolve": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"severity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 4),
					},
				},
			},
		},

		"scopes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"interval": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"rule_group_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r AlertPrometheusRuleGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AlertPrometheusRuleGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AlertPrometheusRuleGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := prometheusrulegroups.NewPrometheusRuleGroupID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := prometheusrulegroups.PrometheusRuleGroupResource{
				Location: location.Normalize(model.Location),
				Properties: prometheusrulegroups.PrometheusRuleGroupProperties{
					Enabled: pointer.To(model.RuleGroupEnabled),
					Scopes:  model.Scopes,
				},
				Tags: pointer.To(model.Tags),
			}

			properties.Properties.ClusterName = pointer.To(model.ClusterName)
			properties.Properties.Description = pointer.To(model.Description)
			if _, ok := metadata.ResourceData.GetOk("interval"); ok {
				properties.Properties.Interval = pointer.To(model.Interval)
			}
			properties.Properties.Rules = expandPrometheusRuleModel(model.Rule, metadata.ResourceData)

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AlertPrometheusRuleGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient

			id, err := prometheusrulegroups.ParsePrometheusRuleGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model AlertPrometheusRuleGroupResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("cluster_name") {
				properties.Properties.ClusterName = pointer.To(model.ClusterName)
			}
			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = pointer.To(model.Description)
			}
			if metadata.ResourceData.HasChange("rule_group_enabled") {
				properties.Properties.Enabled = pointer.To(model.RuleGroupEnabled)
			}
			if metadata.ResourceData.HasChange("interval") {
				properties.Properties.Interval = pointer.To(model.Interval)
			}
			if metadata.ResourceData.HasChange("rule") {
				properties.Properties.Rules = expandPrometheusRuleModel(model.Rule, metadata.ResourceData)
			}
			if metadata.ResourceData.HasChange("scopes") {
				properties.Properties.Scopes = model.Scopes
			}
			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AlertPrometheusRuleGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient

			id, err := prometheusrulegroups.ParsePrometheusRuleGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := AlertPrometheusRuleGroupResourceModel{
				Name:              id.PrometheusRuleGroupName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if resp.Model != nil {
				state.ClusterName = pointer.From(resp.Model.Properties.ClusterName)
				state.Description = pointer.From(resp.Model.Properties.Description)
				state.Interval = pointer.From(resp.Model.Properties.Interval)
				state.Location = location.Normalize(resp.Model.Location)
				state.Rule = flattenPrometheusRuleModel(&resp.Model.Properties.Rules)
				state.RuleGroupEnabled = pointer.From(resp.Model.Properties.Enabled)
				state.Scopes = resp.Model.Properties.Scopes
				state.Tags = pointer.From(resp.Model.Tags)
			}
			return metadata.Encode(&state)
		},
	}
}

func (r AlertPrometheusRuleGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.AlertPrometheusRuleGroupClient

			id, err := prometheusrulegroups.ParsePrometheusRuleGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandPrometheusRuleModel(inputList []PrometheusRuleModel, d *schema.ResourceData) []prometheusrulegroups.PrometheusRule {
	outputList := make([]prometheusrulegroups.PrometheusRule, 0)

	for i, v := range inputList {
		output := prometheusrulegroups.PrometheusRule{
			Enabled:    pointer.To(v.Enabled),
			Expression: v.Expression,
			Labels:     pointer.To(v.Labels),
		}

		if v.Alert != "" {
			output.Actions = expandPrometheusRuleGroupActionModel(v.Action)
			output.Alert = pointer.To(v.Alert)
			if _, ok := d.GetOk("rule." + strconv.Itoa(i) + ".severity"); ok {
				output.Severity = pointer.To(v.Severity)
			}
			output.Annotations = pointer.To(v.Annotations)
			output.For = pointer.To(v.For)
			output.ResolveConfiguration = expandPrometheusRuleAlertResolutionModel(v.AlertResolution)
		} else {
			// action, alert, severity, annotations, for, alert_resolution must be empty when type is recording rule
			output.Record = pointer.To(v.Record)
		}
		outputList = append(outputList, output)
	}

	return outputList
}

func expandPrometheusRuleGroupActionModel(inputList []PrometheusRuleGroupActionModel) *[]prometheusrulegroups.PrometheusRuleGroupAction {
	outputList := make([]prometheusrulegroups.PrometheusRuleGroupAction, 0)
	for _, v := range inputList {
		output := prometheusrulegroups.PrometheusRuleGroupAction{
			ActionProperties: pointer.To(v.ActionProperties),
		}
		output.ActionGroupId = pointer.To(v.ActionGroupId)
		outputList = append(outputList, output)
	}

	return &outputList
}

func expandPrometheusRuleAlertResolutionModel(inputList []PrometheusRuleAlertResolutionModel) *prometheusrulegroups.PrometheusRuleResolveConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := prometheusrulegroups.PrometheusRuleResolveConfiguration{
		AutoResolved: pointer.To(input.AutoResolved),
	}
	output.TimeToResolve = pointer.To(input.TimeToResolve)

	return &output
}

func flattenPrometheusRuleModel(inputList *[]prometheusrulegroups.PrometheusRule) []PrometheusRuleModel {
	outputList := make([]PrometheusRuleModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := PrometheusRuleModel{
			Expression: input.Expression,
		}

		actionsValue := flattenPrometheusRuleGroupActionModel(input.Actions)
		output.Action = actionsValue
		output.Alert = pointer.From(input.Alert)
		output.Annotations = pointer.From(input.Annotations)
		output.Enabled = pointer.From(input.Enabled)
		output.For = pointer.From(input.For)
		output.Labels = pointer.From(input.Labels)
		output.Record = pointer.From(input.Record)
		resolveConfigurationValue := flattenPrometheusRuleAlertResolutionModel(input.ResolveConfiguration)
		output.AlertResolution = resolveConfigurationValue
		output.Severity = pointer.From(input.Severity)
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenPrometheusRuleGroupActionModel(inputList *[]prometheusrulegroups.PrometheusRuleGroupAction) []PrometheusRuleGroupActionModel {
	outputList := make([]PrometheusRuleGroupActionModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := PrometheusRuleGroupActionModel{}
		output.ActionGroupId = pointer.From(input.ActionGroupId)
		output.ActionProperties = pointer.From(input.ActionProperties)
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenPrometheusRuleAlertResolutionModel(input *prometheusrulegroups.PrometheusRuleResolveConfiguration) []PrometheusRuleAlertResolutionModel {
	outputList := make([]PrometheusRuleAlertResolutionModel, 0)
	if input == nil {
		return outputList
	}

	output := PrometheusRuleAlertResolutionModel{}
	output.AutoResolved = pointer.From(input.AutoResolved)
	output.TimeToResolve = pointer.From(input.TimeToResolve)

	return append(outputList, output)
}
