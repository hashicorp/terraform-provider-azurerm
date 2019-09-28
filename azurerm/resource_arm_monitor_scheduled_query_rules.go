package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorScheduledQueryRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorScheduledQueryRulesCreateUpdate,
		Read:   resourceArmMonitorScheduledQueryRulesRead,
		Update: resourceArmMonitorScheduledQueryRulesCreateUpdate,
		Delete: resourceArmMonitorScheduledQueryRulesDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"source": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:         schema.TypeString,
							Required:     true,
						},
						"dataSourceId": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"authorizedResources": {
							Type:			schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
						},
						"queryType": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "ResultCount",
							ValidateFunc: validation.StringInSlice([]string{
								"ResultCount",
							}, false),
						},
					},
				},
				Set: schema.HashString,
			},

			"schedule": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequencyInMinutes": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"timeWindowInMinutes": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"action": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"odata.type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{
								"Microsoft.WindowsAzure.Management.Monitoring.Alerts.Models.Microsoft.AppInsights.Nexus.DataContracts.Resources.ScheduledQueryRules.AlertingAction",
								"Microsoft.WindowsAzure.Management.Monitoring.Alerts.Models.Microsoft.AppInsights.Nexus.DataContracts.Resources.ScheduledQueryRules.LogToMetricAction",
							}, false),
						},
						"severity": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{
								"0",
								"1",
								"2",
								"3",
								"4",
							}, false),
						},
						"throttlingInMin": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"aznsAction": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"actionGroup": {
										Type:         schema.TypeSet,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},
									"customWebhookPayload": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validate.URLIsHTTPOrHTTPS,
									},
									"emailSubject": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"criteria": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dimensions": {
										Type:         schema.TypeSet,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},
									"metricName": {
										Type:         schema.TypeString,
										Required:     true,
									},
								},
							},
						},
						"trigger": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"thresholdOperator": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{
											"GreaterThan",
											"LessThan",
											"Equal",
										}, false),
									},
									"threshold": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
				Set: resourceArmMonitorScheduledQueryRulesActionHash,
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

			"lastUpdatedTime": {
				Type:			schema.TypeString,
				Computed: true,
			},

			"provisioningState": {
				Type:			schema.TypeString,
				Computed: true,
			}

			"displayName": {
				Type:			schema.TypeString,
				Computed: true,
			}

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorScheduledQueryRulesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Monitor Scheduled Query Rules %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_scheduled_query_rules", *existing.ID)
		}
	}

	enabled := d.Get("enabled").(bool)
	description := d.Get("description").(string)
	displayName := d.Get("displayName").(string)
	sourceRaw := d.Get("source").(*schema.Set).List()
	scheduleRaw := d.Get("schedule").(*schema.Set).List()
	actionRaw := d.Get("action").(*schema.Set).List()

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := insights.LogSearchRuleResource{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		LogSearchRule: &insights.LogSearchRule{
			Enabled:     utils.Bool(enabled),
			Description: utils.String(description),
			DisplayName: utils.String(displayName),
			Source:      expandMonitorScheduledQueryRulesSource(sourceRaw),
			Schedule:    expandMonitorScheduledQueryRulesSchedule(scheduleRaw),
			Action:      expandMonitorScheduledQueryRulesAction(actionRaw),
		},
		Tags: expandedTags,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating or updating scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Scheduled query rule %q (resource group %q) ID is empty", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceArmMonitorScheduledQueryRulesRead(d, meta)
}

func resourceArmMonitorScheduledQueryRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ScheduledQueryRules"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Scheduled Query Rule %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if rule := resp.ScheduledQueryRules; rule != nil {
		d.Set("enabled", rule.Enabled)
		d.Set("description", rule.Description)
		d.Set("displayName", rule.DisplayName)
		if err := d.Set("source", flattenMonitorScheduledQueryRulesSource(rule.Source); err != nil {
			return fmt.Errorf("Error setting `source`: %+v", err)
		}
		if err := d.Set("schedule", flattenMonitorScheduledQueryRulesSchedule(rule.Schedule)); err != nil {
			return fmt.Errorf("Error setting `schedule`: %+v", err)
		}
		if err := d.Set("action", flattenMonitorScheduledQueryRulesAction(rule.Action)); err != nil {
			return fmt.Errorf("Error setting `action`: %+v", err)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMonitorScheduledQueryRulesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ScheduledQueryRules"]

	if resp, err := client.Delete(ctx, resourceGroup, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandMonitorScheduledQueryRulesSource(input []interface{}) (*insights.Source,  error) {
	conditions := make([]insights.ScheduledQueryRulesLeafCondition, 0)
	v := input[0].(map[string]interface{})

	if category := v["category"].(string); category != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("category"),
			Equals: utils.String(category),
		})
	}
	if op := v["operation_name"].(string); op != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("operationName"),
			Equals: utils.String(op),
		})
	}
	if caller := v["caller"].(string); caller != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("caller"),
			Equals: utils.String(caller),
		})
	}
	if level := v["level"].(string); level != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("level"),
			Equals: utils.String(level),
		})
	}
	if resourceProvider := v["resource_provider"].(string); resourceProvider != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("resourceProvider"),
			Equals: utils.String(resourceProvider),
		})
	}
	if resourceType := v["resource_type"].(string); resourceType != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("resourceType"),
			Equals: utils.String(resourceType),
		})
	}
	if resourceGroup := v["resource_group"].(string); resourceGroup != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("resourceGroup"),
			Equals: utils.String(resourceGroup),
		})
	}
	if id := v["resource_id"].(string); id != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("resourceId"),
			Equals: utils.String(id),
		})
	}
	if status := v["status"].(string); status != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("status"),
			Equals: utils.String(status),
		})
	}
	if subStatus := v["sub_status"].(string); subStatus != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("subStatus"),
			Equals: utils.String(subStatus),
		})
	}

	return &insights.ScheduledQueryRulesAllOfCondition{
		AllOf: &conditions,
	}
}

func expandMonitorScheduledQueryRulesSchedule(input []interface{}) (*insights.Schedule, error) {
	actions := make([]insights.ScheduledQueryRulesActionGroup, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		if agID := v["action_group_id"].(string); agID != "" {
			props := make(map[string]*string)
			if pVal, ok := v["webhook_properties"]; ok {
				for pk, pv := range pVal.(map[string]interface{}) {
					props[pk] = utils.String(pv.(string))
				}
			}

			actions = append(actions, insights.ScheduledQueryRulesActionGroup{
				ActionGroupID:     utils.String(agID),
				WebhookProperties: props,
			})
		}
	}
	return &insights.ScheduledQueryRulesActionList{
		ActionGroups: &actions,
	}
}

func expandMonitorScheduledQueryRulesAction(input []interface{}) (*insights.BasicAction, error) {
	actions := make([]insights.ScheduledQueryRulesActionGroup, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		if agID := v["action_group_id"].(string); agID != "" {
			props := make(map[string]*string)
			if pVal, ok := v["webhook_properties"]; ok {
				for pk, pv := range pVal.(map[string]interface{}) {
					props[pk] = utils.String(pv.(string))
				}
			}

			actions = append(actions, insights.ScheduledQueryRulesActionGroup{
				ActionGroupID:     utils.String(agID),
				WebhookProperties: props,
			})
		}
	}
	return &insights.ScheduledQueryRulesActionList{
		ActionGroups: &actions,
	}
}

func flattenMonitorScheduledQueryRulesCriteria(input *insights.ScheduledQueryRulesAllOfCondition) []interface{} {
	result := make(map[string]interface{})
	if input == nil || input.AllOf == nil {
		return []interface{}{result}
	}
	for _, condition := range *input.AllOf {
		if condition.Field != nil && condition.Equals != nil {
			switch strings.ToLower(*condition.Field) {
			case "operationname":
				result["operation_name"] = *condition.Equals
			case "resourceprovider":
				result["resource_provider"] = *condition.Equals
			case "resourcetype":
				result["resource_type"] = *condition.Equals
			case "resourcegroup":
				result["resource_group"] = *condition.Equals
			case "resourceid":
				result["resource_id"] = *condition.Equals
			case "substatus":
				result["sub_status"] = *condition.Equals
			case "caller", "category", "level", "status":
				result[*condition.Field] = *condition.Equals
			}
		}
	}
	return []interface{}{result}
}

func flattenMonitorScheduledQueryRulesAction(input *insights.ScheduledQueryRulesActionList) (result []interface{}) {
	result = make([]interface{}, 0)
	if input == nil || input.ActionGroups == nil {
		return
	}
	for _, action := range *input.ActionGroups {
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

func resourceArmMonitorScheduledQueryRulesActionHash(input interface{}) int {
	var buf bytes.Buffer
	if v, ok := input.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", v["action_group_id"].(string)))
	}
	return hashcode.String(buf.String())
}
