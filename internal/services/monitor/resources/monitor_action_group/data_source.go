// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_action_group

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access the properties of an Action Group.",
		Read:        dataSourceMonitorActionGroupRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Action Group.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"short_name": {
				Description: "The short name of the action group.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"enabled": {
				Description: "Whether this action group is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"email_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the email receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"email_address": {
							Description: "The email address of this receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"use_common_alert_schema": {
							Description: "Whether the common alert schema is used.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"itsm_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the ITSM receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"workspace_id": {
							Description: "The Azure Log Analytics Workspace ID where this connection is defined.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"connection_id": {
							Description: "The unique connection identifier of the ITSM connection.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"ticket_configuration": {
							Description: "A JSON blob for the configurations of the ITSM action. CreateMultipleWorkItems option will be part of this blob as well.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"region": {
							Description: "The region of the workspace.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"azure_app_push_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the Azure app push receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"email_address": {
							Description: "The email address of the user signed into the mobile app who will receive push notifications from this receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"sms_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the SMS receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"country_code": {
							Description: "The country code of the SMS receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"phone_number": {
							Description: "The phone number of the SMS receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"webhook_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the webhook receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"service_uri": {
							Description: "The URI where webhooks should be sent.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"use_common_alert_schema": {
							Description: "Whether the common alert schema is used.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
						"aad_auth": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"object_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"identifier_uri": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"tenant_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"automation_runbook_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the Automation Runbook receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"automation_account_id": {
							Description: "The Automation Account ID which holds this runbook and authenticates to Azure resources.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"runbook_name": {
							Description: "The name for this runbook.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"webhook_resource_id": {
							Description: "The resource ID for webhook linked to this runbook.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"is_global_runbook": {
							Description: "Whether this instance is global runbook.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
						"service_uri": {
							Description: "The URI where webhooks should be sent.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"use_common_alert_schema": {
							Description: "Whether the common alert schema is used.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"voice_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the voice receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"country_code": {
							Description: "The country code of the voice receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"phone_number": {
							Description: "The phone number of the voice receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"logic_app_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the Logic App receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"resource_id": {
							Description: "The Azure resource ID of the Logic App.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"callback_url": {
							Description: "The callback URL where HTTP request sent to.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"use_common_alert_schema": {
							Description: "Whether the common alert schema is used.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"azure_function_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the Azure Function receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"function_app_resource_id": {
							Description: "The Azure resource ID of the Function App.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"function_name": {
							Description: "The function name in the Function App.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"http_trigger_url": {
							Description: "The HTTP trigger URL where HTTP request sent to.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"use_common_alert_schema": {
							Description: "Whether the common alert schema is used.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"arm_role_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the ARM role receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"role_id": {
							Description: "The arm role ID.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"use_common_alert_schema": {
							Description: "Whether the common alert schema is used.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"event_hub_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the EventHub Receiver, must be unique within action group.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"event_hub_name": {
							Description: "The name of the specific Event Hub queue.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"event_hub_namespace": {
							Description: "The namespace name of the Event Hub.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"tenant_id": {
							Description: "The Tenant ID for the subscription containing this Event Hub.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
						"use_common_alert_schema": {
							Description: "Whether the common alert schema is used.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
						"subscription_id": {
							Description: "The ID for the subscription containing this Event Hub.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMonitorActionGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.ActionGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)

	id := actiongroupsapis.NewActionGroupID(subscriptionId, resourceGroup, d.Get("name").(string))

	resp, err := client.ActionGroupsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}
	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("short_name", props.GroupShortName)
			d.Set("enabled", props.Enabled)

			if err = d.Set("email_receiver", flattenMonitorActionGroupEmailReceiver(props.EmailReceivers)); err != nil {
				return fmt.Errorf("setting `email_receiver`: %+v", err)
			}

			if err = d.Set("itsm_receiver", flattenMonitorActionGroupItsmReceiver(props.ItsmReceivers)); err != nil {
				return fmt.Errorf("setting `itsm_receiver`: %+v", err)
			}

			if err = d.Set("azure_app_push_receiver", flattenMonitorActionGroupAzureAppPushReceiver(props.AzureAppPushReceivers)); err != nil {
				return fmt.Errorf("setting `azure_app_push_receiver`: %+v", err)
			}

			if err = d.Set("sms_receiver", flattenMonitorActionGroupSmsReceiver(props.SmsReceivers)); err != nil {
				return fmt.Errorf("setting `sms_receiver`: %+v", err)
			}

			if err = d.Set("webhook_receiver", flattenMonitorActionGroupWebHookReceiver(props.WebhookReceivers)); err != nil {
				return fmt.Errorf("setting `webhook_receiver`: %+v", err)
			}

			if err = d.Set("automation_runbook_receiver", flattenMonitorActionGroupAutomationRunbookReceiver(props.AutomationRunbookReceivers)); err != nil {
				return fmt.Errorf("setting `automation_runbook_receiver`: %+v", err)
			}

			if err = d.Set("voice_receiver", flattenMonitorActionGroupVoiceReceiver(props.VoiceReceivers)); err != nil {
				return fmt.Errorf("setting `voice_receiver`: %+v", err)
			}

			if err = d.Set("logic_app_receiver", flattenMonitorActionGroupLogicAppReceiver(props.LogicAppReceivers)); err != nil {
				return fmt.Errorf("setting `logic_app_receiver`: %+v", err)
			}

			if err = d.Set("azure_function_receiver", flattenMonitorActionGroupAzureFunctionReceiver(props.AzureFunctionReceivers)); err != nil {
				return fmt.Errorf("setting `azure_function_receiver`: %+v", err)
			}
			if err = d.Set("arm_role_receiver", flattenMonitorActionGroupRoleReceiver(props.ArmRoleReceivers)); err != nil {
				return fmt.Errorf("setting `arm_role_receiver`: %+v", err)
			}
			if err = d.Set("event_hub_receiver", flattenMonitorActionGroupEventHubReceiver(props.EventHubReceivers)); err != nil {
				return fmt.Errorf("setting `event_hub_receiver`: %+v", err)
			}
		}
	}
	return nil
}
