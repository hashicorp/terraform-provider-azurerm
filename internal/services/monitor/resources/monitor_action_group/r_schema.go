// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_action_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/automationaccount"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorActionGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Action Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": {
			Description: "The Azure Region where the Action Group should exist.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "global",
			ValidateFunc: validation.Any(
				location.EnhancedValidate,
				validation.StringInSlice([]string{
					"global",
				}, false),
			),
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
		},

		"short_name": {
			Description:  "The short name of the action group. This will be used in SMS messages.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 12),
		},

		"enabled": {
			Description: "Whether this action group is enabled. If an action group is not enabled, then none of its receivers will receive communications.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"email_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the email receiver. Names must be unique (case-insensitive) across all receivers within an action group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"email_address": {
						Description:  "The email address of this receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"use_common_alert_schema": {
						Description: "Enables or disables the common alert schema.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
				},
			},
		},

		"itsm_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the ITSM receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"workspace_id": {
						Description:  "The Azure Log Analytics Workspace ID where this connection is defined. Format is `<subscription id>|<workspace id>`, for example `00000000-0000-0000-0000-000000000000|00000000-0000-0000-0000-000000000000`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.WorkspaceID,
					},
					"connection_id": {
						Description:  "The unique connection identifier of the ITSM connection.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"ticket_configuration": {
						Description:      "A JSON blob for the configurations of the ITSM action. CreateMultipleWorkItems option will be part of this blob as well.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
					},
					"region": {
						Description:      "The region of the workspace.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validation.StringIsNotEmpty,
						DiffSuppressFunc: location.DiffSuppressFunc,
					},
				},
			},
		},

		"azure_app_push_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Azure app push receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"email_address": {
						Description:  "The email address of the user signed into the mobile app who will receive push notifications from this receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"sms_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the SMS receiver. Names must be unique (case-insensitive) across all receivers within an action group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"country_code": {
						Description:  "The country code of the SMS receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"phone_number": {
						Description:  "The phone number of the SMS receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"webhook_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the webhook receiver. Names must be unique (case-insensitive) across all receivers within an action group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"service_uri": {
						Description:  "The URI where webhooks should be sent.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},
					"use_common_alert_schema": {
						Description: "Enables or disables the common alert schema.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"aad_auth": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"object_id": {
									Description:  "The webhook application object ID for AAD auth.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsUUID,
								},

								"identifier_uri": {
									Description:  "The identifier URI for AAD auth.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Computed:     true, // azignore:AZS007 - pre-existing violation
									ValidateFunc: validation.IsURLWithScheme([]string{"api", "https"}),
								},

								"tenant_id": {
									Description:  "The tenant ID for AAD auth.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Computed:     true, // azignore:AZS007 - pre-existing violation
									ValidateFunc: validation.IsUUID,
								},
							},
						},
					},
				},
			},
		},

		"automation_runbook_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Automation Runbook receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"automation_account_id": {
						Description:  "The Automation Account ID which holds this runbook and authenticates to Azure resources.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: automationaccount.ValidateAutomationAccountID,
					},
					"runbook_name": {
						Description:  "The name for this runbook.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"webhook_resource_id": {
						Description:  "The resource ID for webhook linked to this runbook.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"is_global_runbook": {
						Description: "Whether this instance is global runbook.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"service_uri": {
						Description:  "The URI where webhooks should be sent.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},
					"use_common_alert_schema": {
						Description: "Enables or disables the common alert schema.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"voice_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the voice receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"country_code": {
						Description:  "The country code of the voice receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"phone_number": {
						Description:  "The phone number of the voice receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"logic_app_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Logic App receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"resource_id": {
						Description:  "The Azure resource ID of the Logic App.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"callback_url": {
						Description:  "The callback URL where HTTP request sent to.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},
					"use_common_alert_schema": {
						Description: "Enables or disables the common alert schema.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
				},
			},
		},

		"azure_function_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Azure Function receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					// TODO: this should be `_id` and not `_resource_id`
					"function_app_resource_id": {
						Description:  "The Azure resource ID of the Function App.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateFunctionAppID,
					},
					"function_name": {
						Description:  "The function name in the Function App.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"http_trigger_url": {
						Description:  "The HTTP trigger URL where HTTP request sent to.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},
					"use_common_alert_schema": {
						Description: "Enables or disables the common alert schema.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
				},
			},
		},

		"arm_role_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the ARM role receiver.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"role_id": {
						Description:  "The arm role ID.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"use_common_alert_schema": {
						Description: "Enables or disables the common alert schema.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
				},
			},
		},

		"event_hub_receiver": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the EventHub Receiver, must be unique within action group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"event_hub_name": {
						Description:  "The name of the specific Event Hub queue.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"event_hub_namespace": {
						Description:  "The namespace name of the Event Hub.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tenant_id": {
						Description:  "The Tenant ID for the subscription containing this Event Hub.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsUUID,
					},
					"use_common_alert_schema": {
						Description: "Whether to use common alert schema.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"subscription_id": {
						Description:  "The ID for the subscription containing this Event Hub.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
