package web

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func schemaAppServiceAadAuthSettings() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
				"allowed_audiences": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
			},
		},
	}
}

func schemaAppServiceFacebookAuthSettings() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
					Required:  true,
					Sensitive: true,
				},
				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
			},
		},
	}
}

func schemaAppServiceGoogleAuthSettings() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
					Required:  true,
					Sensitive: true,
				},
				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
			},
		},
	}
}

func schemaAppServiceMicrosoftAuthSettings() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
					Required:  true,
					Sensitive: true,
				},
				"oauth_scopes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
			},
		},
	}
}

func schemaAppServiceTwitterAuthSettings() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
					Required:  true,
					Sensitive: true,
				},
			},
		},
	}
}

func schemaAppServiceAuthSettings() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},

				"additional_login_params": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"allowed_external_redirect_urls": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},

				"default_provider": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.BuiltInAuthenticationProviderAzureActiveDirectory),
						string(web.BuiltInAuthenticationProviderFacebook),
						// TODO: add GitHub Auth when API bump merged
						// string(web.BuiltInAuthenticationProviderGithub),
						string(web.BuiltInAuthenticationProviderGoogle),
						string(web.BuiltInAuthenticationProviderMicrosoftAccount),
						string(web.BuiltInAuthenticationProviderTwitter),
					}, false),
				},

				"issuer": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
				},

				"runtime_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
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
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AllowAnonymous),
						string(web.RedirectToLoginPage),
					}, false),
				},

				"active_directory": schemaAppServiceAadAuthSettings(),

				"facebook": schemaAppServiceFacebookAuthSettings(),

				"google": schemaAppServiceGoogleAuthSettings(),

				"microsoft": schemaAppServiceMicrosoftAuthSettings(),

				"twitter": schemaAppServiceTwitterAuthSettings(),
			},
		},
	}
}

func schemaAppServiceIdentity() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"identity_ids": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validate.UserAssignedIdentityID,
					},
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedServiceIdentityTypeNone),
						string(web.ManagedServiceIdentityTypeSystemAssigned),
						string(web.ManagedServiceIdentityTypeSystemAssignedUserAssigned),
						string(web.ManagedServiceIdentityTypeUserAssigned),
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"principal_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func schemaAppServiceSiteConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},

				"dotnet_framework_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "v4.0",
					ValidateFunc: validation.StringInSlice([]string{
						"v2.0",
						"v4.0",
						"v5.0",
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": schemaAppServiceIpRestriction(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_ip_restriction": schemaAppServiceIpRestriction(),

				"java_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice([]string{"1.7", "1.8", "11"}, false),
				},

				"java_container": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"JAVA",
						"JETTY",
						"TOMCAT",
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"java_container_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"local_mysql_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.Classic),
						string(web.Integrated),
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"5.5",
						"5.6",
						"7.0",
						"7.1",
						"7.2",
						"7.3",
						"7.4",
					}, false),
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"2.7",
						"3.4",
					}, false),
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
					ValidateFunc: validation.StringInSlice([]string{
						"VS2012", // TODO for 3.0 - remove VS2012, VS2013, VS2015
						"VS2013",
						"VS2015",
						"VS2017",
						"VS2019",
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ScmTypeBitbucketGit),
						string(web.ScmTypeBitbucketHg),
						string(web.ScmTypeCodePlexGit),
						string(web.ScmTypeCodePlexHg),
						string(web.ScmTypeDropbox),
						string(web.ScmTypeExternalGit),
						string(web.ScmTypeExternalHg),
						string(web.ScmTypeGitHub),
						string(web.ScmTypeLocalGit),
						string(web.ScmTypeNone),
						string(web.ScmTypeOneDrive),
						string(web.ScmTypeTfs),
						string(web.ScmTypeVSO),
						string(web.ScmTypeVSTSRM),
					}, false),
				},

				"use_32_bit_worker_process": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AllAllowed),
						string(web.Disabled),
						string(web.FtpsOnly),
					}, false),
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"number_of_workers": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 100),
					Computed:     true,
				},

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},

				"windows_fx_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},

				"min_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.OneFullStopZero),
						string(web.OneFullStopOne),
						string(web.OneFullStopTwo),
					}, false),
				},

				"cors": SchemaWebCorsSettings(),

				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func schemaAppServiceLogsConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"application_logs": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"file_system_level": {
								Type:          pluginsdk.TypeString,
								Optional:      true,
								Default:       "Off",
								ConflictsWith: []string{"logs.0.http_logs.0.azure_blob_storage"},
								ValidateFunc: validation.StringInSlice([]string{
									string(web.Error),
									string(web.Information),
									string(web.Off),
									string(web.Verbose),
									string(web.Warning),
								}, false),
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
											ValidateFunc: validation.StringInSlice([]string{
												string(web.Error),
												string(web.Information),
												string(web.Off),
												string(web.Verbose),
												string(web.Warning),
											}, false),
										},
										"sas_url": {
											Type:      pluginsdk.TypeString,
											Required:  true,
											Sensitive: true,
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
					Computed: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"file_system": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"retention_in_mb": {
											Type:         pluginsdk.TypeInt,
											Required:     true,
											ValidateFunc: validation.IntBetween(25, 100),
										},
										"retention_in_days": {
											Type:         pluginsdk.TypeInt,
											Required:     true,
											ValidateFunc: validation.IntAtLeast(0),
										},
									},
								},
								ConflictsWith: []string{"logs.0.http_logs.0.azure_blob_storage"},
								AtLeastOneOf:  []string{"logs.0.http_logs.0.azure_blob_storage", "logs.0.http_logs.0.file_system"},
							},
							"azure_blob_storage": {
								Type:     pluginsdk.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"sas_url": {
											Type:      pluginsdk.TypeString,
											Required:  true,
											Sensitive: true,
										},
										"retention_in_days": {
											Type:     pluginsdk.TypeInt,
											Required: true,
										},
									},
								},
								ConflictsWith: []string{"logs.0.http_logs.0.file_system"},
								AtLeastOneOf:  []string{"logs.0.http_logs.0.azure_blob_storage", "logs.0.http_logs.0.file_system"},
							},
						},
					},
				},
				"detailed_error_messages_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
				"failed_request_tracing_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func schemaAppServiceStorageAccounts() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AzureBlob),
						string(web.AzureFiles),
					}, false),
				},

				"account_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"share_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"access_key": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"mount_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func schemaAppServiceDataSourceSiteConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},

				"dotnet_framework_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ip_restriction": schemaAppServiceDataSourceIpRestriction(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_ip_restriction": schemaAppServiceDataSourceIpRestriction(),

				"java_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_container": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_container_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"local_mysql_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_32_bit_worker_process": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"number_of_workers": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"windows_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"min_tls_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"cors": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"allowed_origins": {
								Type:     pluginsdk.TypeSet,
								Computed: true,
								Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							},
							"support_credentials": {
								Type:     pluginsdk.TypeBool,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func schemaAppServiceIpRestriction() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:       pluginsdk.TypeList,
		Optional:   true,
		Computed:   true,
		ConfigMode: pluginsdk.SchemaConfigModeAttr,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_address": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"service_tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"virtual_network_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"priority": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      65000,
					ValidateFunc: validation.IntBetween(1, 2147483647),
				},

				"action": {
					Type:     pluginsdk.TypeString,
					Default:  "Allow",
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Allow",
						"Deny",
					}, false),
				},

				//lintignore:XS003
				"headers": {
					Type:       pluginsdk.TypeList,
					Optional:   true,
					Computed:   true,
					MaxItems:   1,
					ConfigMode: pluginsdk.SchemaConfigModeAttr,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{

							// lintignore:S018
							"x_forwarded_host": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},

							// lintignore:S018
							"x_forwarded_for": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.IsCIDR,
								},
							},

							// lintignore:S018
							"x_azure_fdid": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.IsUUID,
								},
							},

							// lintignore:S018
							"x_fd_health_probe": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 1,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										"1",
									}, false),
								},
							},
						},
					},
				},
			},
		},
	}
}

func schemaAppServiceDataSourceIpRestriction() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_address": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"service_tag": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"virtual_network_subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"priority": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func expandAppServiceAuthSettings(input []interface{}) web.SiteAuthSettingsProperties {
	siteAuthSettingsProperties := web.SiteAuthSettingsProperties{}

	if len(input) == 0 {
		return siteAuthSettingsProperties
	}

	setting := input[0].(map[string]interface{})

	if v, ok := setting["enabled"]; ok {
		siteAuthSettingsProperties.Enabled = utils.Bool(v.(bool))
	}

	if v, ok := setting["additional_login_params"]; ok {
		input := v.(map[string]interface{})

		additionalLoginParams := make([]string, 0)
		for k, v := range input {
			additionalLoginParams = append(additionalLoginParams, fmt.Sprintf("%s=%s", k, v.(string)))
		}

		siteAuthSettingsProperties.AdditionalLoginParams = &additionalLoginParams
	}

	if v, ok := setting["allowed_external_redirect_urls"]; ok {
		input := v.([]interface{})

		allowedExternalRedirectUrls := make([]string, 0)
		for _, param := range input {
			allowedExternalRedirectUrls = append(allowedExternalRedirectUrls, param.(string))
		}

		siteAuthSettingsProperties.AllowedExternalRedirectUrls = &allowedExternalRedirectUrls
	}

	if v, ok := setting["default_provider"]; ok {
		siteAuthSettingsProperties.DefaultProvider = web.BuiltInAuthenticationProvider(v.(string))
	}

	if v, ok := setting["issuer"]; ok {
		siteAuthSettingsProperties.Issuer = utils.String(v.(string))
	}

	if v, ok := setting["runtime_version"]; ok {
		siteAuthSettingsProperties.RuntimeVersion = utils.String(v.(string))
	}

	if v, ok := setting["token_refresh_extension_hours"]; ok {
		siteAuthSettingsProperties.TokenRefreshExtensionHours = utils.Float(v.(float64))
	}

	if v, ok := setting["token_store_enabled"]; ok {
		siteAuthSettingsProperties.TokenStoreEnabled = utils.Bool(v.(bool))
	}

	if v, ok := setting["unauthenticated_client_action"]; ok {
		siteAuthSettingsProperties.UnauthenticatedClientAction = web.UnauthenticatedClientAction(v.(string))
	}

	if v, ok := setting["active_directory"]; ok {
		activeDirectorySettings := v.([]interface{})

		for _, setting := range activeDirectorySettings {
			if setting == nil {
				continue
			}

			activeDirectorySetting := setting.(map[string]interface{})

			if v, ok := activeDirectorySetting["client_id"]; ok {
				siteAuthSettingsProperties.ClientID = utils.String(v.(string))
			}

			if v, ok := activeDirectorySetting["client_secret"]; ok {
				siteAuthSettingsProperties.ClientSecret = utils.String(v.(string))
			}

			if v, ok := activeDirectorySetting["allowed_audiences"]; ok {
				input := v.([]interface{})

				allowedAudiences := make([]string, 0)
				for _, param := range input {
					allowedAudiences = append(allowedAudiences, param.(string))
				}

				siteAuthSettingsProperties.AllowedAudiences = &allowedAudiences
			}
		}
	}

	if v, ok := setting["facebook"]; ok {
		facebookSettings := v.([]interface{})

		for _, setting := range facebookSettings {
			facebookSetting := setting.(map[string]interface{})

			if v, ok := facebookSetting["app_id"]; ok {
				siteAuthSettingsProperties.FacebookAppID = utils.String(v.(string))
			}

			if v, ok := facebookSetting["app_secret"]; ok {
				siteAuthSettingsProperties.FacebookAppSecret = utils.String(v.(string))
			}

			if v, ok := facebookSetting["oauth_scopes"]; ok {
				input := v.([]interface{})

				oauthScopes := make([]string, 0)
				for _, param := range input {
					oauthScopes = append(oauthScopes, param.(string))
				}

				siteAuthSettingsProperties.FacebookOAuthScopes = &oauthScopes
			}
		}
	}

	if v, ok := setting["google"]; ok {
		googleSettings := v.([]interface{})

		for _, setting := range googleSettings {
			googleSetting := setting.(map[string]interface{})

			if v, ok := googleSetting["client_id"]; ok {
				siteAuthSettingsProperties.GoogleClientID = utils.String(v.(string))
			}

			if v, ok := googleSetting["client_secret"]; ok {
				siteAuthSettingsProperties.GoogleClientSecret = utils.String(v.(string))
			}

			if v, ok := googleSetting["oauth_scopes"]; ok {
				input := v.([]interface{})

				oauthScopes := make([]string, 0)
				for _, param := range input {
					oauthScopes = append(oauthScopes, param.(string))
				}

				siteAuthSettingsProperties.GoogleOAuthScopes = &oauthScopes
			}
		}
	}

	if v, ok := setting["microsoft"]; ok {
		microsoftSettings := v.([]interface{})

		for _, setting := range microsoftSettings {
			microsoftSetting := setting.(map[string]interface{})

			if v, ok := microsoftSetting["client_id"]; ok {
				siteAuthSettingsProperties.MicrosoftAccountClientID = utils.String(v.(string))
			}

			if v, ok := microsoftSetting["client_secret"]; ok {
				siteAuthSettingsProperties.MicrosoftAccountClientSecret = utils.String(v.(string))
			}

			if v, ok := microsoftSetting["oauth_scopes"]; ok {
				input := v.([]interface{})

				oauthScopes := make([]string, 0)
				for _, param := range input {
					oauthScopes = append(oauthScopes, param.(string))
				}

				siteAuthSettingsProperties.MicrosoftAccountOAuthScopes = &oauthScopes
			}
		}
	}

	if v, ok := setting["twitter"]; ok {
		twitterSettings := v.([]interface{})

		for _, setting := range twitterSettings {
			twitterSetting := setting.(map[string]interface{})

			if v, ok := twitterSetting["consumer_key"]; ok {
				siteAuthSettingsProperties.TwitterConsumerKey = utils.String(v.(string))
			}

			if v, ok := twitterSetting["consumer_secret"]; ok {
				siteAuthSettingsProperties.TwitterConsumerSecret = utils.String(v.(string))
			}
		}
	}

	return siteAuthSettingsProperties
}

func flattenAdditionalLoginParams(input *[]string) map[string]interface{} {
	result := make(map[string]interface{})

	if input == nil {
		return result
	}

	for _, k := range *input {
		parts := strings.Split(k, "=")
		if len(parts) != 2 {
			continue // Params not following the format `key=value` is considered malformed and will be ignored.
		}
		key := parts[0]
		value := parts[1]

		result[key] = value
	}

	return result
}

func flattenAppServiceAuthSettings(input *web.SiteAuthSettingsProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.Enabled != nil {
		result["enabled"] = *input.Enabled
	}

	result["additional_login_params"] = flattenAdditionalLoginParams(input.AdditionalLoginParams)

	allowedExternalRedirectUrls := make([]string, 0)
	if s := input.AllowedExternalRedirectUrls; s != nil {
		allowedExternalRedirectUrls = *s
	}
	result["allowed_external_redirect_urls"] = allowedExternalRedirectUrls

	if input.DefaultProvider != "" {
		result["default_provider"] = input.DefaultProvider
	}

	if input.Issuer != nil {
		result["issuer"] = *input.Issuer
	}

	if input.RuntimeVersion != nil {
		result["runtime_version"] = *input.RuntimeVersion
	}

	if input.TokenRefreshExtensionHours != nil {
		result["token_refresh_extension_hours"] = *input.TokenRefreshExtensionHours
	}

	if input.TokenStoreEnabled != nil {
		result["token_store_enabled"] = *input.TokenStoreEnabled
	}

	if input.UnauthenticatedClientAction != "" {
		result["unauthenticated_client_action"] = input.UnauthenticatedClientAction
	}

	activeDirectorySettings := make([]interface{}, 0)

	if input.ClientID != nil {
		activeDirectorySetting := make(map[string]interface{})

		activeDirectorySetting["client_id"] = *input.ClientID

		if input.ClientSecret != nil {
			activeDirectorySetting["client_secret"] = *input.ClientSecret
		}

		if input.AllowedAudiences != nil {
			activeDirectorySetting["allowed_audiences"] = *input.AllowedAudiences
		}

		activeDirectorySettings = append(activeDirectorySettings, activeDirectorySetting)
	}

	result["active_directory"] = activeDirectorySettings

	facebookSettings := make([]interface{}, 0)

	if input.FacebookAppID != nil {
		facebookSetting := make(map[string]interface{})

		facebookSetting["app_id"] = *input.FacebookAppID

		if input.FacebookAppSecret != nil {
			facebookSetting["app_secret"] = *input.FacebookAppSecret
		}

		if input.FacebookOAuthScopes != nil {
			facebookSetting["oauth_scopes"] = *input.FacebookOAuthScopes
		}

		facebookSettings = append(facebookSettings, facebookSetting)
	}

	result["facebook"] = facebookSettings

	googleSettings := make([]interface{}, 0)

	if input.GoogleClientID != nil {
		googleSetting := make(map[string]interface{})

		googleSetting["client_id"] = *input.GoogleClientID

		if input.GoogleClientSecret != nil {
			googleSetting["client_secret"] = *input.GoogleClientSecret
		}

		if input.GoogleOAuthScopes != nil {
			googleSetting["oauth_scopes"] = *input.GoogleOAuthScopes
		}

		googleSettings = append(googleSettings, googleSetting)
	}

	result["google"] = googleSettings

	microsoftSettings := make([]interface{}, 0)

	if input.MicrosoftAccountClientID != nil {
		microsoftSetting := make(map[string]interface{})

		microsoftSetting["client_id"] = *input.MicrosoftAccountClientID

		if input.MicrosoftAccountClientSecret != nil {
			microsoftSetting["client_secret"] = *input.MicrosoftAccountClientSecret
		}

		if input.MicrosoftAccountOAuthScopes != nil {
			microsoftSetting["oauth_scopes"] = *input.MicrosoftAccountOAuthScopes
		}

		microsoftSettings = append(microsoftSettings, microsoftSetting)
	}

	result["microsoft"] = microsoftSettings

	twitterSettings := make([]interface{}, 0)

	if input.TwitterConsumerKey != nil {
		twitterSetting := make(map[string]interface{})

		twitterSetting["consumer_key"] = *input.TwitterConsumerKey

		if input.TwitterConsumerSecret != nil {
			twitterSetting["consumer_secret"] = *input.TwitterConsumerSecret
		}

		twitterSettings = append(twitterSettings, twitterSetting)
	}

	result["twitter"] = twitterSettings

	return append(results, result)
}

func flattenAppServiceLogs(input *web.SiteLogsConfigProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	appLogs := make([]interface{}, 0)
	if input.ApplicationLogs != nil {
		appLogsItem := make(map[string]interface{})

		if fileSystemInput := input.ApplicationLogs.FileSystem; fileSystemInput != nil {
			appLogsItem["file_system_level"] = string(fileSystemInput.Level)
		}

		blobStorage := make([]interface{}, 0)
		if blobStorageInput := input.ApplicationLogs.AzureBlobStorage; blobStorageInput != nil {
			blobStorageItem := make(map[string]interface{})

			blobStorageItem["level"] = string(blobStorageInput.Level)

			if blobStorageInput.SasURL != nil {
				blobStorageItem["sas_url"] = *blobStorageInput.SasURL
			}

			if blobStorageInput.RetentionInDays != nil {
				blobStorageItem["retention_in_days"] = *blobStorageInput.RetentionInDays
			}

			// The API returns a non nil application logs object when other logs are specified so we'll check that this structure is empty before adding it to the statefile.
			if blobStorageInput.SasURL != nil && *blobStorageInput.SasURL != "" {
				blobStorage = append(blobStorage, blobStorageItem)
			}
		}

		appLogsItem["azure_blob_storage"] = blobStorage
		appLogs = append(appLogs, appLogsItem)
	}
	result["application_logs"] = appLogs

	httpLogs := make([]interface{}, 0)
	if input.HTTPLogs != nil {
		httpLogsItem := make(map[string]interface{})

		fileSystem := make([]interface{}, 0)
		if fileSystemInput := input.HTTPLogs.FileSystem; fileSystemInput != nil {
			fileSystemItem := make(map[string]interface{})

			if fileSystemInput.RetentionInDays != nil {
				fileSystemItem["retention_in_days"] = *fileSystemInput.RetentionInDays
			}

			if fileSystemInput.RetentionInMb != nil {
				fileSystemItem["retention_in_mb"] = *fileSystemInput.RetentionInMb
			}

			// The API returns a non nil filesystem logs object when other logs are specified so we'll check that this is disabled before adding it to the statefile.
			if fileSystemInput.Enabled != nil && *fileSystemInput.Enabled {
				fileSystem = append(fileSystem, fileSystemItem)
			}
		}

		blobStorage := make([]interface{}, 0)
		if blobStorageInput := input.HTTPLogs.AzureBlobStorage; blobStorageInput != nil {
			blobStorageItem := make(map[string]interface{})

			if blobStorageInput.SasURL != nil {
				blobStorageItem["sas_url"] = *blobStorageInput.SasURL
			}

			if blobStorageInput.RetentionInDays != nil {
				blobStorageItem["retention_in_days"] = *blobStorageInput.RetentionInDays
			}

			// The API returns a non nil blob logs object when other logs are specified so we'll check that this is disabled before adding it to the statefile.
			if blobStorageInput.Enabled != nil && *blobStorageInput.Enabled {
				blobStorage = append(blobStorage, blobStorageItem)
			}
		}

		httpLogsItem["file_system"] = fileSystem
		httpLogsItem["azure_blob_storage"] = blobStorage
		httpLogs = append(httpLogs, httpLogsItem)
	}
	result["http_logs"] = httpLogs

	if input.DetailedErrorMessages != nil && input.DetailedErrorMessages.Enabled != nil {
		result["detailed_error_messages_enabled"] = *input.DetailedErrorMessages.Enabled
	}
	if input.FailedRequestsTracing != nil && input.FailedRequestsTracing.Enabled != nil {
		result["failed_request_tracing_enabled"] = *input.FailedRequestsTracing.Enabled
	}

	return append(results, result)
}

func expandAppServiceLogs(input interface{}) web.SiteLogsConfigProperties {
	configs := input.([]interface{})
	logs := web.SiteLogsConfigProperties{}

	if len(configs) == 0 || configs[0] == nil {
		return logs
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["application_logs"]; ok {
		appLogsConfigs := v.([]interface{})

		for _, config := range appLogsConfigs {
			logs.ApplicationLogs = &web.ApplicationLogsConfig{}

			if config == nil {
				continue
			}
			appLogsConfig := config.(map[string]interface{})

			if v, ok := appLogsConfig["file_system_level"]; ok {
				logs.ApplicationLogs.FileSystem = &web.FileSystemApplicationLogsConfig{
					Level: web.LogLevel(v.(string)),
				}
			}

			if v, ok := appLogsConfig["azure_blob_storage"]; ok {
				storageConfigs := v.([]interface{})

				for _, config := range storageConfigs {
					storageConfig := config.(map[string]interface{})

					logs.ApplicationLogs.AzureBlobStorage = &web.AzureBlobStorageApplicationLogsConfig{
						Level:           web.LogLevel(storageConfig["level"].(string)),
						SasURL:          utils.String(storageConfig["sas_url"].(string)),
						RetentionInDays: utils.Int32(int32(storageConfig["retention_in_days"].(int))),
					}
				}
			}
		}
	}

	if v, ok := config["http_logs"]; ok {
		httpLogsConfigs := v.([]interface{})

		for _, config := range httpLogsConfigs {
			logs.HTTPLogs = &web.HTTPLogsConfig{}

			if config == nil {
				continue
			}
			httpLogsConfig := config.(map[string]interface{})

			if v, ok := httpLogsConfig["file_system"]; ok {
				fileSystemConfigs := v.([]interface{})

				for _, config := range fileSystemConfigs {
					fileSystemConfig := config.(map[string]interface{})

					logs.HTTPLogs.FileSystem = &web.FileSystemHTTPLogsConfig{
						RetentionInMb:   utils.Int32(int32(fileSystemConfig["retention_in_mb"].(int))),
						RetentionInDays: utils.Int32(int32(fileSystemConfig["retention_in_days"].(int))),
						Enabled:         utils.Bool(true),
					}
				}
			}

			if v, ok := httpLogsConfig["azure_blob_storage"]; ok {
				storageConfigs := v.([]interface{})

				for _, config := range storageConfigs {
					storageConfig := config.(map[string]interface{})

					logs.HTTPLogs.AzureBlobStorage = &web.AzureBlobStorageHTTPLogsConfig{
						SasURL:          utils.String(storageConfig["sas_url"].(string)),
						RetentionInDays: utils.Int32(int32(storageConfig["retention_in_days"].(int))),
						Enabled:         utils.Bool(true),
					}
				}
			}
		}
	}

	if v, ok := config["detailed_error_messages_enabled"]; ok {
		logs.DetailedErrorMessages = &web.EnabledConfig{
			Enabled: utils.Bool(v.(bool)),
		}
	}

	if v, ok := config["failed_request_tracing_enabled"]; ok {
		logs.FailedRequestsTracing = &web.EnabledConfig{
			Enabled: utils.Bool(v.(bool)),
		}
	}

	return logs
}

func expandAppServiceIdentity(input []interface{}) *web.ManagedServiceIdentity {
	if len(input) == 0 {
		return nil
	}
	identity := input[0].(map[string]interface{})
	identityType := web.ManagedServiceIdentityType(identity["type"].(string))

	identityIds := make(map[string]*web.ManagedServiceIdentityUserAssignedIdentitiesValue)
	for _, id := range identity["identity_ids"].([]interface{}) {
		identityIds[id.(string)] = &web.ManagedServiceIdentityUserAssignedIdentitiesValue{}
	}

	managedServiceIdentity := web.ManagedServiceIdentity{
		Type: identityType,
	}

	if managedServiceIdentity.Type == web.ManagedServiceIdentityTypeUserAssigned || managedServiceIdentity.Type == web.ManagedServiceIdentityTypeSystemAssignedUserAssigned {
		managedServiceIdentity.UserAssignedIdentities = identityIds
	}

	return &managedServiceIdentity
}

func flattenAppServiceIdentity(identity *web.ManagedServiceIdentity) ([]interface{}, error) {
	if identity == nil {
		return make([]interface{}, 0), nil
	}

	principalId := ""
	if identity.PrincipalID != nil {
		principalId = *identity.PrincipalID
	}

	tenantId := ""
	if identity.TenantID != nil {
		tenantId = *identity.TenantID
	}

	identityIds := make([]string, 0)
	if identity.UserAssignedIdentities != nil {
		for key := range identity.UserAssignedIdentities {
			parsedId, err := parse.UserAssignedIdentityID(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
	}

	return []interface{}{
		map[string]interface{}{
			"identity_ids": identityIds,
			"principal_id": principalId,
			"tenant_id":    tenantId,
			"type":         string(identity.Type),
		},
	}, nil
}

func expandAppServiceSiteConfig(input interface{}) (*web.SiteConfig, error) {
	configs := input.([]interface{})
	siteConfig := &web.SiteConfig{}

	if len(configs) == 0 {
		return siteConfig, nil
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["always_on"]; ok {
		siteConfig.AlwaysOn = utils.Bool(v.(bool))
	}

	if v, ok := config["app_command_line"]; ok {
		siteConfig.AppCommandLine = utils.String(v.(string))
	}

	if v, ok := config["default_documents"]; ok {
		input := v.([]interface{})

		documents := make([]string, 0)
		for _, document := range input {
			documents = append(documents, document.(string))
		}

		siteConfig.DefaultDocuments = &documents
	}

	if v, ok := config["dotnet_framework_version"]; ok {
		siteConfig.NetFrameworkVersion = utils.String(v.(string))
	}

	if v, ok := config["java_version"]; ok {
		siteConfig.JavaVersion = utils.String(v.(string))
	}

	if v, ok := config["java_container"]; ok {
		siteConfig.JavaContainer = utils.String(v.(string))
	}

	if v, ok := config["java_container_version"]; ok {
		siteConfig.JavaContainerVersion = utils.String(v.(string))
	}

	if v, ok := config["linux_fx_version"]; ok {
		siteConfig.LinuxFxVersion = utils.String(v.(string))
	}

	if v, ok := config["windows_fx_version"]; ok {
		siteConfig.WindowsFxVersion = utils.String(v.(string))
	}

	if v, ok := config["http2_enabled"]; ok {
		siteConfig.HTTP20Enabled = utils.Bool(v.(bool))
	}

	if v, ok := config["ip_restriction"]; ok {
		ipSecurityRestrictions := v.(interface{})
		restrictions, err := expandAppServiceIpRestriction(ipSecurityRestrictions)
		if err != nil {
			return siteConfig, err
		}
		siteConfig.IPSecurityRestrictions = &restrictions
	}

	if v, ok := config["scm_use_main_ip_restriction"]; ok {
		siteConfig.ScmIPSecurityRestrictionsUseMain = utils.Bool(v.(bool))
	}

	if v, ok := config["scm_ip_restriction"]; ok {
		scmIPSecurityRestrictions := v.([]interface{})
		scmRestrictions, err := expandAppServiceIpRestriction(scmIPSecurityRestrictions)
		if err != nil {
			return siteConfig, err
		}
		siteConfig.ScmIPSecurityRestrictions = &scmRestrictions
	}

	if v, ok := config["local_mysql_enabled"]; ok {
		siteConfig.LocalMySQLEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["managed_pipeline_mode"]; ok {
		siteConfig.ManagedPipelineMode = web.ManagedPipelineMode(v.(string))
	}

	if v, ok := config["php_version"]; ok {
		siteConfig.PhpVersion = utils.String(v.(string))
	}

	if v, ok := config["python_version"]; ok {
		siteConfig.PythonVersion = utils.String(v.(string))
	}

	if v, ok := config["remote_debugging_enabled"]; ok {
		siteConfig.RemoteDebuggingEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["remote_debugging_version"]; ok {
		siteConfig.RemoteDebuggingVersion = utils.String(v.(string))
	}

	if v, ok := config["use_32_bit_worker_process"]; ok {
		siteConfig.Use32BitWorkerProcess = utils.Bool(v.(bool))
	}

	if v, ok := config["websockets_enabled"]; ok {
		siteConfig.WebSocketsEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["scm_type"]; ok {
		siteConfig.ScmType = web.ScmType(v.(string))
	}

	if v, ok := config["ftps_state"]; ok {
		siteConfig.FtpsState = web.FtpsState(v.(string))
	}

	if v, ok := config["health_check_path"]; ok {
		siteConfig.HealthCheckPath = utils.String(v.(string))
	}

	if v, ok := config["number_of_workers"]; ok && v.(int) != 0 {
		siteConfig.NumberOfWorkers = utils.Int32(int32(v.(int)))
	}

	if v, ok := config["min_tls_version"]; ok {
		siteConfig.MinTLSVersion = web.SupportedTLSVersions(v.(string))
	}

	if v, ok := config["cors"]; ok {
		corsSettings := v.(interface{})
		expand := ExpandWebCorsSettings(corsSettings)
		siteConfig.Cors = &expand
	}

	if v, ok := config["auto_swap_slot_name"]; ok {
		siteConfig.AutoSwapSlotName = utils.String(v.(string))
	}

	return siteConfig, nil
}

func flattenAppServiceSiteConfig(input *web.SiteConfig) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	if input.AlwaysOn != nil {
		result["always_on"] = *input.AlwaysOn
	}

	if input.AppCommandLine != nil {
		result["app_command_line"] = *input.AppCommandLine
	}

	documents := make([]string, 0)
	if s := input.DefaultDocuments; s != nil {
		documents = *s
	}
	result["default_documents"] = documents

	if input.NetFrameworkVersion != nil {
		result["dotnet_framework_version"] = *input.NetFrameworkVersion
	}

	if input.JavaVersion != nil {
		result["java_version"] = *input.JavaVersion
	}

	if input.JavaContainer != nil {
		result["java_container"] = *input.JavaContainer
	}

	if input.JavaContainerVersion != nil {
		result["java_container_version"] = *input.JavaContainerVersion
	}

	if input.LocalMySQLEnabled != nil {
		result["local_mysql_enabled"] = *input.LocalMySQLEnabled
	}

	if input.HTTP20Enabled != nil {
		result["http2_enabled"] = *input.HTTP20Enabled
	}

	result["ip_restriction"] = flattenAppServiceIpRestriction(input.IPSecurityRestrictions)

	if input.ScmIPSecurityRestrictionsUseMain != nil {
		result["scm_use_main_ip_restriction"] = *input.ScmIPSecurityRestrictionsUseMain
	}

	result["scm_ip_restriction"] = flattenAppServiceIpRestriction(input.ScmIPSecurityRestrictions)

	result["managed_pipeline_mode"] = string(input.ManagedPipelineMode)

	if input.PhpVersion != nil {
		result["php_version"] = *input.PhpVersion
	}

	if input.PythonVersion != nil {
		result["python_version"] = *input.PythonVersion
	}

	if input.RemoteDebuggingEnabled != nil {
		result["remote_debugging_enabled"] = *input.RemoteDebuggingEnabled
	}

	if input.RemoteDebuggingVersion != nil {
		result["remote_debugging_version"] = *input.RemoteDebuggingVersion
	}

	if input.Use32BitWorkerProcess != nil {
		result["use_32_bit_worker_process"] = *input.Use32BitWorkerProcess
	}

	if input.WebSocketsEnabled != nil {
		result["websockets_enabled"] = *input.WebSocketsEnabled
	}

	if input.LinuxFxVersion != nil {
		result["linux_fx_version"] = *input.LinuxFxVersion
	}

	if input.WindowsFxVersion != nil {
		result["windows_fx_version"] = *input.WindowsFxVersion
	}

	result["scm_type"] = string(input.ScmType)
	result["ftps_state"] = string(input.FtpsState)

	if input.HealthCheckPath != nil {
		result["health_check_path"] = *input.HealthCheckPath
	}

	if input.NumberOfWorkers != nil {
		result["number_of_workers"] = *input.NumberOfWorkers
	}

	result["min_tls_version"] = string(input.MinTLSVersion)

	result["cors"] = FlattenWebCorsSettings(input.Cors)

	if input.AutoSwapSlotName != nil {
		result["auto_swap_slot_name"] = *input.AutoSwapSlotName
	}

	return append(results, result)
}

func flattenAppServiceIpRestriction(input *[]web.IPSecurityRestriction) []interface{} {
	restrictions := make([]interface{}, 0)

	if input == nil {
		return restrictions
	}

	for _, v := range *input {
		restriction := make(map[string]interface{})
		if ip := v.IPAddress; ip != nil {
			if *ip == "Any" {
				continue
			} else {
				switch v.Tag {
				case web.ServiceTag:
					restriction["service_tag"] = *ip
				default:
					restriction["ip_address"] = *ip
				}
			}
		}

		subnetId := ""
		if subnetIdRaw := v.VnetSubnetResourceID; subnetIdRaw != nil {
			subnetId = *subnetIdRaw
		}
		restriction["virtual_network_subnet_id"] = subnetId

		name := ""
		if nameRaw := v.Name; nameRaw != nil {
			name = *nameRaw
		}
		restriction["name"] = name

		priority := 0
		if priorityRaw := v.Priority; priorityRaw != nil {
			priority = int(*priorityRaw)
		}
		restriction["priority"] = priority

		action := ""
		if actionRaw := v.Action; actionRaw != nil {
			action = *actionRaw
		}
		restriction["action"] = action

		if headers := v.Headers; headers != nil {
			restriction["headers"] = flattenHeaders(headers)
		}

		restrictions = append(restrictions, restriction)
	}

	return restrictions
}

func expandAppServiceStorageAccounts(input []interface{}) map[string]*web.AzureStorageInfoValue {
	output := make(map[string]*web.AzureStorageInfoValue, len(input))

	for _, v := range input {
		vals := v.(map[string]interface{})

		saName := vals["name"].(string)
		saType := vals["type"].(string)
		saAccountName := vals["account_name"].(string)
		saShareName := vals["share_name"].(string)
		saAccessKey := vals["access_key"].(string)
		saMountPath := vals["mount_path"].(string)

		output[saName] = &web.AzureStorageInfoValue{
			Type:        web.AzureStorageType(saType),
			AccountName: utils.String(saAccountName),
			ShareName:   utils.String(saShareName),
			AccessKey:   utils.String(saAccessKey),
			MountPath:   utils.String(saMountPath),
		}
	}

	return output
}

func flattenAppServiceStorageAccounts(input map[string]*web.AzureStorageInfoValue) []interface{} {
	results := make([]interface{}, 0)

	for k, v := range input {
		result := make(map[string]interface{})
		result["name"] = k
		result["type"] = string(v.Type)
		if v.AccountName != nil {
			result["account_name"] = *v.AccountName
		}
		if v.ShareName != nil {
			result["share_name"] = *v.ShareName
		}
		if v.AccessKey != nil {
			result["access_key"] = *v.AccessKey
		}
		if v.MountPath != nil {
			result["mount_path"] = *v.MountPath
		}
		results = append(results, result)
	}

	return results
}

func expandAppServiceIpRestriction(input interface{}) ([]web.IPSecurityRestriction, error) {
	restrictions := make([]web.IPSecurityRestriction, 0)

	for _, r := range input.([]interface{}) {
		if r == nil {
			continue
		}

		restriction := r.(map[string]interface{})

		ipAddress := restriction["ip_address"].(string)
		vNetSubnetID := ""

		if subnetID, ok := restriction["virtual_network_subnet_id"]; ok && subnetID != "" {
			vNetSubnetID = subnetID.(string)
		}

		serviceTag := restriction["service_tag"].(string)

		name := restriction["name"].(string)
		priority := restriction["priority"].(int)
		action := restriction["action"].(string)

		if vNetSubnetID != "" && ipAddress != "" && serviceTag != "" {
			return nil, fmt.Errorf("only one of `ip_address`, `service_tag` or `virtual_network_subnet_id` can be set for an IP restriction")
		}

		if vNetSubnetID == "" && ipAddress == "" && serviceTag == "" {
			return nil, fmt.Errorf("one of `ip_address`, `service_tag` or `virtual_network_subnet_id` must be set for an IP restriction")
		}

		ipSecurityRestriction := web.IPSecurityRestriction{}
		if ipAddress == "Any" {
			continue
		}

		if ipAddress != "" {
			ipSecurityRestriction.IPAddress = &ipAddress
		}

		if serviceTag != "" {
			ipSecurityRestriction.IPAddress = &serviceTag
			ipSecurityRestriction.Tag = web.ServiceTag
		}

		if vNetSubnetID != "" {
			ipSecurityRestriction.VnetSubnetResourceID = &vNetSubnetID
		}

		if name != "" {
			ipSecurityRestriction.Name = &name
		}

		if priority != 0 {
			ipSecurityRestriction.Priority = utils.Int32(int32(priority))
		}

		if action != "" {
			ipSecurityRestriction.Action = &action
		}
		if headers, ok := restriction["headers"]; ok {
			ipSecurityRestriction.Headers = expandHeaders(headers.([]interface{}))
		}

		restrictions = append(restrictions, ipSecurityRestriction)
	}

	return restrictions, nil
}

func flattenHeaders(input map[string][]string) []interface{} {
	output := make([]interface{}, 0)
	headers := make(map[string]interface{})
	if input == nil {
		return output
	}

	if forwardedHost, ok := input["x-forwarded-host"]; ok && len(forwardedHost) > 0 {
		headers["x_forwarded_host"] = forwardedHost
	}
	if forwardedFor, ok := input["x-forwarded-for"]; ok && len(forwardedFor) > 0 {
		headers["x_forwarded_for"] = forwardedFor
	}
	if fdids, ok := input["x-azure-fdid"]; ok && len(fdids) > 0 {
		headers["x_azure_fdid"] = fdids
	}
	if healthProbe, ok := input["x-fd-healthprobe"]; ok && len(healthProbe) > 0 {
		headers["x_fd_health_probe"] = healthProbe
	}

	return append(output, headers)
}

func expandHeaders(input interface{}) map[string][]string {
	output := make(map[string][]string)

	for _, r := range input.([]interface{}) {
		if r == nil {
			continue
		}

		val := r.(map[string]interface{})
		if raw := val["x_forwarded_host"].(*pluginsdk.Set).List(); len(raw) > 0 {
			output["x-forwarded-host"] = *utils.ExpandStringSlice(raw)
		}
		if raw := val["x_forwarded_for"].(*pluginsdk.Set).List(); len(raw) > 0 {
			output["x-forwarded-for"] = *utils.ExpandStringSlice(raw)
		}
		if raw := val["x_azure_fdid"].(*pluginsdk.Set).List(); len(raw) > 0 {
			output["x-azure-fdid"] = *utils.ExpandStringSlice(raw)
		}
		if raw := val["x_fd_health_probe"].(*pluginsdk.Set).List(); len(raw) > 0 {
			output["x-fd-healthprobe"] = *utils.ExpandStringSlice(raw)
		}
	}

	return output
}
