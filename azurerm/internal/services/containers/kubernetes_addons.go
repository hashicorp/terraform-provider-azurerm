package containers

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-10-01/containerservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const (
	aciConnectorKey           = "aciConnectorLinux"
	azurePolicyKey            = "azurepolicy"
	kubernetesDashboardKey    = "kubeDashboard"
	httpApplicationRoutingKey = "httpApplicationRouting"
	omsAgentKey               = "omsagent"
)

func SchemaKubernetesAddOnProfiles() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"aci_connector_linux": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Required: true,
							},

							"subnet_name": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
					},
				},

				"azure_policy": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},

				"kube_dashboard": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Required: true,
							},
						},
					},
				},

				"http_application_routing": {
					Type:     schema.TypeList,
					MaxItems: 1,
					ForceNew: true,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								ForceNew: true,
								Required: true,
							},
							"http_application_routing_zone_name": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},

				"oms_agent": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:     schema.TypeBool,
								Required: true,
							},
							"log_analytics_workspace_id": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: azure.ValidateResourceID,
							},
						},
					},
				},
			},
		},
	}
}

func ExpandKubernetesAddOnProfiles(input []interface{}) map[string]*containerservice.ManagedClusterAddonProfile {
	disabled := containerservice.ManagedClusterAddonProfile{
		Enabled: utils.Bool(false),
	}

	profiles := map[string]*containerservice.ManagedClusterAddonProfile{
		// note: the casing on these keys is important
		aciConnectorKey:           &disabled,
		azurePolicyKey:            &disabled,
		kubernetesDashboardKey:    &disabled,
		httpApplicationRoutingKey: &disabled,
		omsAgentKey:               &disabled,
	}
	if len(input) == 0 {
		return profiles
	}

	profile := input[0].(map[string]interface{})
	addonProfiles := map[string]*containerservice.ManagedClusterAddonProfile{}

	httpApplicationRouting := profile["http_application_routing"].([]interface{})
	if len(httpApplicationRouting) > 0 && httpApplicationRouting[0] != nil {
		value := httpApplicationRouting[0].(map[string]interface{})
		enabled := value["enabled"].(bool)
		addonProfiles["httpApplicationRouting"] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
		}
	}

	omsAgent := profile["oms_agent"].([]interface{})
	if len(omsAgent) > 0 && omsAgent[0] != nil {
		value := omsAgent[0].(map[string]interface{})
		config := make(map[string]*string)
		enabled := value["enabled"].(bool)

		if workspaceId, ok := value["log_analytics_workspace_id"]; ok && workspaceId != "" {
			config["logAnalyticsWorkspaceResourceID"] = utils.String(workspaceId.(string))
		}

		addonProfiles["omsagent"] = &containerservice.ManagedClusterAddonProfile{
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

		addonProfiles["aciConnectorLinux"] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  config,
		}
	}

	kubeDashboard := profile["kube_dashboard"].([]interface{})
	if len(kubeDashboard) > 0 && kubeDashboard[0] != nil {
		value := kubeDashboard[0].(map[string]interface{})
		enabled := value["enabled"].(bool)

		addonProfiles["kubeDashboard"] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  nil,
		}
	}

	azurePolicy := profile["azure_policy"].([]interface{})
	if len(azurePolicy) > 0 && azurePolicy[0] != nil {
		value := azurePolicy[0].(map[string]interface{})
		enabled := value["enabled"].(bool)

		addonProfiles["azurepolicy"] = &containerservice.ManagedClusterAddonProfile{
			Enabled: utils.Bool(enabled),
			Config:  nil,
		}
	}

	return addonProfiles
}

func FlattenKubernetesAddOnProfiles(profile map[string]*containerservice.ManagedClusterAddonProfile) []interface{} {
	// when the Kubernetes Cluster is updated in the Portal - Azure updates the casing on the keys
	// meaning what's submitted could be different to what's returned..
	var locateInProfile = func(key string) *containerservice.ManagedClusterAddonProfile {
		for k, v := range profile {
			if strings.EqualFold(k, key) {
				return v
			}
		}

		return nil
	}

	aciConnectors := make([]interface{}, 0)
	if aciConnector := locateInProfile(aciConnectorKey); aciConnector != nil {
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
	if azurePolicy := locateInProfile(azurePolicyKey); azurePolicy != nil {
		enabled := false
		if enabledVal := azurePolicy.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		azurePolicies = append(azurePolicies, map[string]interface{}{
			"enabled": enabled,
		})
	}

	httpApplicationRoutes := make([]interface{}, 0)
	if httpApplicationRouting := locateInProfile(httpApplicationRoutingKey); httpApplicationRouting != nil {
		enabled := false
		if enabledVal := httpApplicationRouting.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		zoneName := ""
		if v := httpApplicationRouting.Config["HTTPApplicationRoutingZoneName"]; v != nil {
			zoneName = *v
		}

		httpApplicationRoutes = append(httpApplicationRoutes, map[string]interface{}{
			"enabled":                            enabled,
			"http_application_routing_zone_name": zoneName,
		})
	}

	kubeDashboards := make([]interface{}, 0)
	if kubeDashboard := locateInProfile(kubernetesDashboardKey); kubeDashboard != nil {
		enabled := false
		if enabledVal := kubeDashboard.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		kubeDashboards = append(kubeDashboards, map[string]interface{}{
			"enabled": enabled,
		})
	}

	omsAgents := make([]interface{}, 0)
	if omsAgent := locateInProfile(omsAgentKey); omsAgent != nil {
		enabled := false
		if enabledVal := omsAgent.Enabled; enabledVal != nil {
			enabled = *enabledVal
		}

		workspaceId := ""
		if workspaceResourceID := omsAgent.Config["logAnalyticsWorkspaceResourceID"]; workspaceResourceID != nil {
			workspaceId = *workspaceResourceID
		}

		omsAgents = append(omsAgents, map[string]interface{}{
			"enabled":                    enabled,
			"log_analytics_workspace_id": workspaceId,
		})
	}

	// this is a UX hack, since if the top level block isn't defined everything should be turned off
	if len(aciConnectors) == 0 && len(azurePolicies) == 0 && len(httpApplicationRoutes) == 0 && len(kubeDashboards) == 0 && len(omsAgents) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"aci_connector_linux":      aciConnectors,
			"azure_policy":             azurePolicies,
			"http_application_routing": httpApplicationRoutes,
			"kube_dashboard":           kubeDashboards,
			"oms_agent":                omsAgents,
		},
	}
}
