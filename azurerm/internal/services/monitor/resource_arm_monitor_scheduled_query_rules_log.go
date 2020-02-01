package monitor

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorScheduledQueryRulesLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorScheduledQueryRulesLogCreateUpdate,
		Read:   resourceArmMonitorScheduledQueryRulesRead,
		Update: resourceArmMonitorScheduledQueryRulesLogCreateUpdate,
		Delete: resourceArmMonitorScheduledQueryRulesDelete,
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

			"location": azure.SchemaLocation(),

			"authorized_resource_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"criteria": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dimension": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "Include",
										ValidateFunc: validation.StringInSlice([]string{
											"Include",
										}, false),
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validate.NoEmptyStrings,
										},
									},
								},
							},
						},
						"metric_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},
			"data_source_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
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
			"last_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorScheduledQueryRulesLogCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	action := expandMonitorScheduledQueryRulesLogToMetricAction(d)

	return resourceArmMonitorScheduledQueryRulesCreateUpdate(d, meta, action, nil)
}

func expandMonitorScheduledQueryRulesLogCriteria(input []interface{}) *[]insights.Criteria {
	criteria := make([]insights.Criteria, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		dimensions := make([]insights.Dimension, 0)
		for _, dimension := range v["dimension"].(*schema.Set).List() {
			dVal := dimension.(map[string]interface{})
			dimensions = append(dimensions, insights.Dimension{
				Name:     utils.String(dVal["name"].(string)),
				Operator: utils.String(dVal["operator"].(string)),
				Values:   utils.ExpandStringSlice(dVal["values"].([]interface{})),
			})
		}

		criteria = append(criteria, insights.Criteria{
			MetricName: utils.String(v["metric_name"].(string)),
			Dimensions: &dimensions,
		})
	}
	return &criteria
}

func expandMonitorScheduledQueryRulesLogSource(d *schema.ResourceData) *insights.Source {
	authorizedResourceIDs := d.Get("authorized_resource_ids").(*schema.Set).List()
	dataSourceID := d.Get("data_source_id").(string)

	source := insights.Source{
		AuthorizedResources: utils.ExpandStringSlice(authorizedResourceIDs),
		DataSourceID:        utils.String(dataSourceID),
		QueryType:           insights.ResultCount,
	}

	return &source
}

func expandMonitorScheduledQueryRulesLogToMetricAction(d *schema.ResourceData) *insights.LogToMetricAction {
	criteriaRaw := d.Get("criteria").(*schema.Set).List()
	criteria := expandMonitorScheduledQueryRulesLogCriteria(criteriaRaw)

	action := insights.LogToMetricAction{
		Criteria:  criteria,
		OdataType: insights.OdataTypeMicrosoftWindowsAzureManagementMonitoringAlertsModelsMicrosoftAppInsightsNexusDataContractsResourcesScheduledQueryRulesLogToMetricAction,
	}

	return &action
}

func flattenAzureRmScheduledQueryRulesLogCriteria(input *[]insights.Criteria) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, criteria := range *input {
			v := make(map[string]interface{})

			v["dimension"] = flattenAzureRmScheduledQueryRulesLogDimension(criteria.Dimensions)
			v["metric_name"] = *criteria.MetricName

			result = append(result, v)
		}
	}

	return result
}

func flattenAzureRmScheduledQueryRulesLogDimension(input *[]insights.Dimension) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, dimension := range *input {
			v := make(map[string]interface{})

			if dimension.Name != nil {
				v["name"] = *dimension.Name
			}

			if dimension.Operator != nil {
				v["operator"] = *dimension.Operator
			}

			if dimension.Values != nil {
				v["values"] = *dimension.Values
			}

			result = append(result, v)
		}
	}

	return result
}
