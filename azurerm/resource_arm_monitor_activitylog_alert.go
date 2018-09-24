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
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"level": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"operation_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"resource_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"sub_status": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.NoZeroValues,
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

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := insights.ActivityLogAlertResource{
		Location:         utils.String(azureRMNormalizeLocation("Global")),
		ActivityLogAlert: &insights.ActivityLogAlert{},
		Tags:             expandedTags,
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
