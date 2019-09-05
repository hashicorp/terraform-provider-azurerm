package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"email_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"email_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"itsm_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"workspace_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.UUID,
						},
						"connection_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.UUID,
						},
						"ticket_configuration": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validation.ValidateJsonString,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},
						"region": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"azure_app_push_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"email_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"country_code": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"phone_number": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"service_uri": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.URLIsHTTPOrHTTPS,
						},
					},
				},
			},

			"automation_runbook_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"automation_account_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.UUID,
						},
						"runbook_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"webhook_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"is_global_runbook": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"service_uri": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.URLIsHTTPOrHTTPS,
						},
					},
				},
			},

			"voice_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"country_code": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"phone_number": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"logic_app_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"callback_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.URLIsHTTPOrHTTPS,
						},
					},
				},
			},

			"azure_function_receiver": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"function_app_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"function_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"http_trigger_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.URLIsHTTPOrHTTPS,
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

		if err = d.Set("itsm_receiver", flattenMonitorActionGroupItsmReceiver(group.ItsmReceivers)); err != nil {
			return fmt.Errorf("Error setting `itsm_receiver`: %+v", err)
		}

		if err = d.Set("azure_app_push_receiver", flattenMonitorActionGroupAzureAppPushReceiver(group.AzureAppPushReceivers)); err != nil {
			return fmt.Errorf("Error setting `azure_app_push_receiver`: %+v", err)
		}

		if err = d.Set("sms_receiver", flattenMonitorActionGroupSmsReceiver(group.SmsReceivers)); err != nil {
			return fmt.Errorf("Error setting `sms_receiver`: %+v", err)
		}

		if err = d.Set("webhook_receiver", flattenMonitorActionGroupWebHookReceiver(group.WebhookReceivers)); err != nil {
			return fmt.Errorf("Error setting `webhook_receiver`: %+v", err)
		}

		if err = d.Set("automation_runbook_receiver", flattenMonitorActionGroupAutomationRunbookReceiver(group.AutomationRunbookReceivers)); err != nil {
			return fmt.Errorf("Error setting `automation_runbook_receiver`: %+v", err)
		}

		if err = d.Set("voice_receiver", flattenMonitorActionGroupVoiceReceiver(group.VoiceReceivers)); err != nil {
			return fmt.Errorf("Error setting `voice_receiver`: %+v", err)
		}

		if err = d.Set("logic_app_receiver", flattenMonitorActionGroupLogicAppReceiver(group.LogicAppReceivers)); err != nil {
			return fmt.Errorf("Error setting `logic_app_receiver`: %+v", err)
		}

		if err = d.Set("azure_function_receiver", flattenMonitorActionGroupAzureFunctionReceiver(group.AzureFunctionReceivers)); err != nil {
			return fmt.Errorf("Error setting `azure_function_receiver`: %+v", err)
		}
	}

	return nil
}
