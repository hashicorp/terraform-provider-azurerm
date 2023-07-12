// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMonitorActionGroup() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceMonitorActionGroupRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"short_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"email_receiver": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"connection_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"ticket_configuration": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"region": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"phone_number": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"service_uri": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"automation_account_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"runbook_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"webhook_resource_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"is_global_runbook": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"service_uri": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"phone_number": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"callback_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"function_app_resource_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"function_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"http_trigger_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
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
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"role_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["event_hub_receiver"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"event_hub_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: eventhubs.ValidateEventhubID,
						Deprecated:   "This property is deprecated and will be removed in version 4.0 of the provider, please use 'event_hub_name' and 'event_hub_namespace' instead.",
					},
					"event_hub_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"event_hub_namespace": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsUUID,
					},
					"use_common_alert_schema": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"subscription_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		}
	} else {
		resource.Schema["event_hub_receiver"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"event_hub_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"event_hub_namespace": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsUUID,
					},
					"use_common_alert_schema": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"subscription_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		}
	}
	return resource
}

func dataSourceMonitorActionGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			if err = d.Set("event_hub_receiver", flattenMonitorActionGroupEventHubReceiver(resourceGroup, props.EventHubReceivers)); err != nil {
				return fmt.Errorf("setting `event_hub_receiver`: %+v", err)
			}
		}
	}
	return nil
}
