// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinuxWebAppSlotV0toV1 struct{}

var _ pluginsdk.StateUpgrade = LinuxWebAppSlotV0toV1{}

func (l LinuxWebAppSlotV0toV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
		},

		"service_plan_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auth_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"additional_login_parameters": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"allowed_external_redirect_urls": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"default_provider": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // Once set, cannot be unset
					},

					"issuer": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"runtime_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"token_refresh_extension_hours": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
						Default:  72,
					},

					"token_store_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"unauthenticated_client_action": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // Once set, cannot be removed
					},

					"active_directory": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_secret": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"allowed_audiences": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"facebook": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"app_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"app_secret": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},

								"app_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"oauth_scopes": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"github": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_secret": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"oauth_scopes": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"google": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_secret": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"oauth_scopes": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"microsoft": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_secret": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"oauth_scopes": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"twitter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"consumer_key": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"consumer_secret": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},

								"consumer_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},

		"auth_settings_v2": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"auth_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"runtime_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  `~1`,
					},

					"config_file_path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"require_authentication": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"unauthenticated_action": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "RedirectToLoginPage",
					},

					"default_provider": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"excluded_paths": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"apple_v2": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"login_scopes": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"active_directory_v2": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"tenant_auth_endpoint": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"client_secret_certificate_thumbprint": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"jwt_allowed_groups": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"jwt_allowed_client_applications": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"www_authentication_disabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"allowed_groups": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"allowed_identities": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"allowed_applications": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"login_parameters": {
									Type:     pluginsdk.TypeMap,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"allowed_audiences": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"azure_static_web_app_v2": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},

					"custom_oidc_v2": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"openid_configuration_endpoint": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"name_claim_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"scopes": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"client_credential_method": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"authorisation_endpoint": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"token_endpoint": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"issuer_endpoint": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"certification_uri": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"facebook_v2": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						AtLeastOneOf: []string{
							"auth_settings_v2.0.apple_v2",
							"auth_settings_v2.0.active_directory_v2",
							"auth_settings_v2.0.azure_static_web_app_v2",
							"auth_settings_v2.0.custom_oidc_v2",
							"auth_settings_v2.0.facebook_v2",
							"auth_settings_v2.0.github_v2",
							"auth_settings_v2.0.google_v2",
							"auth_settings_v2.0.microsoft_v2",
							"auth_settings_v2.0.twitter_v2",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"app_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"app_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"graph_api_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},

								"login_scopes": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"github_v2": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"login_scopes": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"google_v2": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"allowed_audiences": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"login_scopes": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"microsoft_v2": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"client_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"allowed_audiences": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"login_scopes": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"twitter_v2": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"consumer_key": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"consumer_secret_setting_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},

					"login": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"logout_endpoint": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"token_store_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"token_refresh_extension_time": {
									Type:     pluginsdk.TypeFloat,
									Optional: true,
									Default:  72,
								},

								"token_store_path": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"token_store_sas_setting_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"preserve_url_fragments_for_logins": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"allowed_external_redirect_urls": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"cookie_expiration_convention": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "FixedTime",
								},

								"cookie_expiration_time": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "08:00:00",
								},

								"validate_nonce": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},

								"nonce_expiration_time": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "00:05:00",
								},
							},
						},
					},

					"require_https": {
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
						Description: "Should HTTPS be required on connections? Defaults to true.",
					},

					"http_route_api_prefix": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "/.auth",
					},

					"forward_proxy_convention": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "NoProxy",
					},

					"forward_proxy_custom_host_header_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"forward_proxy_custom_scheme_header_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"backup": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"storage_account_url": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"schedule": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"frequency_interval": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},

								"frequency_unit": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"keep_at_least_one_backup": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"retention_period_days": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  30,
								},

								"start_time": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},

								"last_execution_time": {
									Type:        pluginsdk.TypeString,
									Computed:    true,
									Description: "The time the backup was last attempted.",
								},
							},
						},
					},
				},
			},
		},

		"client_affinity_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_certificate_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_certificate_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "Optional",
		},

		"client_certificate_exclusion_paths": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"connection_string": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"value": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"identity": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"identity_ids": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"principal_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"key_vault_reference_identity_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"logs": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"application_logs": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"file_system_level": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"azure_blob_storage": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"level": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"sas_url": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"retention_in_days": {
												Type:     pluginsdk.TypeInt,
												Required: true,
											},
										},
									},
								},
							},
						},
					},

					"http_logs": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"file_system": {
									Type:          pluginsdk.TypeList,
									Optional:      true,
									MaxItems:      1,
									ConflictsWith: []string{"logs.0.http_logs.0.azure_blob_storage"},
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"retention_in_mb": {
												Type:     pluginsdk.TypeInt,
												Required: true,
											},

											"retention_in_days": {
												Type:     pluginsdk.TypeInt,
												Required: true,
											},
										},
									},
								},

								"azure_blob_storage": {
									Type:          pluginsdk.TypeList,
									Optional:      true,
									MaxItems:      1,
									ConflictsWith: []string{"logs.0.http_logs.0.file_system"},
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"sas_url": {
												Type:      pluginsdk.TypeString,
												Required:  true,
												Sensitive: true,
											},
											"retention_in_days": {
												Type:     pluginsdk.TypeInt,
												Optional: true,
												Default:  0,
											},
										},
									},
								},
							},
						},
					},

					"failed_request_tracing": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"detailed_error_messages": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"site_config": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"always_on": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},

					"api_management_api_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"api_definition_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"app_command_line": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"application_stack": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dotnet_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"go_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"php_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"python_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"node_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"ruby_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"java_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"java_server": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"java_server_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"docker_image": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"docker_image_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"docker_image_tag": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"docker_registry_url": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},

								"docker_registry_username": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},

								"docker_registry_password": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Computed:  true,
									Sensitive: true,
								},
							},
						},
					},

					"auto_heal_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"auto_heal_setting": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"trigger": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"requests": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:     pluginsdk.TypeInt,
															Required: true,
														},

														"interval": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},
													},
												},
											},

											"status_code": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"status_code_range": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},

														"count": {
															Type:     pluginsdk.TypeInt,
															Required: true,
														},

														"interval": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},

														"sub_status": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},

														"win32_status_code": {
															Type:     pluginsdk.TypeInt,
															Optional: true,
														},

														"path": {
															Type:     pluginsdk.TypeString,
															Optional: true,
														},
													},
												},
											},

											"slow_request": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"time_taken": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},

														"interval": {
															Type:     pluginsdk.TypeString,
															Required: true,
														},

														"count": {
															Type:     pluginsdk.TypeInt,
															Required: true,
														},

														"path": {
															Type:     pluginsdk.TypeString,
															Optional: true,
														},
													},
												},
											},
										},
									},
								},

								"action": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"action_type": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},

											"minimum_process_execution_time": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},

					"container_registry_use_managed_identity": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"container_registry_managed_identity_client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"default_documents": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"http2_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"ip_restriction": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ip_address": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"service_tag": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"virtual_network_subnet_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},

								"priority": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  65000,
								},

								"action": {
									Type:     pluginsdk.TypeString,
									Default:  "Allow",
									Optional: true,
								},

								"headers": {
									Type:       pluginsdk.TypeList,
									MaxItems:   1,
									Optional:   true,
									ConfigMode: pluginsdk.SchemaConfigModeAttr,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"x_forwarded_host": {
												Type:     pluginsdk.TypeList,
												MaxItems: 8,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"x_forwarded_for": {
												Type:     pluginsdk.TypeList,
												MaxItems: 8,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"x_azure_fdid": { // Front Door ID (UUID)
												Type:     pluginsdk.TypeList,
												MaxItems: 8,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"x_fd_health_probe": { // 1 or absent
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
							},
						},
					},

					"scm_use_main_ip_restriction": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"scm_ip_restriction": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ip_address": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"service_tag": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"virtual_network_subnet_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
								},

								"priority": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  65000,
								},

								"action": {
									Type:     pluginsdk.TypeString,
									Default:  "Allow",
									Optional: true,
								},

								"headers": {
									Type:       pluginsdk.TypeList,
									MaxItems:   1,
									Optional:   true,
									ConfigMode: pluginsdk.SchemaConfigModeAttr,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"x_forwarded_host": {
												Type:     pluginsdk.TypeList,
												MaxItems: 8,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"x_forwarded_for": {
												Type:     pluginsdk.TypeList,
												MaxItems: 8,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"x_azure_fdid": { // Front Door ID (UUID)
												Type:     pluginsdk.TypeList,
												MaxItems: 8,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"x_fd_health_probe": { // 1 or absent
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
							},
						},
					},

					"local_mysql_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"load_balancing_mode": { // Supported on Function Apps?
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "LeastRequests",
					},

					"managed_pipeline_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "Integrated",
					},

					"remote_debugging_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"remote_debugging_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"scm_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"use_32_bit_worker": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"websockets_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"ftps_state": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "Disabled",
					},

					"health_check_path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"health_check_eviction_time_in_min": { // NOTE: Will evict the only node in single node configurations.
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"worker_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"minimum_tls_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "1.2",
					},

					"scm_minimum_tls_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "1.2",
					},

					"cors": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allowed_origins": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"support_credentials": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},

					"auto_swap_slot_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"vnet_route_all_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"detailed_error_logging_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"linux_fx_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"storage_account": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"share_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"access_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},

					"mount_path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"virtual_network_subnet_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"zip_deploy_file": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func (l LinuxWebAppSlotV0toV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId, ok := rawState["service_plan_id"].(string)
		if !ok || oldId == "" {
			return rawState, nil
		}
		parsedId, err := commonids.ParseAppServicePlanIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := parsedId.ID()
		rawState["service_plan_id"] = newId
		return rawState, nil
	}
}
