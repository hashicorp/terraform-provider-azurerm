package containers

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-03-01/containerservice"
	"github.com/Azure/go-autorest/autorest/azure"
	commonValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	laparse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	logAnalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	applicationGatewayValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	subnetValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const (
	// note: the casing on these keys is important
	aciConnectorKey              = "aciConnectorLinux"
	azurePolicyKey               = "azurepolicy"
	kubernetesDashboardKey       = "kubeDashboard"
	httpApplicationRoutingKey    = "httpApplicationRouting"
	omsAgentKey                  = "omsagent"
	ingressApplicationGatewayKey = "ingressApplicationGateway"
)

// The AKS API hard-codes which add-ons are supported in which environment
// as such unfortunately we can't just send "disabled" - we need to strip
// the unsupported addons from the HTTP response. As such this defines
// the list of unsupported addons in the defined region - e.g. by being
// omitted from this list an addon/environment combination will be supported
var unsupportedAddonsForEnvironment = map[string][]string{
	azure.ChinaCloud.Name: {
		aciConnectorKey,           // https://github.com/terraform-providers/terraform-provider-azurerm/issues/5510
		httpApplicationRoutingKey, // https://github.com/terraform-providers/terraform-provider-azurerm/issues/5960
		kubernetesDashboardKey,    // https://github.com/terraform-providers/terraform-provider-azurerm/issues/7487
	},
	azure.USGovernmentCloud.Name: {
		azurePolicyKey,            // https://github.com/terraform-providers/terraform-provider-azurerm/issues/6702
		httpApplicationRoutingKey, // https://github.com/terraform-providers/terraform-provider-azurerm/issues/5960
		kubernetesDashboardKey,    // https://github.com/terraform-providers/terraform-provider-azurerm/issues/7136
	},
}

func schemaKubernetesAddOnProfiles() *pluginsdk.Schema {
	//lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"aci_connector_linux": {
					Type:     pluginsdk.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:     pluginsdk.TypeBool,
								Required: true,
							},

							"subnet_name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"azure_policy": {
					Type:     pluginsdk.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:     pluginsdk.TypeBool,
								Required: true,
							},
						},
					},
				},

				"kube_dashboard": {
					Type:     pluginsdk.TypeList,
					MaxItems: 1,
					Optional: true,
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
					Type:     pluginsdk.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:     pluginsdk.TypeBool,
								Required: true,
							},
							"http_application_routing_zone_name": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"oms_agent": {
					Type:     pluginsdk.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"enabled": {
								Type:     pluginsdk.TypeBool,
								Required: true,
							},
							"log_analytics_workspace_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
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
							"enabled": {
								Type:     pluginsdk.TypeBool,
								Required: true,
							},
							"gateway_id": {
								Type:          pluginsdk.TypeString,
								Optional:      true,
								ConflictsWith: []string{"addon_profile.0.ingress_application_gateway.0.subnet_cidr", "addon_profile.0.ingress_application_gateway.0.subnet_id"},
								ValidateFunc:  applicationGatewayValidate.ApplicationGatewayID,
							},
							"gateway_name": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"subnet_cidr": {
								Type:          pluginsdk.TypeString,
								Optional:      true,
								ConflictsWith: []string{"addon_profile.0.ingress_application_gateway.0.gateway_id", "addon_profile.0.ingress_application_gateway.0.subnet_id"},
								ValidateFunc:  commonValidate.CIDR,
							},
							"subnet_id": {
								Type:          pluginsdk.TypeString,
								Optional:      true,
								ConflictsWith: []string{"addon_profile.0.ingress_application_gateway.0.gateway_id", "addon_profile.0.ingress_application_gateway.0.subnet_cidr"},
								ValidateFunc:  subnetValidate.SubnetID,
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
			},
		},
	}
}

func expandKubernetesAddOnProfiles(input []interface{}, env azure.Environment) (*map[string]*containerservice.ManagedClusterAddonProfile, error) {
	disabled := containerservice.ManagedClusterAddonProfile{
		Enabled: utils.Bool(false),
	}

	profiles := map[string]*containerservice.ManagedClusterAddonProfile{
		aciConnectorKey:              &disabled,
		azurePolicyKey:               &disabled,
		kubernetesDashboardKey:       &disabled,
		httpApplicationRoutingKey:    &disabled,
		omsAgentKey:                  &disabled,
		ingressApplicationGatewayKey: &disabled,
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

	// this is a UX hack, since if the top level block isn't defined everything should be turned off
	if len(aciConnectors) == 0 && len(azurePolicies) == 0 && len(httpApplicationRoutes) == 0 && len(kubeDashboards) == 0 && len(omsAgents) == 0 && len(ingressApplicationGateways) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"aci_connector_linux":         aciConnectors,
			"azure_policy":                azurePolicies,
			"http_application_routing":    httpApplicationRoutes,
			"kube_dashboard":              kubeDashboards,
			"oms_agent":                   omsAgents,
			"ingress_application_gateway": ingressApplicationGateways,
		},
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
