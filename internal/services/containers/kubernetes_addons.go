package containers

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-08-01/containerservice"
	"github.com/Azure/go-autorest/autorest/azure"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	laparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	logAnalyticsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	applicationGatewayValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	subnetValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	// note: the casing on these keys is important
	aciConnectorKey                 = "aciConnectorLinux"
	azurePolicyKey                  = "azurepolicy"
	kubernetesDashboardKey          = "kubeDashboard"
	httpApplicationRoutingKey       = "httpApplicationRouting"
	omsAgentKey                     = "omsagent"
	ingressApplicationGatewayKey    = "ingressApplicationGateway"
	openServiceMeshKey              = "openServiceMesh"
	azureKeyvaultSecretsProviderKey = "azureKeyvaultSecretsProvider"
)

// The AKS API hard-codes which add-ons are supported in which environment
// as such unfortunately we can't just send "disabled" - we need to strip
// the unsupported addons from the HTTP response. As such this defines
// the list of unsupported addons in the defined region - e.g. by being
// omitted from this list an addon/environment combination will be supported
var unsupportedAddonsForEnvironment = map[string][]string{
	azure.ChinaCloud.Name: {
		aciConnectorKey,                 // https://github.com/hashicorp/terraform-provider-azurerm/issues/5510
		httpApplicationRoutingKey,       // https://github.com/hashicorp/terraform-provider-azurerm/issues/5960
		kubernetesDashboardKey,          // https://github.com/hashicorp/terraform-provider-azurerm/issues/7487
		openServiceMeshKey,              // Preview features are not supported in Azure China
		azureKeyvaultSecretsProviderKey, // Preview features are not supported in Azure China
	},
	azure.USGovernmentCloud.Name: {
		httpApplicationRoutingKey,       // https://github.com/hashicorp/terraform-provider-azurerm/issues/5960
		kubernetesDashboardKey,          // https://github.com/hashicorp/terraform-provider-azurerm/issues/7136
		openServiceMeshKey,              // Preview features are not supported in Azure Government
		azureKeyvaultSecretsProviderKey, // Preview features are not supported in Azure China
	},
}

// TODO 3.0 - Remove this schema as it's deprecated
func schemaKubernetesAddOnProfiles() *pluginsdk.Schema {
	//lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		ConflictsWith: []string{
			"aci_connector_linux",
			"azure_policy_enabled",
			"http_application_routing_enabled",
			"oms_agent",
			"ingress_application_gateway",
			"open_service_mesh_enabled",
			"key_vault_secrets_provider",
		},
		Deprecated: "`addon_profile` block has been deprecated and will be removed in version 3.0 of the AzureRM Provider. All properties within the block will move to the top level.",
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"aci_connector_linux": {
					Type:       pluginsdk.TypeList,
					MaxItems:   1,
					Optional:   true,
					Deprecated: "`addon_profile.0.aci_connector_linux` block has been deprecated in favour of the `aci_connector_linux` block and will be removed in version 3.0 of the AzureRM Provider.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:       pluginsdk.TypeBool,
								Required:   true,
								Deprecated: "`addon_profile.0.aci_connector_linux.0.enabled` has been deprecated and will be removed in version 3.0 of the AzureRM Provider.",
							},

							"subnet_name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Deprecated:   "`addon_profile.0.aci_connector_linux.0.subnet_name` has been deprecated in favour of `aci_connector_linux.0.subnet_name` and will be removed in version 3.0 of the AzureRM Provider.",
							},
						},
					},
				},

				"azure_policy": {
					Type:       pluginsdk.TypeList,
					MaxItems:   1,
					Optional:   true,
					Deprecated: "`addon_profile.0.azure_policy` has been deprecated in favour of `azure_policy_enabled` and will be removed in version 3.0 of the AzureRM Provider.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:       pluginsdk.TypeBool,
								Required:   true,
								Deprecated: "`addon_profile.0.azure_policy.0.enabled` has been deprecated and will be removed in version 3.0 of the AzureRM Provider.",
							},
						},
					},
				},

				"kube_dashboard": {
					Type:       pluginsdk.TypeList,
					MaxItems:   1,
					Optional:   true,
					Deprecated: "`kube_dashboard` has been deprecated since it is no longer supported by Kubernetes versions 1.19 or above, this property will be removed in version 3.0 of the AzureRM Provider.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:     pluginsdk.TypeBool,
								Required: true,
							},
						},
					},
				},

				"http_application_routing": {
					Type:       pluginsdk.TypeList,
					MaxItems:   1,
					Optional:   true,
					Deprecated: "`addon_profile.0.http_application_routing` block has been deprecated in favour of the `http_application_routing_enabled` property and will be removed in version 3.0 of the AzureRM Provider.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:       pluginsdk.TypeBool,
								Required:   true,
								Deprecated: "`addon_profile.0.http_application_routing.0.enabled` has been deprecated and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"http_application_routing_zone_name": {
								Type:       pluginsdk.TypeString,
								Computed:   true,
								Deprecated: "`addon_profile.0.http_application_routing.0.http_application_routing_zone_name` has been deprecated in favour of `http_application_routing_zone_name` and will be removed in version 3.0 of the AzureRM Provider.",
							},
						},
					},
				},

				"oms_agent": {
					Type:       pluginsdk.TypeList,
					MaxItems:   1,
					Optional:   true,
					Deprecated: "`addon_profile.0.oms_agent` block has been deprecated in favour of the `oms_agent` block and will be removed in version 3.0 of the AzureRM Provider.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:       pluginsdk.TypeBool,
								Required:   true,
								Deprecated: "`addon_profile.0.oms_agent.0.enabled` has been deprecated and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"log_analytics_workspace_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: logAnalyticsValidate.LogAnalyticsWorkspaceID,
								Deprecated:   "`addon_profile.0.oms_agent.0.log_analytics_workspace_id` has been deprecated in favour of `oms_agent.0.log_analytics_workspace_id` and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"oms_agent_identity": {
								Type:       pluginsdk.TypeList,
								Computed:   true,
								Deprecated: "`addon_profile.0.oms_agent.0.oms_agent_identity` has been deprecated in favour of `oms_agent.0.oms_agent_identity` and will be removed in version 3.0 of the AzureRM Provider.",
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"client_id": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"object_id": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"user_assigned_identity_id": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},

				"ingress_application_gateway": {
					Type:       pluginsdk.TypeList,
					MaxItems:   1,
					Optional:   true,
					Deprecated: "`addon_profile.0.ingress_application_gateway` block has been deprecated in favour of the `ingress_application_gateway` block and will be removed in version 3.0 of the AzureRM Provider.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:       pluginsdk.TypeBool,
								Required:   true,
								Deprecated: "`addon_profile.0.ingress_application_gateway.0.enabled` has been deprecated and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"gateway_id": {
								Type:          pluginsdk.TypeString,
								Optional:      true,
								ConflictsWith: []string{"addon_profile.0.ingress_application_gateway.0.subnet_cidr", "addon_profile.0.ingress_application_gateway.0.subnet_id"},
								ValidateFunc:  applicationGatewayValidate.ApplicationGatewayID,
								Deprecated:    "`addon_profile.0.ingress_application_gateway.0.gateway_id` has been deprecated in favour of `ingress_application_gateway.0.gateway_id` and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"gateway_name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Deprecated:   "`addon_profile.0.ingress_application_gateway.0.gateway_name` has been deprecated in favour of `ingress_application_gateway.0.gateway_name` and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"subnet_cidr": {
								Type:          pluginsdk.TypeString,
								Optional:      true,
								ConflictsWith: []string{"addon_profile.0.ingress_application_gateway.0.gateway_id", "addon_profile.0.ingress_application_gateway.0.subnet_id"},
								ValidateFunc:  commonValidate.CIDR,
								Deprecated:    "`addon_profile.0.ingress_application_gateway.0.subnet_cidr` has been deprecated in favour of `ingress_application_gateway.0.subnet_cidr` and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"subnet_id": {
								Type:          pluginsdk.TypeString,
								Optional:      true,
								ConflictsWith: []string{"addon_profile.0.ingress_application_gateway.0.gateway_id", "addon_profile.0.ingress_application_gateway.0.subnet_cidr"},
								ValidateFunc:  subnetValidate.SubnetID,
								Deprecated:    "`addon_profile.0.ingress_application_gateway.0.subnet_id` has been deprecated in favour of `ingress_application_gateway.0.subnet_id` and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"effective_gateway_id": {
								Type:       pluginsdk.TypeString,
								Computed:   true,
								Deprecated: "`addon_profile.0.ingress_application_gateway.0.effective_gateway_id` has been deprecated in favour of `ingress_application_gateway.0.effective_gateway_id` and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"ingress_application_gateway_identity": {
								Type:       pluginsdk.TypeList,
								Computed:   true,
								Deprecated: "`addon_profile.0.ingress_application_gateway.0.ingress_application_gateway_identity` has been deprecated in favour of `ingress_application_gateway.0.ingress_application_gateway_identity` and will be removed in version 3.0 of the AzureRM Provider.",
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"client_id": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"object_id": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"user_assigned_identity_id": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},

				"open_service_mesh": {
					Type:       pluginsdk.TypeList,
					MaxItems:   1,
					Optional:   true,
					Deprecated: "`addon_profile.0.open_service_mesh` has been deprecated in favour of `open_service_mesh_enabled` and will be removed in version 3.0 of the AzureRM Provider.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:       pluginsdk.TypeBool,
								Required:   true,
								Deprecated: "`addon_profile.0.open_service_mesh.0.enabled` has been deprecated and will be removed in version 3.0 of the AzureRM Provider.",
							},
						},
					},
				},

				"azure_keyvault_secrets_provider": {
					Type:       pluginsdk.TypeList,
					MaxItems:   1,
					Optional:   true,
					Deprecated: "`addon_profile.0.azure_keyvault_secrets_provider` block has been deprecated in favour of the `key_vault_secrets_provider` block and will be removed in version 3.0 of the AzureRM Provider.",
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:       pluginsdk.TypeBool,
								Required:   true,
								Deprecated: "`addon_profile.0.azure_keyvault_secrets_provider.0.enabled` has been deprecated and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"secret_rotation_enabled": {
								Type:       pluginsdk.TypeBool,
								Default:    false,
								Optional:   true,
								Deprecated: "`addon_profile.0.azure_keyvault_secrets_provider.0.secret_rotation_enabled` has been deprecated in favour of `key_vault_secrets_provider.0.secret_rotation_enabled` and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"secret_rotation_interval": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Default:      "2m",
								ValidateFunc: containerValidate.Duration,
								Deprecated:   "`addon_profile.0.azure_keyvault_secrets_provider.0.secret_rotation_interval` has been deprecated in favour of `key_vault_secrets_provider.0.secret_rotation_interval` and will be removed in version 3.0 of the AzureRM Provider.",
							},
							"secret_identity": {
								Type:       pluginsdk.TypeList,
								Computed:   true,
								Deprecated: "`addon_profile.0.azure_keyvault_secrets_provider.0.secret_identity` has been deprecated in favour of `key_vault_secrets_provider.0.secret_identity` and will be removed in version 3.0 of the AzureRM Provider.",
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"client_id": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"object_id": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
										"user_assigned_identity_id": {
											Type:     pluginsdk.TypeString,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func schemaKubernetesAddOns() map[string]*pluginsdk.Schema {
	out := map[string]*pluginsdk.Schema{
		"aci_connector_linux": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: !features.ThreePointOhBeta(),
			ConflictsWith: func() []string {
				if !features.ThreePointOhBeta() {
					return []string{"addon_profile"}
				}
				return []string{}
			}(),
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"azure_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: !features.ThreePointOhBeta(),
			ConflictsWith: func() []string {
				if !features.ThreePointOhBeta() {
					return []string{"addon_profile"}
				}
				return []string{}
			}(),
		},
		"http_application_routing_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: !features.ThreePointOhBeta(),
			ConflictsWith: func() []string {
				if !features.ThreePointOhBeta() {
					return []string{"addon_profile"}
				}
				return []string{}
			}(),
		},
		"http_application_routing_zone_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"oms_agent": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: !features.ThreePointOhBeta(),
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"log_analytics_workspace_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: logAnalyticsValidate.LogAnalyticsWorkspaceID,
					},
					"oms_agent_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"user_assigned_identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
			ConflictsWith: func() []string {
				if !features.ThreePointOhBeta() {
					return []string{"addon_profile"}
				}
				return []string{}
			}(),
		},
		"ingress_application_gateway": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: !features.ThreePointOhBeta(),
			ConflictsWith: func() []string {
				if !features.ThreePointOhBeta() {
					return []string{"addon_profile"}
				}
				return []string{}
			}(),
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"gateway_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						AtLeastOneOf: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						ValidateFunc: applicationGatewayValidate.ApplicationGatewayID,
					},
					"gateway_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"subnet_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_id",
						},
						AtLeastOneOf: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						ValidateFunc: commonValidate.CIDR,
					},
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
						},
						AtLeastOneOf: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						ValidateFunc: subnetValidate.SubnetID,
					},
					"effective_gateway_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"ingress_application_gateway_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"user_assigned_identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"open_service_mesh_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: !features.ThreePointOhBeta(),
			ConflictsWith: func() []string {
				if !features.ThreePointOhBeta() {
					return []string{"addon_profile"}
				}
				return []string{}
			}(),
		},
		"key_vault_secrets_provider": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: !features.ThreePointOhBeta(),
			ConflictsWith: func() []string {
				if !features.ThreePointOhBeta() {
					return []string{"addon_profile"}
				}
				return []string{}
			}(),
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"secret_rotation_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  false,
						Optional: true,
						AtLeastOneOf: []string{
							"key_vault_secrets_provider.0.secret_rotation_enabled",
							"key_vault_secrets_provider.0.secret_rotation_interval",
						},
					},
					"secret_rotation_interval": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "2m",
						AtLeastOneOf: []string{
							"key_vault_secrets_provider.0.secret_rotation_enabled",
							"key_vault_secrets_provider.0.secret_rotation_interval",
						},
						ValidateFunc: containerValidate.Duration,
					},
					"secret_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"user_assigned_identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}

	if !features.ThreePointOhBeta() {
		out["addon_profile"] = schemaKubernetesAddOnProfiles()
	}

	return out
}

// TODO 3.0 - Remove this function
func expandKubernetesAddOnProfiles(input []interface{}, env azure.Environment) (*map[string]*containerservice.ManagedClusterAddonProfile, error) {
	disabled := containerservice.ManagedClusterAddonProfile{
		Enabled: utils.Bool(false),
	}

	profiles := map[string]*containerservice.ManagedClusterAddonProfile{
		aciConnectorKey:                 &disabled,
		azurePolicyKey:                  &disabled,
		kubernetesDashboardKey:          &disabled,
		httpApplicationRoutingKey:       &disabled,
		omsAgentKey:                     &disabled,
		ingressApplicationGatewayKey:    &disabled,
		openServiceMeshKey:              &disabled,
		azureKeyvaultSecretsProviderKey: &disabled,
	}

	if len(input) == 0 || input[0] == nil {
		return filterUnsupportedKubernetesAddOns(profiles, env)
	}

	profile := input[0].(map[string]interface{})
	addonProfiles := map[string]*containerservice.ManagedClusterAddonProfile{}

	httpApplicationRouting := profile["http_application_routing"].([]interface{})
	if len(httpApplicationRouting) > 0 && httpApplicationRouting[0] != nil {
		value := httpApplicationRouting[0].(map[string]interface{})
		enabled := value["enabled"].(bool)
		addonProfiles[httpApplicationRoutingKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
		}
	}

	omsAgent := profile["oms_agent"].([]interface{})
	if len(omsAgent) > 0 && omsAgent[0] != nil {
		value := omsAgent[0].(map[string]interface{})
		config := make(map[string]*string)
		enabled := value["enabled"].(bool)

		if workspaceID, ok := value["log_analytics_workspace_id"]; ok && workspaceID != "" {
			lawid, err := laparse.LogAnalyticsWorkspaceID(workspaceID.(string))
			if err != nil {
				return nil, fmt.Errorf("parsing Log Analytics Workspace ID: %+v", err)
			}
			config["logAnalyticsWorkspaceResourceID"] = utils.String(lawid.ID())
		}

		addonProfiles[omsAgentKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  config,
		}
	}

	aciConnector := profile["aci_connector_linux"].([]interface{})
	if len(aciConnector) > 0 && aciConnector[0] != nil {
		value := aciConnector[0].(map[string]interface{})
		config := make(map[string]*string)
		enabled := value["enabled"].(bool)

		if subnetName, ok := value["subnet_name"]; ok && subnetName != "" {
			config["SubnetName"] = utils.String(subnetName.(string))
		}

		addonProfiles[aciConnectorKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  config,
		}
	}

	kubeDashboard := profile["kube_dashboard"].([]interface{})
	if len(kubeDashboard) > 0 && kubeDashboard[0] != nil {
		value := kubeDashboard[0].(map[string]interface{})
		enabled := value["enabled"].(bool)

		addonProfiles[kubernetesDashboardKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  nil,
		}
	}

	azurePolicy := profile["azure_policy"].([]interface{})
	if len(azurePolicy) > 0 && azurePolicy[0] != nil {
		value := azurePolicy[0].(map[string]interface{})
		enabled := value["enabled"].(bool)

		addonProfiles[azurePolicyKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config: map[string]*string{
				"version": utils.String("v2"),
			},
		}
	}

	ingressApplicationGateway := profile["ingress_application_gateway"].([]interface{})
	if len(ingressApplicationGateway) > 0 && ingressApplicationGateway[0] != nil {
		value := ingressApplicationGateway[0].(map[string]interface{})
		config := make(map[string]*string)
		enabled := value["enabled"].(bool)

		if gatewayId, ok := value["gateway_id"]; ok && gatewayId != "" {
			config["applicationGatewayId"] = utils.String(gatewayId.(string))
		}

		if gatewayName, ok := value["gateway_name"]; ok && gatewayName != "" {
			config["applicationGatewayName"] = utils.String(gatewayName.(string))
		}

		if subnetCIDR, ok := value["subnet_cidr"]; ok && subnetCIDR != "" {
			config["subnetCIDR"] = utils.String(subnetCIDR.(string))
		}

		if subnetId, ok := value["subnet_id"]; ok && subnetId != "" {
			config["subnetId"] = utils.String(subnetId.(string))
		}

		addonProfiles[ingressApplicationGatewayKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  config,
		}
	}

	openServiceMesh := profile["open_service_mesh"].([]interface{})
	if len(openServiceMesh) > 0 && openServiceMesh[0] != nil {
		value := openServiceMesh[0].(map[string]interface{})
		enabled := value["enabled"].(bool)

		addonProfiles[openServiceMeshKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  nil,
		}
	}

	azureKeyvaultSecretsProvider := profile["azure_keyvault_secrets_provider"].([]interface{})
	if len(azureKeyvaultSecretsProvider) > 0 && azureKeyvaultSecretsProvider[0] != nil {
		value := azureKeyvaultSecretsProvider[0].(map[string]interface{})
		config := make(map[string]*string)
		enabled := value["enabled"].(bool)

		enableSecretRotation := "false"
		if value["secret_rotation_enabled"].(bool) {
			enableSecretRotation = "true"
		}
		config["enableSecretRotation"] = utils.String(enableSecretRotation)
		config["rotationPollInterval"] = utils.String(value["secret_rotation_interval"].(string))

		addonProfiles[azureKeyvaultSecretsProviderKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  config,
		}
	}

	return filterUnsupportedKubernetesAddOns(addonProfiles, env)
}

func expandKubernetesAddOns(d *pluginsdk.ResourceData, input map[string]interface{}, env azure.Environment) (*map[string]*containerservice.ManagedClusterAddonProfile, error) {
	disabled := containerservice.ManagedClusterAddonProfile{
		Enabled: utils.Bool(false),
	}

	addonProfiles := map[string]*containerservice.ManagedClusterAddonProfile{}
	if d.HasChange("http_application_routing_enabled") {
		addonProfiles[httpApplicationRoutingKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(input["http_application_routing_enabled"].(bool)),
		}
	}

	omsAgent := input["oms_agent"].([]interface{})
	if len(omsAgent) > 0 && omsAgent[0] != nil {
		value := omsAgent[0].(map[string]interface{})
		config := make(map[string]*string)

		if workspaceID, ok := value["log_analytics_workspace_id"]; ok && workspaceID != "" {
			lawid, err := laparse.LogAnalyticsWorkspaceID(workspaceID.(string))
			if err != nil {
				return nil, fmt.Errorf("parsing Log Analytics Workspace ID: %+v", err)
			}
			config["logAnalyticsWorkspaceResourceID"] = utils.String(lawid.ID())
		}

		addonProfiles[omsAgentKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(true),
			Config:  config,
		}
	} else if len(omsAgent) == 0 && d.HasChange("oms_agent") {
		addonProfiles[omsAgentKey] = &disabled
	}

	aciConnector := input["aci_connector_linux"].([]interface{})
	if len(aciConnector) > 0 && aciConnector[0] != nil {
		value := aciConnector[0].(map[string]interface{})
		config := make(map[string]*string)

		if subnetName, ok := value["subnet_name"]; ok && subnetName != "" {
			config["SubnetName"] = utils.String(subnetName.(string))
		}

		addonProfiles[aciConnectorKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(true),
			Config:  config,
		}
	} else if len(aciConnector) == 0 && d.HasChange("aci_connector_linux") {
		addonProfiles[aciConnectorKey] = &disabled
	}

	if ok := d.HasChange("azure_policy_enabled"); ok {
		v := input["azure_policy_enabled"].(bool)
		props := &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(v),
			Config: map[string]*string{
				"version": utils.String("v2"),
			},
		}
		addonProfiles[azurePolicyKey] = props
	}

	ingressApplicationGateway := input["ingress_application_gateway"].([]interface{})
	if len(ingressApplicationGateway) > 0 && ingressApplicationGateway[0] != nil {
		value := ingressApplicationGateway[0].(map[string]interface{})
		config := make(map[string]*string)

		if gatewayId, ok := value["gateway_id"]; ok && gatewayId != "" {
			config["applicationGatewayId"] = utils.String(gatewayId.(string))
		}

		if gatewayName, ok := value["gateway_name"]; ok && gatewayName != "" {
			config["applicationGatewayName"] = utils.String(gatewayName.(string))
		}

		if subnetCIDR, ok := value["subnet_cidr"]; ok && subnetCIDR != "" {
			config["subnetCIDR"] = utils.String(subnetCIDR.(string))
		}

		if subnetId, ok := value["subnet_id"]; ok && subnetId != "" {
			config["subnetId"] = utils.String(subnetId.(string))
		}

		addonProfiles[ingressApplicationGatewayKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(true),
			Config:  config,
		}
	} else if len(ingressApplicationGateway) == 0 && d.HasChange("ingress_application_gateway") {
		addonProfiles[ingressApplicationGatewayKey] = &disabled
	}

	if ok := d.HasChange("open_service_mesh_enabled"); ok {
		addonProfiles[openServiceMeshKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(input["open_service_mesh_enabled"].(bool)),
			Config:  nil,
		}
	}

	azureKeyVaultSecretsProvider := input["key_vault_secrets_provider"].([]interface{})
	if len(azureKeyVaultSecretsProvider) > 0 && azureKeyVaultSecretsProvider[0] != nil {
		value := azureKeyVaultSecretsProvider[0].(map[string]interface{})
		config := make(map[string]*string)

		enableSecretRotation := fmt.Sprintf("%t", value["secret_rotation_enabled"].(bool))
		config["enableSecretRotation"] = utils.String(enableSecretRotation)
		config["rotationPollInterval"] = utils.String(value["secret_rotation_interval"].(string))

		addonProfiles[azureKeyvaultSecretsProviderKey] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(true),
			Config:  config,
		}
	} else if len(azureKeyVaultSecretsProvider) == 0 && d.HasChange("key_vault_secrets_provider") {
		addonProfiles[azureKeyvaultSecretsProviderKey] = &disabled
	}

	return filterUnsupportedKubernetesAddOns(addonProfiles, env)
}

func filterUnsupportedKubernetesAddOns(input map[string]*containerservice.ManagedClusterAddonProfile, env azure.Environment) (*map[string]*containerservice.ManagedClusterAddonProfile, error) {
	filter := func(input map[string]*containerservice.ManagedClusterAddonProfile, key string) (*map[string]*containerservice.ManagedClusterAddonProfile, error) {
		output := input
		if v, ok := output[key]; ok {
			if v.Enabled != nil && *v.Enabled {
				return nil, fmt.Errorf("The addon %q is not supported for a Kubernetes Cluster located in %q", key, env.Name)
			}

			// otherwise it's disabled by default, so just remove it
			delete(output, key)
		}

		return &output, nil
	}

	output := input
	if unsupportedAddons, ok := unsupportedAddonsForEnvironment[env.Name]; ok {
		for _, key := range unsupportedAddons {
			out, err := filter(output, key)
			if err != nil {
				return nil, err
			}

			output = *out
		}
	}
	return &output, nil
}

// TODO 3.0 - Remove this function
func flattenKubernetesAddOnProfiles(profile map[string]*containerservice.ManagedClusterAddonProfile) []interface{} {
	aciConnectors := make([]interface{}, 0)
	if aciConnector := kubernetesAddonProfileLocate(profile, aciConnectorKey); aciConnector != nil {
		enabled := false
		if enabledVal := aciConnector.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		subnetName := ""
		if v := aciConnector.Config["SubnetName"]; v != nil {
			subnetName = *v
		}

		aciConnectors = append(aciConnectors, map[string]interface{}{
			"enabled":     enabled,
			"subnet_name": subnetName,
		})
	}

	azurePolicies := make([]interface{}, 0)
	if azurePolicy := kubernetesAddonProfileLocate(profile, azurePolicyKey); azurePolicy != nil {
		enabled := false
		if enabledVal := azurePolicy.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		azurePolicies = append(azurePolicies, map[string]interface{}{
			"enabled": enabled,
		})
	}

	httpApplicationRoutes := make([]interface{}, 0)
	if httpApplicationRouting := kubernetesAddonProfileLocate(profile, httpApplicationRoutingKey); httpApplicationRouting != nil {
		enabled := false
		if enabledVal := httpApplicationRouting.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		zoneName := ""
		if v := kubernetesAddonProfilelocateInConfig(httpApplicationRouting.Config, "HTTPApplicationRoutingZoneName"); v != nil {
			zoneName = *v
		}

		httpApplicationRoutes = append(httpApplicationRoutes, map[string]interface{}{
			"enabled":                            enabled,
			"http_application_routing_zone_name": zoneName,
		})
	}

	kubeDashboards := make([]interface{}, 0)
	if kubeDashboard := kubernetesAddonProfileLocate(profile, kubernetesDashboardKey); kubeDashboard != nil {
		enabled := false
		if enabledVal := kubeDashboard.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		kubeDashboards = append(kubeDashboards, map[string]interface{}{
			"enabled": enabled,
		})
	}

	omsAgents := make([]interface{}, 0)
	if omsAgent := kubernetesAddonProfileLocate(profile, omsAgentKey); omsAgent != nil {
		enabled := false
		if enabledVal := omsAgent.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		workspaceID := ""
		if v := kubernetesAddonProfilelocateInConfig(omsAgent.Config, "logAnalyticsWorkspaceResourceID"); v != nil {
			if lawid, err := laparse.LogAnalyticsWorkspaceID(*v); err == nil {
				workspaceID = lawid.ID()
			}
		}

		omsagentIdentity := flattenKubernetesClusterAddOnIdentityProfile(omsAgent.Identity)

		omsAgents = append(omsAgents, map[string]interface{}{
			"enabled":                    enabled,
			"log_analytics_workspace_id": workspaceID,
			"oms_agent_identity":         omsagentIdentity,
		})
	}

	ingressApplicationGateways := make([]interface{}, 0)
	if ingressApplicationGateway := kubernetesAddonProfileLocate(profile, ingressApplicationGatewayKey); ingressApplicationGateway != nil {
		enabled := false
		if enabledVal := ingressApplicationGateway.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		gatewayId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayId"); v != nil {
			gatewayId = *v
		}

		gatewayName := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayName"); v != nil {
			gatewayName = *v
		}

		effectiveGatewayId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "effectiveApplicationGatewayId"); v != nil {
			effectiveGatewayId = *v
		}

		subnetCIDR := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetCIDR"); v != nil {
			subnetCIDR = *v
		}

		subnetId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetId"); v != nil {
			subnetId = *v
		}

		ingressApplicationGatewayIdentity := flattenKubernetesClusterAddOnIdentityProfile(ingressApplicationGateway.Identity)

		ingressApplicationGateways = append(ingressApplicationGateways, map[string]interface{}{
			"enabled":                              enabled,
			"gateway_id":                           gatewayId,
			"gateway_name":                         gatewayName,
			"effective_gateway_id":                 effectiveGatewayId,
			"subnet_cidr":                          subnetCIDR,
			"subnet_id":                            subnetId,
			"ingress_application_gateway_identity": ingressApplicationGatewayIdentity,
		})
	}

	openServiceMeshes := make([]interface{}, 0)
	if openServiceMesh := kubernetesAddonProfileLocate(profile, openServiceMeshKey); openServiceMesh != nil {
		enabled := false
		if enabledVal := openServiceMesh.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		openServiceMeshes = append(openServiceMeshes, map[string]interface{}{
			"enabled": enabled,
		})
	}

	azureKeyvaultSecretsProviders := make([]interface{}, 0)
	if azureKeyvaultSecretsProvider := kubernetesAddonProfileLocate(profile, azureKeyvaultSecretsProviderKey); azureKeyvaultSecretsProvider != nil {
		enabled := false
		if enabledVal := azureKeyvaultSecretsProvider.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}
		enableSecretRotation := false
		if v := kubernetesAddonProfilelocateInConfig(azureKeyvaultSecretsProvider.Config, "enableSecretRotation"); v != nil && *v != "false" {
			enableSecretRotation = true
		}
		rotationPollInterval := ""
		if v := kubernetesAddonProfilelocateInConfig(azureKeyvaultSecretsProvider.Config, "rotationPollInterval"); v != nil {
			rotationPollInterval = *v
		}

		azureKeyvaultSecretsProviderIdentity := flattenKubernetesClusterAddOnIdentityProfile(azureKeyvaultSecretsProvider.Identity)

		azureKeyvaultSecretsProviders = append(azureKeyvaultSecretsProviders, map[string]interface{}{
			"enabled":                  enabled,
			"secret_rotation_enabled":  enableSecretRotation,
			"secret_rotation_interval": rotationPollInterval,
			"secret_identity":          azureKeyvaultSecretsProviderIdentity,
		})
	}

	// this is a UX hack, since if the top level block isn't defined everything should be turned off
	if len(aciConnectors) == 0 && len(azurePolicies) == 0 && len(httpApplicationRoutes) == 0 && len(kubeDashboards) == 0 && len(omsAgents) == 0 && len(ingressApplicationGateways) == 0 && len(openServiceMeshes) == 0 && len(azureKeyvaultSecretsProviders) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"aci_connector_linux":             aciConnectors,
			"azure_policy":                    azurePolicies,
			"http_application_routing":        httpApplicationRoutes,
			"kube_dashboard":                  kubeDashboards,
			"oms_agent":                       omsAgents,
			"ingress_application_gateway":     ingressApplicationGateways,
			"open_service_mesh":               openServiceMeshes,
			"azure_keyvault_secrets_provider": azureKeyvaultSecretsProviders,
		},
	}
}

func flattenKubernetesAddOns(profile map[string]*containerservice.ManagedClusterAddonProfile) map[string]interface{} {
	aciConnectors := make([]interface{}, 0)
	if aciConnector := kubernetesAddonProfileLocate(profile, aciConnectorKey); aciConnector != nil {
		if enabled := aciConnector.Enabled; enabled != nil && *enabled {
			subnetName := ""
			if v := aciConnector.Config["SubnetName"]; v != nil {
				subnetName = *v
			}

			aciConnectors = append(aciConnectors, map[string]interface{}{
				"subnet_name": subnetName,
			})
		}

	}

	azurePolicyEnabled := false
	if azurePolicy := kubernetesAddonProfileLocate(profile, azurePolicyKey); azurePolicy != nil {
		if enabledVal := azurePolicy.Enabled; enabledVal != nil {
			azurePolicyEnabled = *enabledVal
		}
	}

	httpApplicationRoutingEnabled := false
	httpApplicationRoutingZone := ""
	if httpApplicationRouting := kubernetesAddonProfileLocate(profile, httpApplicationRoutingKey); httpApplicationRouting != nil {
		if enabledVal := httpApplicationRouting.Enabled; enabledVal != nil {
			httpApplicationRoutingEnabled = *enabledVal
		}

		if v := kubernetesAddonProfilelocateInConfig(httpApplicationRouting.Config, "HTTPApplicationRoutingZoneName"); v != nil {
			httpApplicationRoutingZone = *v
		}
	}

	omsAgents := make([]interface{}, 0)
	if omsAgent := kubernetesAddonProfileLocate(profile, omsAgentKey); omsAgent != nil {
		if enabled := omsAgent.Enabled; enabled != nil && *enabled {
			workspaceID := ""
			if v := kubernetesAddonProfilelocateInConfig(omsAgent.Config, "logAnalyticsWorkspaceResourceID"); v != nil {
				if lawid, err := laparse.LogAnalyticsWorkspaceID(*v); err == nil {
					workspaceID = lawid.ID()
				}
			}

			omsAgentIdentity := flattenKubernetesClusterAddOnIdentityProfile(omsAgent.Identity)

			omsAgents = append(omsAgents, map[string]interface{}{
				"log_analytics_workspace_id": workspaceID,
				"oms_agent_identity":         omsAgentIdentity,
			})
		}
	}

	ingressApplicationGateways := make([]interface{}, 0)
	if ingressApplicationGateway := kubernetesAddonProfileLocate(profile, ingressApplicationGatewayKey); ingressApplicationGateway != nil {
		if enabled := ingressApplicationGateway.Enabled; enabled != nil && *enabled {
			gatewayId := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayId"); v != nil {
				gatewayId = *v
			}

			gatewayName := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayName"); v != nil {
				gatewayName = *v
			}

			effectiveGatewayId := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "effectiveApplicationGatewayId"); v != nil {
				effectiveGatewayId = *v
			}

			subnetCIDR := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetCIDR"); v != nil {
				subnetCIDR = *v
			}

			subnetId := ""
			if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetId"); v != nil {
				subnetId = *v
			}

			ingressApplicationGatewayIdentity := flattenKubernetesClusterAddOnIdentityProfile(ingressApplicationGateway.Identity)

			ingressApplicationGateways = append(ingressApplicationGateways, map[string]interface{}{
				"gateway_id":                           gatewayId,
				"gateway_name":                         gatewayName,
				"effective_gateway_id":                 effectiveGatewayId,
				"subnet_cidr":                          subnetCIDR,
				"subnet_id":                            subnetId,
				"ingress_application_gateway_identity": ingressApplicationGatewayIdentity,
			})
		}
	}

	openServiceMeshEnabled := false
	if openServiceMesh := kubernetesAddonProfileLocate(profile, openServiceMeshKey); openServiceMesh != nil {
		if enabledVal := openServiceMesh.Enabled; enabledVal != nil {
			openServiceMeshEnabled = *enabledVal
		}
	}

	azureKeyVaultSecretsProviders := make([]interface{}, 0)
	if azureKeyVaultSecretsProvider := kubernetesAddonProfileLocate(profile, azureKeyvaultSecretsProviderKey); azureKeyVaultSecretsProvider != nil {
		if enabled := azureKeyVaultSecretsProvider.Enabled; enabled != nil && *enabled {
			enableSecretRotation := false
			if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "enableSecretRotation"); v != nil && *v != "false" {
				enableSecretRotation = true
			}

			rotationPollInterval := ""
			if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "rotationPollInterval"); v != nil {
				rotationPollInterval = *v
			}

			azureKeyvaultSecretsProviderIdentity := flattenKubernetesClusterAddOnIdentityProfile(azureKeyVaultSecretsProvider.Identity)

			azureKeyVaultSecretsProviders = append(azureKeyVaultSecretsProviders, map[string]interface{}{
				"secret_rotation_enabled":  enableSecretRotation,
				"secret_rotation_interval": rotationPollInterval,
				"secret_identity":          azureKeyvaultSecretsProviderIdentity,
			})
		}
	}

	return map[string]interface{}{
		"aci_connector_linux":                aciConnectors,
		"azure_policy_enabled":               azurePolicyEnabled,
		"http_application_routing_enabled":   httpApplicationRoutingEnabled,
		"http_application_routing_zone_name": httpApplicationRoutingZone,
		"oms_agent":                          omsAgents,
		"ingress_application_gateway":        ingressApplicationGateways,
		"open_service_mesh_enabled":          openServiceMeshEnabled,
		"key_vault_secrets_provider":         azureKeyVaultSecretsProviders,
	}
}

func flattenKubernetesClusterAddOnIdentityProfile(profile *containerservice.ManagedClusterAddonProfileIdentity) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	identity := make([]interface{}, 0)
	clientID := ""
	if clientid := profile.ClientID; clientid != nil {
		clientID = *clientid
	}

	objectID := ""
	if objectid := profile.ObjectID; objectid != nil {
		objectID = *objectid
	}

	userAssignedIdentityID := ""
	if resourceid := profile.ResourceID; resourceid != nil {
		userAssignedIdentityID = *resourceid
	}

	identity = append(identity, map[string]interface{}{
		"client_id":                 clientID,
		"object_id":                 objectID,
		"user_assigned_identity_id": userAssignedIdentityID,
	})

	return identity
}

func collectKubernetesAddons(d *pluginsdk.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"aci_connector_linux":              d.Get("aci_connector_linux").([]interface{}),
		"azure_policy_enabled":             d.Get("azure_policy_enabled").(bool),
		"http_application_routing_enabled": d.Get("http_application_routing_enabled").(bool),
		"oms_agent":                        d.Get("oms_agent").([]interface{}),
		"ingress_application_gateway":      d.Get("ingress_application_gateway").([]interface{}),
		"open_service_mesh_enabled":        d.Get("open_service_mesh_enabled").(bool),
		"key_vault_secrets_provider":       d.Get("key_vault_secrets_provider").([]interface{}),
	}
}

// when the Kubernetes Cluster is updated in the Portal - Azure updates the casing on the keys
// meaning what's submitted could be different to what's returned..
func kubernetesAddonProfileLocate(profile map[string]*containerservice.ManagedClusterAddonProfile, key string) *containerservice.ManagedClusterAddonProfile {
	for k, v := range profile {
		if strings.EqualFold(k, key) {
			return v
		}
	}

	return nil
}

// when the Kubernetes Cluster is updated in the Portal - Azure updates the casing on the keys
// meaning what's submitted could be different to what's returned..
// Related issue: https://github.com/Azure/azure-rest-api-specs/issues/10716
func kubernetesAddonProfilelocateInConfig(config map[string]*string, key string) *string {
	for k, v := range config {
		if strings.EqualFold(k, key) {
			return v
		}
	}

	return nil
}
