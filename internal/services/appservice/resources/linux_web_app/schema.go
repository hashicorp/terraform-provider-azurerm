// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package linux_web_app

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func SiteConfigSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Description: "Is this Linux Web App is Always On enabled?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"api_management_api_id": {
					Description: "The ID of the API Management API for this Linux Web App.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"api_definition_url": {
					Description: "The ID of the APIM configuration for this Linux Web App.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"app_command_line": {
					Description: "The command line used to launch this app.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"application_stack": helpers.LinuxApplicationStackSchemaComputed(),

				"auto_heal_setting": autoHealSettingSchemaLinuxComputed(),

				"container_registry_use_managed_identity": {
					Description: "Do connections for Azure Container Registry use Managed Identity?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"container_registry_managed_identity_client_id": {
					Description: "The Client ID of the Managed Service Identity used for connections to the Azure Container Registry.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"default_documents": {
					Description: "The list of Default Documents for the Linux Web App.",
					Type:        pluginsdk.TypeList,
					Computed:    true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"http2_enabled": {
					Description: "Is HTTP2.0 enabled?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"ip_restriction": helpers.IpRestrictionSchemaComputed(),

				"ip_restriction_default_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_use_main_ip_restriction": {
					Description: "Is the Linux Web App `ip_restriction` configuration used for the SCM also?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"scm_ip_restriction": helpers.IpRestrictionSchemaComputed(),

				"scm_ip_restriction_default_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"local_mysql_enabled": {
					Description: "Is the Local MySQL enabled?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"load_balancing_mode": {
					Description: "The site Load Balancing Mode.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"managed_pipeline_mode": {
					Description: "The Managed Pipeline Mode.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"remote_debugging_enabled": {
					Description: "Is Remote Debugging enabled?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"remote_debugging_version": {
					Description: "The Remote Debugging Version.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"scm_type": {
					Description: "The Source Control Management Type in use.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"use_32_bit_worker": {
					Description: "Does the Linux Web App use a 32-bit worker?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"websockets_enabled": {
					Description: "Are Web Sockets enabled?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"ftps_state": {
					Description: "The State of FTP / FTPS service.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"health_check_path": {
					Description: "The path to the Health Check endpoint.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"health_check_eviction_time_in_min": {
					Description: "The amount of time in minutes that a node can be unhealthy before being removed from the load balancer.",
					Type:        pluginsdk.TypeInt,
					Computed:    true,
				},

				"worker_count": {
					Description: "The number of Workers for this Linux App Service.",
					Type:        pluginsdk.TypeInt,
					Computed:    true,
				},

				"minimum_tls_version": {
					Description: "The Minimum version of TLS for requests.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"scm_minimum_tls_version": {
					Description: "The Minimum version of TLS for requests to SCM.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"minimum_tls_cipher_suite": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"cors": helpers.CorsSettingsSchemaComputed(),

				"detailed_error_logging_enabled": {
					Description: "Is Detailed Error Logging enabled?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"linux_fx_version": {
					Description: "The `LinuxFXVersion` string.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"vnet_route_all_enabled": {
					Description: "Are all outbound traffic to NAT Gateways, Network Security Groups and User Defined Routes applied?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},
			},
		},
	}
}

func autoHealSettingSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"trigger": autoHealTriggerSchemaLinuxComputed(),

				"action": autoHealActionSchemaLinuxComputed(),
			},
		},
	}
}

func autoHealActionSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action_type": {
					Description: "The predefined action to be taken to an Auto Heal trigger.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"minimum_process_execution_time": {
					Description: "The minimum amount of time in `hh:mm:ss` the Linux Web App must have been running before the defined action will be run in the event of a trigger.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},
			},
		},
	}
}

func autoHealTriggerSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"requests": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"count": {
								Description: "The number of requests in the specified `interval` to trigger this rule.",
								Type:        pluginsdk.TypeInt,
								Computed:    true,
							},

							"interval": {
								Description: "The interval in `hh:mm:ss`.",
								Type:        pluginsdk.TypeString,
								Computed:    true,
							},
						},
					},
				},

				"status_code": {
					Type:     pluginsdk.TypeSet,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
								Description: "The status code or range for this rule.",
								Type:        pluginsdk.TypeString,
								Computed:    true,
							},

							"count": {
								Description: "The number of occurrences of the defined `status_code` in the specified `interval` on which to trigger this rule.",
								Type:        pluginsdk.TypeInt,
								Computed:    true,
							},

							"interval": {
								Description: "The time interval in the form `hh:mm:ss`.",
								Type:        pluginsdk.TypeString,
								Computed:    true,
							},

							"sub_status": {
								Description: "The Request Sub Status of the Status Code.",
								Type:        pluginsdk.TypeInt,
								Computed:    true,
							},

							"win32_status_code": {
								Description: "The Win32 Status Code of the Request.",
								Type:        pluginsdk.TypeInt,
								Computed:    true,
							},

							"path": {
								Description: "The path to which this rule status code applies.",
								Type:        pluginsdk.TypeString,
								Computed:    true,
							},
						},
					},
				},

				"slow_request": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Description: "The amount of time that qualifies as slow for this rule.",
								Type:        pluginsdk.TypeString,
								Computed:    true,
							},

							"interval": {
								Description: "The time interval.",
								Type:        pluginsdk.TypeString,
								Computed:    true,
							},

							"count": {
								Description: "The number of requests within the interval at which to trigger.",
								Type:        pluginsdk.TypeInt,
								Computed:    true,
							},
						},
					},
				},

				"slow_request_with_path": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Description: "The threshold of time passed to qualify as a Slow Request in `hh:mm:ss`.",
								Type:        pluginsdk.TypeString,
								Computed:    true,
							},

							"interval": {
								Description: "The time interval in the form `hh:mm:ss`.",
								Type:        pluginsdk.TypeString,
								Computed:    true,
							},

							"count": {
								Description: "The number of Slow Requests in the time `interval` to trigger this rule.",
								Type:        pluginsdk.TypeInt,
								Computed:    true,
							},

							"path": {
								Description: "The path for which this slow request rule applies.",
								Type:        pluginsdk.TypeString,
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}
