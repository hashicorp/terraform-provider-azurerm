// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceArmSignalRServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the SignalR service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "Which tier to use.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Free_F1",
							"Standard_S1",
							"Premium_P1",
							"Premium_P2",
						}, false),
					},

					"capacity": {
						Description: "The number of units associated with this SignalR service.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
						ValidateFunc: validation.IntInSlice([]int{
							1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 200,
							300, 400, 500, 600, 700, 800, 900, 1000,
						}),
					},
				},
			},
		},

		"connectivity_logs_enabled": {
			Description: "Whether Connectivity Logs are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"messaging_logs_enabled": {
			Description: "Whether Messaging Logs are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"http_request_logs_enabled": {
			Description: "Whether HTTP Request Logs are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"live_trace_enabled": { // azignore:AZP003 - deprecated in favor of `live_trace` and not added to the data source
			Type:       pluginsdk.TypeBool,
			Optional:   true,
			Default:    false,
			Deprecated: "`live_trace_enabled` has been deprecated in favour of `live_trace` and will be removed in 4.0.",
		},

		"live_trace": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Description: "Whether the live trace is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"connectivity_logs_enabled": {
						Description: "Whether the log category `ConnectivityLogs` is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"messaging_logs_enabled": {
						Description: "Whether the log category `MessagingLogs` is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"http_request_logs_enabled": {
						Description: "Whether the log category `HttpRequestLogs` is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Description: "Whether to enable public network access.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"local_auth_enabled": {
			Description: "Whether to enable local auth.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"aad_auth_enabled": {
			Description: "Whether to enable AAD auth.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"tls_client_cert_enabled": {
			Description: "Whether to request client certificate during TLS handshake.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"serverless_connection_timeout_in_seconds": {
			Description: "The client connection timeout.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     30,
		},

		"service_mode": {
			Description: "The service mode.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "Default",
			ValidateFunc: validation.StringInSlice([]string{
				"Serverless",
				"Classic",
				"Default",
			}, false),
		},

		"upstream_endpoint": {
			Description: "Using this block requires the SignalR service to be Serverless. When creating multiple blocks they will be processed in the order they are defined in.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category_pattern": {
						Description: "The categories to match on, or `*` for all.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"event_pattern": {
						Description: "The events to match on, or `*` for all.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"hub_pattern": {
						Description: "The hubs to match on, or `*` for all.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"url_template": {
						Description:  "The upstream URL Template. This can be a URL or a template such as `http://host.com/{hub}/api/{category}/{event}`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.UrlTemplate,
					},

					"user_assigned_identity_id": {
						Description:  "The Managed Identity IDs to be assigned to this signalR upstream setting by using resource UUID as both system assigned and User Assigned Identity is supported.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"cors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Description: "A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"hostname": {
			Description: "The FQDN of the SignalR service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ip_address": {
			Description: "The publicly accessible IP of the SignalR service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"public_port": {
			Description: "The publicly accessible port of the SignalR service which is designed for browser/client use.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"server_port": {
			Description: "The publicly accessible port of the SignalR service which is designed for customer server side use.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"primary_access_key": {
			Description: "The primary access key for the SignalR service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"primary_connection_string": {
			Description: "The primary connection string for the SignalR service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_access_key": {
			Description: "The secondary access key for the SignalR service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"secondary_connection_string": {
			Description: "The secondary connection string for the SignalR service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"tags": commonschema.Tags(),
	}
}
