package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorMetricAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorMetricAlertCreateOrUpdate,
		Read:   resourceArmMonitorMetricAlertRead,
		Update: resourceArmMonitorMetricAlertCreateOrUpdate,
		Delete: resourceArmMonitorMetricAlertDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"target_ids": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"criteria": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_namespace": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"metric_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"aggregation": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Average",
								"Minimum",
								"Maximum",
								"Total",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
						"operator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Equals",
								"NotEquals",
								"GreaterThan",
								"GreaterThanOrEqual",
								"LessThan",
								"LessThanOrEqual",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"dimension": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.NoZeroValues,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"action": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_group_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"webhook_properties": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"frequency": {
				Type:     schema.TypeString,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(0, 4),
			},

			"window_size": {
				Type:     schema.TypeString,
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmMonitorMetricAlertCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorMetricAlertsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	enabled := d.Get("enabled").(bool)
	description := d.Get("description").(string)
	targetsRaw := d.Get("target_ids").([]interface{})
	severity := d.Get("severity").(int)
	frequency := d.Get("frequency").(string)
	windowSize := d.Get("window_size").(string)
	criteriaRaw := d.Get("criteria").([]interface{})
	actionRaw := d.Get("action").([]interface{})

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := insights.MetricAlertResource{
		Location: utils.String(azureRMNormalizeLocation("Global")),
		MetricAlertProperties: &insights.MetricAlertProperties{
			Enabled:             utils.Bool(enabled),
			Description:         utils.String(description),
			Severity:            utils.Int32(int32(severity)),
			EvaluationFrequency: utils.String(frequency),
			WindowSize:          utils.String(windowSize),
			Scopes:              expandMonitorMetricAlertStringArray(targetsRaw),
			Criteria:            expandMonitorMetricAlertCriteria(criteriaRaw),
			Actions:             expandMonitorMetricAlertAction(actionRaw),
		},
		Tags: expandedTags,
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating or updating metric alert %q (resource group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Metric alert %q (resource group %q) ID is empty", name, resGroup)
	}
	d.SetId(*read.ID)

	return resourceArmMonitorMetricAlertRead(d, meta)
}

func resourceArmMonitorMetricAlertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorMetricAlertsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["metricAlerts"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting metric alert %q (resource group %q): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if alert := resp.MetricAlertProperties; alert != nil {
		d.Set("enabled", alert.Enabled)
		d.Set("description", alert.Description)
		d.Set("severity", alert.Severity)
		d.Set("frequency", alert.EvaluationFrequency)
		d.Set("window_size", alert.WindowSize)
		if err := d.Set("target_ids", flattenMonitorMetricAlertStringArray(alert.Scopes)); err != nil {
			return err
		}
		if err := d.Set("criteria", flattenMonitorMetricAlertCriteria(alert.Criteria)); err != nil {
			return err
		}
		if err := d.Set("action", flattenMonitorMetricAlertAction(alert.Actions)); err != nil {
			return err
		}
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmMonitorMetricAlertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorMetricAlertsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["metricAlerts"]

	if resp, err := client.Delete(ctx, resGroup, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting metric alert %q (resource group %q): %+v", name, resGroup, err)
		}
	}

	return nil
}

func expandMonitorMetricAlertStringArray(v []interface{}) *[]string {
	result := make([]string, 0)
	for _, item := range v {
		result = append(result, item.(string))
	}
	return &result
}

func expandMonitorMetricAlertCriteria(v []interface{}) *insights.MetricAlertSingleResourceMultipleMetricCriteria {
	criterias := make([]insights.MetricCriteria, 0)
	for i, criteriaValue := range v {
		val := criteriaValue.(map[string]interface{})

		mNS := val["metric_namespace"].(string)
		mName := val["metric_name"].(string)
		aggr := val["aggregation"].(string)
		op := val["operator"].(string)
		threshold := val["threshold"].(float64)

		dimensions := make([]insights.MetricDimension, 0)
		for _, dimension := range val["dimension"].([]interface{}) {
			dVal := dimension.(map[string]interface{})
			dName := dVal["name"].(string)
			dOp := dVal["operator"].(string)
			dValuesRaw := dVal["values"].([]interface{})
			dimensions = append(dimensions, insights.MetricDimension{
				Name:     utils.String(dName),
				Operator: utils.String(dOp),
				Values:   expandMonitorMetricAlertStringArray(dValuesRaw),
			})
		}

		criterias = append(criterias, insights.MetricCriteria{
			Name:            utils.String(fmt.Sprintf("Metric%d", i+1)),
			MetricNamespace: utils.String(mNS),
			MetricName:      utils.String(mName),
			TimeAggregation: aggr,
			Operator:        op,
			Threshold:       utils.Float(threshold),
			Dimensions:      &dimensions,
		})
	}
	return &insights.MetricAlertSingleResourceMultipleMetricCriteria{
		AllOf:     &criterias,
		OdataType: insights.OdataTypeMicrosoftAzureMonitorSingleResourceMultipleMetricCriteria,
	}
}

func expandMonitorMetricAlertAction(v []interface{}) *[]insights.MetricAlertAction {
	actions := make([]insights.MetricAlertAction, 0)
	for _, actionValue := range v {
		val := actionValue.(map[string]interface{})

		agID := val["action_group_id"].(string)
		props := make(map[string]*string)
		if propsValue, ok := val["webhook_properties"]; ok {
			for k, v := range propsValue.(map[string]interface{}) {
				props[k] = utils.String(v.(string))
			}
		}

		actions = append(actions, insights.MetricAlertAction{
			ActionGroupID:     utils.String(agID),
			WebhookProperties: props,
		})
	}
	return &actions
}

func flattenMonitorMetricAlertStringArray(v *[]string) []interface{} {
	result := make([]interface{}, 0)
	if v != nil {
		for _, item := range *v {
			result = append(result, item)
		}
	}
	return result
}

func flattenMonitorMetricAlertCriteria(v insights.BasicMetricAlertCriteria) []interface{} {
	result := make([]interface{}, 0)
	if v != nil {
		if val, ok := v.AsMetricAlertSingleResourceMultipleMetricCriteria(); ok && val != nil && val.AllOf != nil {
			for _, metric := range *val.AllOf {
				item := make(map[string]interface{})

				if metric.MetricNamespace != nil {
					item["metric_namespace"] = *metric.MetricNamespace
				}
				if metric.MetricName != nil {
					item["metric_name"] = *metric.MetricName
				}
				if aggr, ok := metric.TimeAggregation.(string); ok {
					item["aggregation"] = aggr
				}
				if op, ok := metric.Operator.(string); ok {
					item["operator"] = op
				}
				if metric.Threshold != nil {
					item["threshold"] = *metric.Threshold
				}
				if metric.Dimensions != nil {
					dimResult := make([]map[string]interface{}, 0)
					for _, dimension := range *metric.Dimensions {
						dimItem := make(map[string]interface{})
						if dimension.Name != nil {
							dimItem["name"] = *dimension.Name
						}
						if dimension.Operator != nil {
							dimItem["operator"] = *dimension.Operator
						}
						dimItem["values"] = flattenMonitorMetricAlertStringArray(dimension.Values)
						dimResult = append(dimResult, dimItem)
					}
					item["dimension"] = dimResult
				}

				result = append(result, item)
			}
		}
	}
	return result
}

func flattenMonitorMetricAlertAction(v *[]insights.MetricAlertAction) []interface{} {
	result := make([]interface{}, 0)
	if v != nil {
		for _, action := range *v {
			item := make(map[string]interface{}, 0)

			if action.ActionGroupID != nil {
				item["action_group_id"] = *action.ActionGroupID
			}

			props := make(map[string]string)
			for k, v := range action.WebhookProperties {
				if v != nil {
					props[k] = *v
				}
			}
			item["webhook_properties"] = props

			result = append(result, item)
		}
	}
	return result
}
