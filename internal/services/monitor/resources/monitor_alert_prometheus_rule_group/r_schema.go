// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_prometheus_rule_group

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
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

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Alert Management Prometheus Rule Group.",
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
									Description:  "The resource ID of the monitor action group.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"action_properties": {
									Description: "The properties of an action group object.",
									Type:        pluginsdk.TypeMap,
									Optional:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"alert": {
						Description:  "The Alert rule name.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"annotations": {
						Description: "A set of informational labels that can be used to store longer additional information such as alert descriptions or runbook links.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"enabled": {
						Description: "Whether this rule is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"expression": {
						Description:  "The Prometheus Query Language expression to evaluate. For more details see [this doc](https://prometheus.io/docs/prometheus/latest/querying/basics). Evaluate at the period given by `interval` and record the result as a new set of time series with the metric name given by `record`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"for": {
						Description:  "The amount of time alert must be active before firing, represented in ISO 8601 duration format.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"labels": {
						Description: "The labels to add or overwrite before storing the result.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: schema.TypeString,
						},
					},

					"record": {
						Description:  "The recorded metrics name.",
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
									Description: "Whether the alert is automatically resolved.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
								},

								"time_to_resolve": {
									Description:  "The alert auto-resolution interval, represented in ISO 8601 duration format.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"severity": {
						Description:  "The severity of the alerts fired by the rule.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 4),
					},
				},
			},
		},

		"scopes": {
			Description: "The resource ID of the Azure Monitor Workspace.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"cluster_name": {
			Description:  "The name of the Managed Kubernetes Cluster.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Description:  "The description of the Alert Management Prometheus Rule Group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"interval": {
			Description:  "The interval in which to run the Alert Management Prometheus Rule Group represented in ISO 8601 duration format.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"rule_group_enabled": {
			Description: "Whether this Alert Management Prometheus Rule Group is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
