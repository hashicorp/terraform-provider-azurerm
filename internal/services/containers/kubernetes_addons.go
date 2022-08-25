package containers

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerservice/mgmt/2022-03-02-preview/containerservice"
	"github.com/Azure/go-autorest/autorest/azure"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
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
		aciConnectorKey,           // https://github.com/hashicorp/terraform-provider-azurerm/issues/5510
		httpApplicationRoutingKey, // https://github.com/hashicorp/terraform-provider-azurerm/issues/5960
	},
	azure.USGovernmentCloud.Name: {
		httpApplicationRoutingKey, // https://github.com/hashicorp/terraform-provider-azurerm/issues/5960
	},
}

func schemaKubernetesAddOns() map[string]*pluginsdk.Schema {
	out := map[string]*pluginsdk.Schema{
		"aci_connector_linux": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
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
		},
		"http_application_routing_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"http_application_routing_zone_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"oms_agent": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
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
		},
		"ingress_application_gateway": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
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
		},
		"key_vault_secrets_provider": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
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

	return out
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
