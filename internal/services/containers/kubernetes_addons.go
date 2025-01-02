// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	commonValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	// note: the casing on these keys is important
	aciConnectorKey                 = "aciConnectorLinux"
	azureKeyvaultSecretsProviderKey = "azureKeyvaultSecretsProvider"
	azurePolicyKey                  = "azurepolicy"
	confidentialComputingKey        = "ACCSGXDevicePlugin"
	httpApplicationRoutingKey       = "httpApplicationRouting"
	ingressApplicationGatewayKey    = "ingressApplicationGateway"
	omsAgentKey                     = "omsagent"
	openServiceMeshKey              = "openServiceMesh"
)

// The AKS API hard-codes which add-ons are supported in which environment
// as such unfortunately we can't just send "disabled" - we need to strip
// the unsupported addons from the HTTP response. As such this defines
// the list of unsupported addons in the defined region - e.g. by being
// omitted from this list an addon/environment combination will be supported
var unsupportedAddonsForEnvironment = map[string][]string{
	environments.AzureChinaCloud: {
		aciConnectorKey,           // https://github.com/hashicorp/terraform-provider-azurerm/issues/5510
		httpApplicationRoutingKey, // https://github.com/hashicorp/terraform-provider-azurerm/issues/5960
	},
	environments.AzureUSGovernmentCloud: {
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
					"connector_identity": {
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
		"azure_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"confidential_computing": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"sgx_quote_helper_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
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
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
					"msi_auth_for_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
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
						ValidateFunc: applicationgateways.ValidateApplicationGatewayID,
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
						ValidateFunc: commonids.ValidateSubnetID,
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

func expandKubernetesAddOns(d *pluginsdk.ResourceData, input map[string]interface{}, env environments.Environment) (*map[string]managedclusters.ManagedClusterAddonProfile, error) {
	disabled := managedclusters.ManagedClusterAddonProfile{
		Enabled: false,
	}

	addonProfiles := map[string]managedclusters.ManagedClusterAddonProfile{}

	confidentialComputing := input["confidential_computing"].([]interface{})
	if len(confidentialComputing) > 0 && confidentialComputing[0] != nil {
		value := confidentialComputing[0].(map[string]interface{})
		config := make(map[string]string)
		quoteHelperEnabled := "false"
		if value["sgx_quote_helper_enabled"].(bool) {
			quoteHelperEnabled = "true"
		}
		config["ACCSGXQuoteHelperEnabled"] = quoteHelperEnabled
		addonProfiles[confidentialComputingKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	} else if len(confidentialComputing) == 0 && d.HasChange("confidential_computing") {
		addonProfiles[confidentialComputingKey] = disabled
	}

	if d.HasChange("http_application_routing_enabled") {
		addonProfiles[httpApplicationRoutingKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: input["http_application_routing_enabled"].(bool),
		}
	}

	omsAgent := input["oms_agent"].([]interface{})
	if len(omsAgent) > 0 && omsAgent[0] != nil {
		value := omsAgent[0].(map[string]interface{})
		config := make(map[string]string)

		if workspaceID, ok := value["log_analytics_workspace_id"]; ok && workspaceID != "" {
			lawid, err := workspaces.ParseWorkspaceIDInsensitively(workspaceID.(string))
			if err != nil {
				return nil, fmt.Errorf("parsing Log Analytics Workspace ID: %+v", err)
			}
			config["logAnalyticsWorkspaceResourceID"] = lawid.ID()
		}

		if useAADAuth, ok := value["msi_auth_for_monitoring_enabled"].(bool); ok {
			config["useAADAuth"] = fmt.Sprintf("%t", useAADAuth)
		}

		addonProfiles[omsAgentKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	} else if len(omsAgent) == 0 && d.HasChange("oms_agent") {
		addonProfiles[omsAgentKey] = disabled
	}

	aciConnector := input["aci_connector_linux"].([]interface{})
	if len(aciConnector) > 0 && aciConnector[0] != nil {
		value := aciConnector[0].(map[string]interface{})
		config := make(map[string]string)

		if subnetName, ok := value["subnet_name"]; ok && subnetName != "" {
			config["SubnetName"] = subnetName.(string)
		}

		addonProfiles[aciConnectorKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	} else if len(aciConnector) == 0 && d.HasChange("aci_connector_linux") {
		addonProfiles[aciConnectorKey] = disabled
	}

	if ok := d.HasChange("azure_policy_enabled"); ok {
		v := input["azure_policy_enabled"].(bool)
		props := managedclusters.ManagedClusterAddonProfile{
			Enabled: v,
			Config: pointer.To(map[string]string{
				"version": "v2",
			}),
		}
		addonProfiles[azurePolicyKey] = props
	}

	ingressApplicationGateway := input["ingress_application_gateway"].([]interface{})
	if len(ingressApplicationGateway) > 0 && ingressApplicationGateway[0] != nil {
		value := ingressApplicationGateway[0].(map[string]interface{})
		config := make(map[string]string)

		if gatewayId, ok := value["gateway_id"]; ok && gatewayId != "" {
			config["applicationGatewayId"] = gatewayId.(string)
		}

		if gatewayName, ok := value["gateway_name"]; ok && gatewayName != "" {
			config["applicationGatewayName"] = gatewayName.(string)
		}

		if subnetCIDR, ok := value["subnet_cidr"]; ok && subnetCIDR != "" {
			config["subnetCIDR"] = subnetCIDR.(string)
		}

		if subnetId, ok := value["subnet_id"]; ok && subnetId != "" {
			config["subnetId"] = subnetId.(string)
		}

		addonProfiles[ingressApplicationGatewayKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	} else if len(ingressApplicationGateway) == 0 && d.HasChange("ingress_application_gateway") {
		addonProfiles[ingressApplicationGatewayKey] = disabled
	}

	if ok := d.HasChange("open_service_mesh_enabled"); ok {
		addonProfiles[openServiceMeshKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: input["open_service_mesh_enabled"].(bool),
			Config:  nil,
		}
	}

	azureKeyVaultSecretsProvider := input["key_vault_secrets_provider"].([]interface{})
	if len(azureKeyVaultSecretsProvider) > 0 && azureKeyVaultSecretsProvider[0] != nil {
		value := azureKeyVaultSecretsProvider[0].(map[string]interface{})
		config := make(map[string]string)

		enableSecretRotation := fmt.Sprintf("%t", value["secret_rotation_enabled"].(bool))
		config["enableSecretRotation"] = enableSecretRotation
		config["rotationPollInterval"] = value["secret_rotation_interval"].(string)

		addonProfiles[azureKeyvaultSecretsProviderKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  &config,
		}
	} else if len(azureKeyVaultSecretsProvider) == 0 && d.HasChange("key_vault_secrets_provider") {
		addonProfiles[azureKeyvaultSecretsProviderKey] = disabled
	}

	return filterUnsupportedKubernetesAddOns(addonProfiles, env)
}

func filterUnsupportedKubernetesAddOns(input map[string]managedclusters.ManagedClusterAddonProfile, env environments.Environment) (*map[string]managedclusters.ManagedClusterAddonProfile, error) {
	filter := func(input map[string]managedclusters.ManagedClusterAddonProfile, key string) (map[string]managedclusters.ManagedClusterAddonProfile, error) {
		output := input
		if v, ok := output[key]; ok {
			if v.Enabled {
				return nil, fmt.Errorf("The addon %q is not supported for a Kubernetes Cluster located in %q", key, env.Name)
			}

			// otherwise it's disabled by default, so just remove it
			delete(output, key)
		}

		return output, nil
	}

	output := input
	if unsupportedAddons, ok := unsupportedAddonsForEnvironment[env.Name]; ok {
		for _, key := range unsupportedAddons {
			out, err := filter(output, key)
			if err != nil {
				return nil, err
			}

			output = out
		}
	}
	return &output, nil
}

func flattenKubernetesAddOns(profile map[string]managedclusters.ManagedClusterAddonProfile) map[string]interface{} {
	aciConnectors := make([]interface{}, 0)
	aciConnector := kubernetesAddonProfileLocate(profile, aciConnectorKey)
	if enabled := aciConnector.Enabled; enabled {
		subnetName := ""
		if v := aciConnector.Config; v != nil && (*v)["SubnetName"] != "" {
			subnetName = (*v)["SubnetName"]
		}

		identity := flattenKubernetesClusterAddOnIdentityProfile(aciConnector.Identity)

		aciConnectors = append(aciConnectors, map[string]interface{}{
			"subnet_name":        subnetName,
			"connector_identity": identity,
		})
	}

	azurePolicyEnabled := false
	azurePolicy := kubernetesAddonProfileLocate(profile, azurePolicyKey)
	if enabledVal := azurePolicy.Enabled; enabledVal {
		azurePolicyEnabled = enabledVal
	}

	confidentialComputings := make([]interface{}, 0)
	confidentialComputing := kubernetesAddonProfileLocate(profile, confidentialComputingKey)
	if enabled := confidentialComputing.Enabled; enabled {
		quoteHelperEnabled := false
		if v := kubernetesAddonProfilelocateInConfig(confidentialComputing.Config, "ACCSGXQuoteHelperEnabled"); v != "" && v != "false" {
			quoteHelperEnabled = true
		}
		confidentialComputings = append(confidentialComputings, map[string]interface{}{
			"sgx_quote_helper_enabled": quoteHelperEnabled,
		})
	}

	httpApplicationRoutingEnabled := false
	httpApplicationRoutingZone := ""
	httpApplicationRouting := kubernetesAddonProfileLocate(profile, httpApplicationRoutingKey)
	if enabledVal := httpApplicationRouting.Enabled; enabledVal {
		httpApplicationRoutingEnabled = enabledVal
	}

	if v := kubernetesAddonProfilelocateInConfig(httpApplicationRouting.Config, "HTTPApplicationRoutingZoneName"); v != "" {
		httpApplicationRoutingZone = v
	}

	omsAgents := make([]interface{}, 0)
	omsAgent := kubernetesAddonProfileLocate(profile, omsAgentKey)
	if enabled := omsAgent.Enabled; enabled {
		workspaceID := ""
		useAADAuth := false

		if v := kubernetesAddonProfilelocateInConfig(omsAgent.Config, "logAnalyticsWorkspaceResourceID"); v != "" {
			if lawid, err := workspaces.ParseWorkspaceID(v); err == nil {
				workspaceID = lawid.ID()
			}
		}

		if v := kubernetesAddonProfilelocateInConfig(omsAgent.Config, "useAADAuth"); v != "false" && v != "" {
			useAADAuth = true
		}

		omsAgentIdentity := flattenKubernetesClusterAddOnIdentityProfile(omsAgent.Identity)

		omsAgents = append(omsAgents, map[string]interface{}{
			"log_analytics_workspace_id":      workspaceID,
			"msi_auth_for_monitoring_enabled": useAADAuth,
			"oms_agent_identity":              omsAgentIdentity,
		})
	}

	ingressApplicationGateways := make([]interface{}, 0)
	ingressApplicationGateway := kubernetesAddonProfileLocate(profile, ingressApplicationGatewayKey)
	if enabled := ingressApplicationGateway.Enabled; enabled {
		gatewayId := ""

		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayId"); v != "" {
			gatewayId = v
		}

		gatewayName := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayName"); v != "" {
			gatewayName = v
		}

		effectiveGatewayId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "effectiveApplicationGatewayId"); v != "" {
			effectiveGatewayId = v
		}

		subnetCIDR := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetCIDR"); v != "" {
			subnetCIDR = v
		}

		subnetId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetId"); v != "" {
			subnetId = v
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

	openServiceMeshEnabled := false
	openServiceMesh := kubernetesAddonProfileLocate(profile, openServiceMeshKey)
	if enabledVal := openServiceMesh.Enabled; enabledVal {
		openServiceMeshEnabled = enabledVal
	}

	azureKeyVaultSecretsProviders := make([]interface{}, 0)
	azureKeyVaultSecretsProvider := kubernetesAddonProfileLocate(profile, azureKeyvaultSecretsProviderKey)
	if enabled := azureKeyVaultSecretsProvider.Enabled; enabled {
		enableSecretRotation := false
		rotationPollInterval := ""

		if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "enableSecretRotation"); v != "false" {
			enableSecretRotation = true
		}

		if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "rotationPollInterval"); v != "" {
			rotationPollInterval = v
		}

		azureKeyvaultSecretsProviderIdentity := flattenKubernetesClusterAddOnIdentityProfile(azureKeyVaultSecretsProvider.Identity)

		azureKeyVaultSecretsProviders = append(azureKeyVaultSecretsProviders, map[string]interface{}{
			"secret_rotation_enabled":  enableSecretRotation,
			"secret_rotation_interval": rotationPollInterval,
			"secret_identity":          azureKeyvaultSecretsProviderIdentity,
		})
	}

	return map[string]interface{}{
		"aci_connector_linux":                aciConnectors,
		"azure_policy_enabled":               azurePolicyEnabled,
		"confidential_computing":             confidentialComputings,
		"http_application_routing_enabled":   httpApplicationRoutingEnabled,
		"http_application_routing_zone_name": httpApplicationRoutingZone,
		"ingress_application_gateway":        ingressApplicationGateways,
		"key_vault_secrets_provider":         azureKeyVaultSecretsProviders,
		"oms_agent":                          omsAgents,
		"open_service_mesh_enabled":          openServiceMeshEnabled,
	}
}

func flattenKubernetesClusterAddOnIdentityProfile(profile *managedclusters.UserAssignedIdentity) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	identity := make([]interface{}, 0)
	clientID := ""
	if clientid := profile.ClientId; clientid != nil {
		clientID = *clientid
	}

	objectID := ""
	if objectid := profile.ObjectId; objectid != nil {
		objectID = *objectid
	}

	userAssignedIdentityID := ""
	if resourceid := profile.ResourceId; resourceid != nil {
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
		"confidential_computing":           d.Get("confidential_computing").([]interface{}),
		"http_application_routing_enabled": d.Get("http_application_routing_enabled").(bool),
		"oms_agent":                        d.Get("oms_agent").([]interface{}),
		"ingress_application_gateway":      d.Get("ingress_application_gateway").([]interface{}),
		"open_service_mesh_enabled":        d.Get("open_service_mesh_enabled").(bool),
		"key_vault_secrets_provider":       d.Get("key_vault_secrets_provider").([]interface{}),
	}
}

// when the Kubernetes Cluster is updated in the Portal - Azure updates the casing on the keys
// meaning what's submitted could be different to what's returned..
func kubernetesAddonProfileLocate(profile map[string]managedclusters.ManagedClusterAddonProfile, key string) managedclusters.ManagedClusterAddonProfile {
	for k, v := range profile {
		if strings.EqualFold(k, key) {
			return v
		}
	}

	return managedclusters.ManagedClusterAddonProfile{}
}

// when the Kubernetes Cluster is updated in the Portal - Azure updates the casing on the keys
// meaning what's submitted could be different to what's returned..
// Related issue: https://github.com/Azure/azure-rest-api-specs/issues/10716
func kubernetesAddonProfilelocateInConfig(config *map[string]string, key string) string {
	if config == nil {
		return ""
	}

	for k, v := range *config {
		if strings.EqualFold(k, key) {
			return v
		}
	}

	return ""
}
