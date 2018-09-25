package azurerm

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorActivityLogAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorActivityLogAlertCreateOrUpdate,
		Read:   resourceArmMonitorActivityLogAlertRead,
		Update: resourceArmMonitorActivityLogAlertCreateOrUpdate,
		Delete: resourceArmMonitorActivityLogAlertDelete,

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

			"scopes": {
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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"caller": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"level": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"operation_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"status": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sub_status": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"action": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmMonitorActivityLogAlertCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorActivityLogAlertsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	enabled := d.Get("enabled").(bool)
	description := d.Get("description").(string)
	scopesRaw := d.Get("scopes").([]interface{})
	criteriaRaw := d.Get("criteria.0").(map[string]interface{})
	actionRaw := d.Get("action").([]interface{})

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := insights.ActivityLogAlertResource{
		Location: utils.String(azureRMNormalizeLocation("Global")),
		ActivityLogAlert: &insights.ActivityLogAlert{
			Enabled:     utils.Bool(enabled),
			Description: utils.String(description),
			Scopes:      expandMonitorActivityLogAlertScopes(scopesRaw),
			Condition:   expandMonitorActivityLogAlertCriteria(criteriaRaw),
			Actions:     expandMonitorActivityLogAlertAction(actionRaw),
		},
		Tags: expandedTags,
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating or updating activity log alert %q (resource group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Activity log alert %q (resource group %q) ID is empty", name, resGroup)
	}
	d.SetId(*read.ID)

	return resourceArmMonitorActivityLogAlertRead(d, meta)
}

func resourceArmMonitorActivityLogAlertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorActivityLogAlertsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["activityLogAlerts"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting activity log alert %q (resource group %q): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if alert := resp.ActivityLogAlert; alert != nil {
		d.Set("enabled", alert.Enabled)
		d.Set("description", alert.Description)
		if err := d.Set("scopes", flattenMonitorActivityLogAlertScopes(alert.Scopes)); err != nil {
			return err
		}
		if err := d.Set("criteria", flattenMonitorActivityLogAlertCriteria(alert.Condition)); err != nil {
			return err
		}
		if err := d.Set("action", flattenMonitorActivityLogAlertAction(alert.Actions)); err != nil {
			return err
		}
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmMonitorActivityLogAlertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorActivityLogAlertsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["activityLogAlerts"]

	if resp, err := client.Delete(ctx, resGroup, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting activity log alert %q (resource group %q): %+v", name, resGroup, err)
		}
	}

	return nil
}

func expandMonitorActivityLogAlertScopes(v []interface{}) *[]string {
	scopes := make([]string, 0)
	for _, scope := range v {
		scopes = append(scopes, scope.(string))
	}
	return &scopes
}

func expandMonitorActivityLogAlertCriteria(v map[string]interface{}) *insights.ActivityLogAlertAllOfCondition {
	conditions := make([]insights.ActivityLogAlertLeafCondition, 0)

	appendCondition := func(schemaName, fieldName string) {
		if val, ok := v[schemaName]; ok {
			conditions = append(conditions, insights.ActivityLogAlertLeafCondition{
				Field:  utils.String(fieldName),
				Equals: utils.String(val.(string)),
			})
		}
	}
	appendCondition("caller", "caller")
	appendCondition("category", "category")
	appendCondition("level", "level")
	appendCondition("operation_name", "operationName")
	appendCondition("resource_id", "resourceId")
	appendCondition("status", "status")
	appendCondition("sub_status", "subStatus")

	return &insights.ActivityLogAlertAllOfCondition{
		AllOf: &conditions,
	}
}

func expandMonitorActivityLogAlertAction(v []interface{}) *insights.ActivityLogAlertActionList {
	ags := make([]insights.ActivityLogAlertActionGroup, 0)
	for _, agValue := range v {
		val := agValue.(map[string]interface{})

		agID := val["action_group_id"].(string)
		props := make(map[string]*string)
		if propsValue, ok := val["webhook_properties"]; ok {
			for k, v := range propsValue.(map[string]interface{}) {
				props[k] = utils.String(v.(string))
			}
		}

		ag := insights.ActivityLogAlertActionGroup{
			ActionGroupID:     utils.String(agID),
			WebhookProperties: props,
		}
		ags = append(ags, ag)
	}
	return &insights.ActivityLogAlertActionList{
		ActionGroups: &ags,
	}
}

func flattenMonitorActivityLogAlertScopes(v *[]string) []interface{} {
	result := make([]interface{}, 0)
	if v != nil {
		for _, scope := range *v {
			result = append(result, scope)
		}
	}
	return result
}

func flattenMonitorActivityLogAlertCriteria(v *insights.ActivityLogAlertAllOfCondition) []interface{} {
	result := make(map[string]interface{})
	if v != nil && v.AllOf != nil {
		for _, condition := range *v.AllOf {
			if condition.Field != nil && condition.Equals != nil {
				switch strings.ToLower(*condition.Field) {
				case "operationname":
					result["operation_name"] = *condition.Equals
				case "resourceid":
					result["resource_id"] = *condition.Equals
				case "substatus":
					result["sub_status"] = *condition.Equals
				case "caller", "category", "level", "status":
					result[*condition.Field] = *condition.Equals
				}
			}
		}
	}
	return []interface{}{result}
}

func flattenMonitorActivityLogAlertAction(v *insights.ActivityLogAlertActionList) []interface{} {
	result := make([]interface{}, 0)
	if v != nil && v.ActionGroups != nil {
		for _, ag := range *v.ActionGroups {
			val := make(map[string]interface{}, 0)

			if ag.ActionGroupID != nil {
				val["action_group_id"] = *ag.ActionGroupID
			}

			props := make(map[string]string)
			for k, v := range ag.WebhookProperties {
				if v != nil {
					props[k] = *v
				}
			}
			val["webhook_properties"] = props

			result = append(result, val)
		}
	}
	return result
}
