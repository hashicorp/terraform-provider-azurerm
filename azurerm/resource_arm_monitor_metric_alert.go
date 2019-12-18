package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorMetricAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorMetricAlertCreateUpdate,
		Read:   resourceArmMonitorMetricAlertRead,
		Update: resourceArmMonitorMetricAlertCreateUpdate,
		Delete: resourceArmMonitorMetricAlertDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			// TODO: Multiple resource IDs (Remove MaxItems) support is missing in SDK
			// Issue to track: https://github.com/Azure/azure-sdk-for-go/issues/2920
			// But to prevent potential state migration in the future, let's stick to use Set now

			//lintignore:S018
			"scopes": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
				Set: schema.HashString,
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
							ValidateFunc: validate.NoEmptyStrings,
						},
						"metric_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"aggregation": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Average",
								"Count",
								"Minimum",
								"Maximum",
								"Total",
							}, false),
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
							}, false),
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
										ValidateFunc: validate.NoEmptyStrings,
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Include",
											"Exclude",
										}, false),
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"action": {
				Type:     schema.TypeSet,
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
				Set: resourceArmMonitorMetricAlertActionHash,
			},

			"auto_mitigate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorMetricAlertCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Monitor.MetricAlertsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Monitor Metric Alert %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_metric_alert", *existing.ID)
		}
	}

	enabled := d.Get("enabled").(bool)
	autoMitigate := d.Get("auto_mitigate").(bool)
	description := d.Get("description").(string)
	scopesRaw := d.Get("scopes").(*schema.Set).List()
	severity := d.Get("severity").(int)
	frequency := d.Get("frequency").(string)
	windowSize := d.Get("window_size").(string)
	criteriaRaw := d.Get("criteria").([]interface{})
	actionRaw := d.Get("action").(*schema.Set).List()

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := insights.MetricAlertResource{
		Location: utils.String(azure.NormalizeLocation("Global")),
		MetricAlertProperties: &insights.MetricAlertProperties{
			Enabled:             utils.Bool(enabled),
			AutoMitigate:        utils.Bool(autoMitigate),
			Description:         utils.String(description),
			Severity:            utils.Int32(int32(severity)),
			EvaluationFrequency: utils.String(frequency),
			WindowSize:          utils.String(windowSize),
			Scopes:              utils.ExpandStringSlice(scopesRaw),
			Criteria:            expandMonitorMetricAlertCriteria(criteriaRaw),
			Actions:             expandMonitorMetricAlertAction(actionRaw),
		},
		Tags: expandedTags,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating or updating metric alert %q (resource group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Metric alert %q (resource group %q) ID is empty", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceArmMonitorMetricAlertRead(d, meta)
}

func resourceArmMonitorMetricAlertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Monitor.MetricAlertsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["metricAlerts"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Metric Alert %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting metric alert %q (resource group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if alert := resp.MetricAlertProperties; alert != nil {
		d.Set("enabled", alert.Enabled)
		d.Set("auto_mitigate", alert.AutoMitigate)
		d.Set("description", alert.Description)
		d.Set("severity", alert.Severity)
		d.Set("frequency", alert.EvaluationFrequency)
		d.Set("window_size", alert.WindowSize)
		if err := d.Set("scopes", utils.FlattenStringSlice(alert.Scopes)); err != nil {
			return fmt.Errorf("Error setting `scopes`: %+v", err)
		}
		if err := d.Set("criteria", flattenMonitorMetricAlertCriteria(alert.Criteria)); err != nil {
			return fmt.Errorf("Error setting `criteria`: %+v", err)
		}
		if err := d.Set("action", flattenMonitorMetricAlertAction(alert.Actions)); err != nil {
			return fmt.Errorf("Error setting `action`: %+v", err)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMonitorMetricAlertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Monitor.MetricAlertsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["metricAlerts"]

	if resp, err := client.Delete(ctx, resourceGroup, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting metric alert %q (resource group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandMonitorMetricAlertCriteria(input []interface{}) *insights.MetricAlertSingleResourceMultipleMetricCriteria {
	criteria := make([]insights.MetricCriteria, 0)
	for i, item := range input {
		v := item.(map[string]interface{})

		dimensions := make([]insights.MetricDimension, 0)
		for _, dimension := range v["dimension"].([]interface{}) {
			dVal := dimension.(map[string]interface{})
			dimensions = append(dimensions, insights.MetricDimension{
				Name:     utils.String(dVal["name"].(string)),
				Operator: utils.String(dVal["operator"].(string)),
				Values:   utils.ExpandStringSlice(dVal["values"].([]interface{})),
			})
		}

		criteria = append(criteria, insights.MetricCriteria{
			Name:            utils.String(fmt.Sprintf("Metric%d", i+1)),
			MetricNamespace: utils.String(v["metric_namespace"].(string)),
			MetricName:      utils.String(v["metric_name"].(string)),
			TimeAggregation: v["aggregation"].(string),
			Operator:        v["operator"].(string),
			Threshold:       utils.Float(v["threshold"].(float64)),
			Dimensions:      &dimensions,
		})
	}
	return &insights.MetricAlertSingleResourceMultipleMetricCriteria{
		AllOf:     &criteria,
		OdataType: insights.OdataTypeMicrosoftAzureMonitorSingleResourceMultipleMetricCriteria,
	}
}

func expandMonitorMetricAlertAction(input []interface{}) *[]insights.MetricAlertAction {
	actions := make([]insights.MetricAlertAction, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		if agID := v["action_group_id"].(string); agID != "" {
			props := make(map[string]*string)
			if pVal, ok := v["webhook_properties"]; ok {
				for pk, pv := range pVal.(map[string]interface{}) {
					props[pk] = utils.String(pv.(string))
				}
			}

			actions = append(actions, insights.MetricAlertAction{
				ActionGroupID:     utils.String(agID),
				WebhookProperties: props,
			})
		}
	}
	return &actions
}

func flattenMonitorMetricAlertCriteria(input insights.BasicMetricAlertCriteria) (result []interface{}) {
	result = make([]interface{}, 0)
	if input == nil {
		return
	}
	criteria, ok := input.AsMetricAlertSingleResourceMultipleMetricCriteria()
	if !ok || criteria == nil || criteria.AllOf == nil {
		return
	}
	for _, metric := range *criteria.AllOf {
		v := make(map[string]interface{})

		if metric.MetricNamespace != nil {
			v["metric_namespace"] = *metric.MetricNamespace
		}
		if metric.MetricName != nil {
			v["metric_name"] = *metric.MetricName
		}
		if aggr, ok := metric.TimeAggregation.(string); ok {
			v["aggregation"] = aggr
		}
		if op, ok := metric.Operator.(string); ok {
			v["operator"] = op
		}
		if metric.Threshold != nil {
			v["threshold"] = *metric.Threshold
		}
		if metric.Dimensions != nil {
			dimResult := make([]map[string]interface{}, 0)
			for _, dimension := range *metric.Dimensions {
				dVal := make(map[string]interface{})
				if dimension.Name != nil {
					dVal["name"] = *dimension.Name
				}
				if dimension.Operator != nil {
					dVal["operator"] = *dimension.Operator
				}
				dVal["values"] = utils.FlattenStringSlice(dimension.Values)
				dimResult = append(dimResult, dVal)
			}
			v["dimension"] = dimResult
		}

		result = append(result, v)
	}

	return result
}

func flattenMonitorMetricAlertAction(input *[]insights.MetricAlertAction) (result []interface{}) {
	result = make([]interface{}, 0)
	if input == nil {
		return
	}
	for _, action := range *input {
		v := make(map[string]interface{})

		if action.ActionGroupID != nil {
			v["action_group_id"] = *action.ActionGroupID
		}

		props := make(map[string]string)
		for pk, pv := range action.WebhookProperties {
			if pv != nil {
				props[pk] = *pv
			}
		}
		v["webhook_properties"] = props

		result = append(result, v)
	}

	return result
}

func resourceArmMonitorMetricAlertActionHash(input interface{}) int {
	var buf bytes.Buffer
	if v, ok := input.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", v["action_group_id"].(string)))
	}
	return hashcode.String(buf.String())
}
