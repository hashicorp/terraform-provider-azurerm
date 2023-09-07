// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorMetricAlert() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorMetricAlertCreateUpdate,
		Read:   resourceMonitorMetricAlertRead,
		Update: resourceMonitorMetricAlertCreateUpdate,
		Delete: resourceMonitorMetricAlertDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := metricalerts.ParseMetricAlertID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.MetricAlertUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"scopes": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
				Set: pluginsdk.HashString,
			},

			"target_resource_type": {
				Type:        pluginsdk.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The resource type (e.g. Microsoft.Compute/virtualMachines) of the target pluginsdk. Required when using subscription, resource group scope or multiple scopes.`,
			},

			"target_resource_location": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
				Description:      `The location of the target pluginsdk. Required when using subscription, resource group scope or multiple scopes.`,
			},

			// static criteria
			"criteria": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MinItems:     1,
				ExactlyOneOf: []string{"criteria", "dynamic_criteria", "application_insights_web_test_location_availability_criteria"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"metric_namespace": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"metric_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"aggregation": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Average",
								"Count",
								"Minimum",
								"Maximum",
								"Total",
							}, false),
						},
						"dimension": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Include",
											"Exclude",
											"StartsWith",
										}, false),
									},
									"values": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						"operator": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(metricalerts.OperatorEquals),
								string(metricalerts.OperatorGreaterThan),
								string(metricalerts.OperatorGreaterThanOrEqual),
								string(metricalerts.OperatorLessThan),
								string(metricalerts.OperatorLessThanOrEqual),
							}, false),
						},
						"threshold": {
							Type:     pluginsdk.TypeFloat,
							Required: true,
						},
						"skip_metric_validation": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			// lintignore: S018
			"dynamic_criteria": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				// Curently, it allows to define only one dynamic criteria in one metric alert.
				MaxItems:     1,
				ExactlyOneOf: []string{"criteria", "dynamic_criteria", "application_insights_web_test_location_availability_criteria"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"metric_namespace": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"metric_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"aggregation": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Average",
								"Count",
								"Minimum",
								"Maximum",
								"Total",
							}, false),
						},
						"dimension": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Include",
											"Exclude",
											"StartsWith",
										}, false),
									},
									"values": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						"operator": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(metricalerts.DynamicThresholdOperatorLessThan),
								string(metricalerts.DynamicThresholdOperatorGreaterThan),
								string(metricalerts.DynamicThresholdOperatorGreaterOrLessThan),
							}, false),
						},
						"alert_sensitivity": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(metricalerts.DynamicThresholdSensitivityLow),
								string(metricalerts.DynamicThresholdSensitivityMedium),
								string(metricalerts.DynamicThresholdSensitivityHigh),
							}, false),
						},

						"evaluation_total_count": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Default:      4,
						},

						"evaluation_failure_count": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Default:      4,
						},

						"ignore_data_before": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
						"skip_metric_validation": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"application_insights_web_test_location_availability_criteria": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MinItems:     1,
				MaxItems:     1,
				ExactlyOneOf: []string{"criteria", "dynamic_criteria", "application_insights_web_test_location_availability_criteria"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"web_test_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.WebTestID,
						},
						"component_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ComponentID,
						},
						"failed_location_count": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},

			"action": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"action_group_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"webhook_properties": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
				Set: resourceMonitorMetricAlertActionHash,
			},

			"auto_mitigate": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"frequency": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "PT1M",
				ValidateFunc: validation.StringInSlice([]string{
					"PT1M",
					"PT5M",
					"PT15M",
					"PT30M",
					"PT1H",
				}, false),
			},

			"severity": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(0, 4),
			},

			"window_size": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "PT5M",
				ValidateFunc: validation.StringInSlice([]string{
					"PT1M",
					"PT5M",
					"PT15M",
					"PT30M",
					"PT1H",
					"PT6H",
					"PT12H",
					"P1D",
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorMetricAlertCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.MetricAlertsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := metricalerts.NewMetricAlertID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_metric_alert", id.ID())
		}
	}

	enabled := d.Get("enabled").(bool)
	autoMitigate := d.Get("auto_mitigate").(bool)
	description := d.Get("description").(string)
	scopesRaw := d.Get("scopes").(*pluginsdk.Set).List()
	severity := d.Get("severity").(int)
	frequency := d.Get("frequency").(string)
	windowSize := d.Get("window_size").(string)
	actionRaw := d.Get("action").(*pluginsdk.Set).List()
	targetResourceType := d.Get("target_resource_type").(string)
	targetResourceLocation := d.Get("target_resource_location").(string)

	t := d.Get("tags").(map[string]interface{})

	// The criteria type of "old" resource is `MetricAlertSingleResourceMultipleMetricCriteria` (rather than `MetricAlertMultipleResourceMultipleMetricCriteria`).
	// We need to keep using that type in order to keep backward compatibility. Otherwise, changing the criteria type will cause error as reported in issue:
	// https://github.com/hashicorp/terraform-provider-azurerm/issues/7910
	var isLegacy bool
	if !d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Monitor %s: %+v", id, err)
		}
		if existing.Model == nil || existing.Model.Properties.Criteria == nil {
			return fmt.Errorf("unexpected nil properties of Monitor %s", id)
		}
		_, isLegacy = existing.Model.Properties.Criteria.(metricalerts.MetricAlertSingleResourceMultipleMetricCriteria)

	}

	criteria, err := expandMonitorMetricAlertCriteria(d, isLegacy)
	if err != nil {
		return fmt.Errorf(`expanding criteria: %+v`, err)
	}

	parameters := metricalerts.MetricAlertResource{
		Location: azure.NormalizeLocation("Global"),
		Properties: metricalerts.MetricAlertProperties{
			Enabled:              enabled,
			AutoMitigate:         utils.Bool(autoMitigate),
			Description:          utils.String(description),
			Severity:             int64(severity),
			EvaluationFrequency:  frequency,
			WindowSize:           windowSize,
			Scopes:               expandStringValues(scopesRaw),
			Criteria:             criteria,
			Actions:              expandMonitorMetricAlertAction(actionRaw),
			TargetResourceType:   utils.String(targetResourceType),
			TargetResourceRegion: utils.String(targetResourceLocation),
		},
		Tags: utils.ExpandPtrMapStringString(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating or updating Monitor %s: %+v", id, err)
	}

	// Monitor Metric Alert API would return 404 while creating multiple Monitor Metric Alerts and get each resource immediately once it's created successfully in parallel.
	// Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/10973
	log.Printf("[DEBUG] Waiting for %s to be created", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   monitorMetricAlertStateRefreshFunc(ctx, client, id),
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Monitor %s to finish provisioning: %s", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorMetricAlertRead(d, meta)
}

func resourceMonitorMetricAlertRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.MetricAlertsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := metricalerts.ParseMetricAlertID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Monitor %s: %+v", *id, err)
	}

	d.Set("name", id.MetricAlertName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		props := model.Properties
		d.Set("enabled", props.Enabled)
		d.Set("auto_mitigate", props.AutoMitigate)
		d.Set("description", props.Description)
		d.Set("severity", props.Severity)
		d.Set("frequency", props.EvaluationFrequency)
		d.Set("window_size", props.WindowSize)
		if err := d.Set("scopes", props.Scopes); err != nil {
			return fmt.Errorf("setting `scopes`: %+v", err)
		}

		// Determine the correct criteria schema to set
		var criteriaSchema string
		switch c := props.Criteria.(type) {
		case metricalerts.MetricAlertSingleResourceMultipleMetricCriteria:
			criteriaSchema = "criteria"
		case metricalerts.MetricAlertMultipleResourceMultipleMetricCriteria:
			if c.AllOf == nil || len(*c.AllOf) == 0 {
				return fmt.Errorf("nil or empty contained criteria of MultipleResourceMultipleMetricCriteria")
			}
			// `MinItems` defined in schema guaranteed there is at least one element.
			switch (*c.AllOf)[0].(type) {
			case metricalerts.DynamicMetricCriteria:
				criteriaSchema = "dynamic_criteria"
			case metricalerts.MetricCriteria:
				criteriaSchema = "criteria"
			}
		case metricalerts.WebtestLocationAvailabilityCriteria:
			criteriaSchema = "application_insights_web_test_location_availability_criteria"
		default:
			return fmt.Errorf("unknown criteria type")
		}

		monitorMetricAlertCriteria := flattenMonitorMetricAlertCriteria(props.Criteria)
		// lintignore:R001
		if err := d.Set(criteriaSchema, monitorMetricAlertCriteria); err != nil {
			return fmt.Errorf("failed setting `%s`: %+v", criteriaSchema, err)
		}

		if err := d.Set("action", flattenMonitorMetricAlertAction(props.Actions)); err != nil {
			return fmt.Errorf("setting `action`: %+v", err)
		}
		d.Set("target_resource_type", props.TargetResourceType)
		d.Set("target_resource_location", props.TargetResourceRegion)

		if err = d.Set("tags", utils.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}

func resourceMonitorMetricAlertDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.MetricAlertsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := metricalerts.ParseMetricAlertID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Monitor %s: %+v", *id, err)
		}
	}

	return nil
}

func expandMonitorMetricAlertCriteria(d *pluginsdk.ResourceData, isLegacy bool) (metricalerts.MetricAlertCriteria, error) {
	switch {
	case len(d.Get("criteria").([]interface{})) != 0:
		if isLegacy {
			return expandMonitorMetricAlertSingleResourceMultiMetricCriteria(d.Get("criteria").([]interface{})), nil
		}
		return expandMonitorMetricAlertMultiResourceMultiMetricForStaticMetricCriteria(d.Get("criteria").([]interface{})), nil
	case len(d.Get("dynamic_criteria").([]interface{})) != 0:
		return expandMonitorMetricAlertMultiResourceMultiMetricForDynamicMetricCriteria(d.Get("dynamic_criteria").([]interface{})), nil
	case len(d.Get("application_insights_web_test_location_availability_criteria").([]interface{})) != 0:
		return expandMonitorMetricAlertWebtestLocAvailCriteria(d.Get("application_insights_web_test_location_availability_criteria").([]interface{})), nil
	default:
		// Guaranteed by schema `AtLeastOne` constraint
		return nil, fmt.Errorf("unknown criteria type")
	}
}

func expandMonitorMetricAlertSingleResourceMultiMetricCriteria(input []interface{}) metricalerts.MetricAlertCriteria {
	criteria := make([]metricalerts.MetricCriteria, 0)
	for i, item := range input {
		v := item.(map[string]interface{})
		dimensions := expandMonitorMetricDimension(v["dimension"].([]interface{}))
		criteria = append(criteria, metricalerts.MetricCriteria{
			Name:                 fmt.Sprintf("Metric%d", i+1),
			MetricNamespace:      utils.String(v["metric_namespace"].(string)),
			MetricName:           v["metric_name"].(string),
			TimeAggregation:      metricalerts.AggregationTypeEnum(v["aggregation"].(string)),
			Dimensions:           &dimensions,
			Operator:             metricalerts.Operator(v["operator"].(string)),
			Threshold:            v["threshold"].(float64),
			SkipMetricValidation: utils.Bool(v["skip_metric_validation"].(bool)),
		})
	}
	return &metricalerts.MetricAlertSingleResourceMultipleMetricCriteria{
		AllOf: &criteria,
	}
}

func expandMonitorMetricAlertMultiResourceMultiMetricForStaticMetricCriteria(input []interface{}) metricalerts.MetricAlertCriteria {
	criteria := make([]metricalerts.MultiMetricCriteria, 0)
	for i, item := range input {
		v := item.(map[string]interface{})
		dimensions := expandMonitorMetricDimension(v["dimension"].([]interface{}))
		criteria = append(criteria, metricalerts.MetricCriteria{
			Name:                 fmt.Sprintf("Metric%d", i+1),
			MetricNamespace:      utils.String(v["metric_namespace"].(string)),
			MetricName:           v["metric_name"].(string),
			TimeAggregation:      metricalerts.AggregationTypeEnum(v["aggregation"].(string)),
			Dimensions:           &dimensions,
			Operator:             metricalerts.Operator(v["operator"].(string)),
			Threshold:            v["threshold"].(float64),
			SkipMetricValidation: utils.Bool(v["skip_metric_validation"].(bool)),
		})
	}
	return &metricalerts.MetricAlertMultipleResourceMultipleMetricCriteria{
		AllOf: &criteria,
	}
}

func expandMonitorMetricAlertMultiResourceMultiMetricForDynamicMetricCriteria(input []interface{}) metricalerts.MetricAlertCriteria {
	criteria := make([]metricalerts.MultiMetricCriteria, 0)
	for i, item := range input {
		v := item.(map[string]interface{})
		dimensions := expandMonitorMetricDimension(v["dimension"].([]interface{}))

		dynamicMetricCriteria := metricalerts.DynamicMetricCriteria{
			Name:             fmt.Sprintf("Metric%d", i+1),
			MetricNamespace:  utils.String(v["metric_namespace"].(string)),
			MetricName:       v["metric_name"].(string),
			TimeAggregation:  metricalerts.AggregationTypeEnum(v["aggregation"].(string)),
			Dimensions:       &dimensions,
			Operator:         metricalerts.DynamicThresholdOperator(v["operator"].(string)),
			AlertSensitivity: metricalerts.DynamicThresholdSensitivity(v["alert_sensitivity"].(string)),
			FailingPeriods: metricalerts.DynamicThresholdFailingPeriods{
				NumberOfEvaluationPeriods: float64(v["evaluation_total_count"].(int)),
				MinFailingPeriodsToAlert:  float64(v["evaluation_failure_count"].(int)),
			},
			SkipMetricValidation: utils.Bool(v["skip_metric_validation"].(bool)),
		}

		if datetime := v["ignore_data_before"].(string); datetime != "" {
			// Guaranteed in schema validation func.
			t, _ := time.Parse(time.RFC3339, datetime)
			ignoreDataBefore := &date.Time{Time: t}
			dynamicMetricCriteria.IgnoreDataBefore = pointer.To(ignoreDataBefore.String())
		}

		criteria = append(criteria, dynamicMetricCriteria)
	}
	return &metricalerts.MetricAlertMultipleResourceMultipleMetricCriteria{
		AllOf: &criteria,
	}
}

func expandMonitorMetricAlertWebtestLocAvailCriteria(input []interface{}) metricalerts.MetricAlertCriteria {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &metricalerts.WebtestLocationAvailabilityCriteria{
		WebTestId:           v["web_test_id"].(string),
		ComponentId:         v["component_id"].(string),
		FailedLocationCount: float64(v["failed_location_count"].(int)),
	}
}

func expandMonitorMetricDimension(input []interface{}) []metricalerts.MetricDimension {
	result := make([]metricalerts.MetricDimension, 0)
	for _, dimension := range input {
		dVal := dimension.(map[string]interface{})
		result = append(result, metricalerts.MetricDimension{
			Name:     dVal["name"].(string),
			Operator: dVal["operator"].(string),
			Values:   expandStringValues(dVal["values"].([]interface{})),
		})
	}
	return result
}

func expandMonitorMetricAlertAction(input []interface{}) *[]metricalerts.MetricAlertAction {
	actions := make([]metricalerts.MetricAlertAction, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		if agID := v["action_group_id"].(string); agID != "" {
			props := make(map[string]string)
			if pVal, ok := v["webhook_properties"]; ok {
				for pk, pv := range pVal.(map[string]interface{}) {
					props[pk] = pv.(string)
				}
			}

			actions = append(actions, metricalerts.MetricAlertAction{
				ActionGroupId:     utils.String(agID),
				WebHookProperties: &props,
			})
		}
	}
	return &actions
}

func flattenMonitorMetricAlertCriteria(input metricalerts.MetricAlertCriteria) []interface{} {
	switch criteria := input.(type) {
	case metricalerts.MetricAlertSingleResourceMultipleMetricCriteria:
		return flattenMonitorMetricAlertSingleResourceMultiMetricCriteria(criteria.AllOf)
	case metricalerts.MetricAlertMultipleResourceMultipleMetricCriteria:
		return flattenMonitorMetricAlertMultiResourceMultiMetricCriteria(criteria.AllOf)
	case metricalerts.WebtestLocationAvailabilityCriteria:
		return flattenMonitorMetricAlertWebtestLocAvailCriteria(&criteria)
	default:
		return nil
	}
}

func flattenMonitorMetricAlertSingleResourceMultiMetricCriteria(input *[]metricalerts.MetricCriteria) []interface{} {
	if input == nil || len(*input) == 0 {
		return nil
	}
	criteria := (*input)[0]
	metricName := criteria.MetricName
	metricNamespace := criteria.MetricNamespace
	timeAggregation := criteria.TimeAggregation

	dimResult := make([]map[string]interface{}, 0)
	if criteria.Dimensions != nil {
		for _, dimension := range *criteria.Dimensions {
			dVal := make(map[string]interface{})
			dVal["name"] = dimension.Name
			dVal["operator"] = dimension.Operator
			dVal["values"] = dimension.Values
			dimResult = append(dimResult, dVal)
		}
	}

	operator := string(criteria.Operator)
	threshold := criteria.Threshold

	var skipMetricValidation bool
	if criteria.SkipMetricValidation != nil {
		skipMetricValidation = *criteria.SkipMetricValidation
	}

	return []interface{}{
		map[string]interface{}{
			"metric_namespace":       metricNamespace,
			"metric_name":            metricName,
			"aggregation":            timeAggregation,
			"dimension":              dimResult,
			"operator":               operator,
			"threshold":              threshold,
			"skip_metric_validation": skipMetricValidation,
		},
	}
}

func flattenMonitorMetricAlertMultiResourceMultiMetricCriteria(input *[]metricalerts.MultiMetricCriteria) []interface{} {
	if input == nil {
		return nil
	}
	result := make([]interface{}, 0)

	for _, criteria := range *input {
		v := make(map[string]interface{})
		var (
			metricName           string
			metricNamespace      string
			timeAggregation      interface{}
			dimensions           []metricalerts.MetricDimension
			skipMetricValidation bool
		)

		switch criteria := criteria.(type) {
		case metricalerts.MetricCriteria:
			metricName = criteria.MetricName

			if criteria.MetricNamespace != nil {
				metricNamespace = *criteria.MetricNamespace
			}
			timeAggregation = criteria.TimeAggregation
			if criteria.Dimensions != nil {
				dimensions = *criteria.Dimensions
			}

			// MetricCriteria specific properties
			v["operator"] = string(criteria.Operator)
			v["threshold"] = criteria.Threshold
			if criteria.SkipMetricValidation != nil {
				skipMetricValidation = *criteria.SkipMetricValidation
			}
		case metricalerts.DynamicMetricCriteria:
			metricName = criteria.MetricName

			if criteria.MetricNamespace != nil {
				metricNamespace = *criteria.MetricNamespace
			}
			timeAggregation = criteria.TimeAggregation
			if criteria.Dimensions != nil {
				dimensions = *criteria.Dimensions
			}
			if criteria.SkipMetricValidation != nil {
				skipMetricValidation = *criteria.SkipMetricValidation
			}
			// DynamicMetricCriteria specific properties
			v["operator"] = string(criteria.Operator)
			v["alert_sensitivity"] = string(criteria.AlertSensitivity)

			v["evaluation_total_count"] = int(criteria.FailingPeriods.NumberOfEvaluationPeriods)
			v["evaluation_failure_count"] = int(criteria.FailingPeriods.MinFailingPeriodsToAlert)

			ignoreDataBefore := ""
			if criteria.IgnoreDataBefore != nil {
				ignoreDataBefore = *criteria.IgnoreDataBefore
			}
			v["ignore_data_before"] = ignoreDataBefore
		}

		// Common properties
		v["metric_name"] = metricName
		v["metric_namespace"] = metricNamespace
		v["aggregation"] = timeAggregation
		v["skip_metric_validation"] = skipMetricValidation
		if dimensions != nil {
			dimResult := make([]map[string]interface{}, 0)
			for _, dimension := range dimensions {
				dVal := make(map[string]interface{})
				dVal["name"] = dimension.Name
				dVal["operator"] = dimension.Operator
				dVal["values"] = dimension.Values
				dimResult = append(dimResult, dVal)
			}
			v["dimension"] = dimResult
		}

		result = append(result, v)
	}
	return result
}

func flattenMonitorMetricAlertWebtestLocAvailCriteria(input *metricalerts.WebtestLocationAvailabilityCriteria) []interface{} {
	if input == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"web_test_id":           input.WebTestId,
			"component_id":          input.ComponentId,
			"failed_location_count": int(input.FailedLocationCount),
		},
	}
}

func flattenMonitorMetricAlertAction(input *[]metricalerts.MetricAlertAction) (result []interface{}) {
	result = make([]interface{}, 0)
	if input == nil {
		return
	}
	for _, action := range *input {
		v := make(map[string]interface{})

		if action.ActionGroupId != nil {
			v["action_group_id"] = *action.ActionGroupId
		}

		props := make(map[string]string)
		if action.WebHookProperties != nil {
			for pk, pv := range *action.WebHookProperties {
				props[pk] = pv
			}
		}
		v["webhook_properties"] = props

		result = append(result, v)
	}

	return result
}

func resourceMonitorMetricAlertActionHash(input interface{}) int {
	var buf bytes.Buffer
	if v, ok := input.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", v["action_group_id"].(string)))
		if m, ok := v["webhook_properties"].(map[string]interface{}); ok && m != nil {
			buf.WriteString(fmt.Sprintf("%v-", m))
		}
	}
	return pluginsdk.HashString(buf.String())
}

func monitorMetricAlertStateRefreshFunc(ctx context.Context, client *metricalerts.MetricAlertsClient, id metricalerts.MetricAlertId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return nil, "404", nil
			}

			return nil, "", fmt.Errorf("retrieving %s: %s", id, err)
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
}
