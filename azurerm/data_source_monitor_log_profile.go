package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMonitorLogProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmLogProfileRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"servicebus_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"retention_policy": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmLogProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.LogProfilesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Log Profile %q was not found", name)
		}
		return fmt.Errorf("Error reading Log Profile: %+v", err)
	}

	d.SetId(*resp.ID)

	if props := resp.LogProfileProperties; props != nil {
		d.Set("storage_account_id", props.StorageAccountID)
		d.Set("servicebus_rule_id", props.ServiceBusRuleID)
		d.Set("categories", props.Categories)

		if err := d.Set("locations", flattenAzureRmLogProfileLocations(props.Locations)); err != nil {
			return fmt.Errorf("Error setting `locations`: %+v", err)
		}

		if err := d.Set("retention_policy", flattenAzureRmLogProfileRetentionPolicy(props.RetentionPolicy)); err != nil {
			return fmt.Errorf("Error setting `retention_policy`: %+v", err)
		}
	}

	return nil
}
