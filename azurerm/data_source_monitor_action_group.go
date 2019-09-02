package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMonitorActionGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMonitorActionGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"short_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"email_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"sms_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phone_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"webhook_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmMonitorActionGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.ActionGroupsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Action Group %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on Action Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	if group := resp.ActionGroup; group != nil {
		d.Set("short_name", group.GroupShortName)
		d.Set("enabled", group.Enabled)

		if err = d.Set("email_receiver", flattenMonitorActionGroupEmailReceiver(group.EmailReceivers)); err != nil {
			return fmt.Errorf("Error setting `email_receiver`: %+v", err)
		}

		if err = d.Set("sms_receiver", flattenMonitorActionGroupSmsReceiver(group.SmsReceivers)); err != nil {
			return fmt.Errorf("Error setting `sms_receiver`: %+v", err)
		}

		if err = d.Set("webhook_receiver", flattenMonitorActionGroupWebHookReceiver(group.WebhookReceivers)); err != nil {
			return fmt.Errorf("Error setting `webhook_receiver`: %+v", err)
		}
	}

	return nil
}
