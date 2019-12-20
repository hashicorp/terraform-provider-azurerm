package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

var olderKubernetesVersion = "1.14.8"
var currentKubernetesVersion = "1.15.5"

func TestAccAzureRMKubernetes_all(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests to
	// ease the load on the kubernetes api
	testCases := map[string]map[string]func(t *testing.T){
		"clusterAddOn": {
			"addonProfileAciConnectorLinux":         testAccAzureRMKubernetesCluster_addonProfileAciConnectorLinux,
			"addonProfileAciConnectorLinuxDisabled": testAccAzureRMKubernetesCluster_addonProfileAciConnectorLinuxDisabled,
			"addonProfileAzurePolicy":               testAccAzureRMKubernetesCluster_addonProfileAzurePolicy,
			"addonProfileKubeDashboard":             testAccAzureRMKubernetesCluster_addonProfileKubeDashboard,
			"addonProfileOMS":                       testAccAzureRMKubernetesCluster_addonProfileOMS,
			"addonProfileOMSToggle":                 testAccAzureRMKubernetesCluster_addonProfileOMSToggle,
			"addonProfileRouting":                   testAccAzureRMKubernetesCluster_addonProfileRouting,
		},
		"auth": {
			"apiServerAuthorizedIPRanges": testAccAzureRMKubernetesCluster_apiServerAuthorizedIPRanges,
			"enablePodSecurityPolicy":     testAccAzureRMKubernetesCluster_enablePodSecurityPolicy,
			"roleBasedAccessControl":      testAccAzureRMKubernetesCluster_roleBasedAccessControl,
			"roleBasedAccessControlAAD":   testAccAzureRMKubernetesCluster_roleBasedAccessControlAAD,
		},
		"legacy": {
			"legacyAgentPoolProfileAvailabilitySet": testAccAzureRMKubernetesCluster_legacyAgentPoolProfileAvailabilitySet,
			"legacyAgentPoolProfileVMSS":            testAccAzureRMKubernetesCluster_legacyAgentPoolProfileVMSS,
		},
		"network": {
			"advancedNetworkingKubenet":                   testAccAzureRMKubernetesCluster_advancedNetworkingKubenet,
			"advancedNetworkingKubenetComplete":           testAccAzureRMKubernetesCluster_advancedNetworkingKubenetComplete,
			"advancedNetworkingAzure":                     testAccAzureRMKubernetesCluster_advancedNetworkingAzure,
			"advancedNetworkingAzureComplete":             testAccAzureRMKubernetesCluster_advancedNetworkingAzureComplete,
			"advancedNetworkingAzureCalicoPolicy":         testAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy,
			"advancedNetworkingAzureCalicoPolicyComplete": testAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete,
			"advancedNetworkingAzureNPMPolicy":            testAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy,
			"advancedNetworkingAzureNPMPolicyComplete":    testAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete,
			"enableNodePublicIP":                          testAccAzureRMKubernetesCluster_enableNodePublicIP,
			"internalNetwork":                             testAccAzureRMKubernetesCluster_internalNetwork,
			"standardLoadBalancer":                        testAccAzureRMKubernetesCluster_standardLoadBalancer,
			"standardLoadBalancerComplete":                testAccAzureRMKubernetesCluster_standardLoadBalancerComplete,
		},
		"nodePool": {
			"autoScale":                      testAccAzureRMKubernetesClusterNodePool_autoScale,
			"autoScaleUpdate":                testAccAzureRMKubernetesClusterNodePool_autoScaleUpdate,
			"availabilityZones":              testAccAzureRMKubernetesClusterNodePool_availabilityZones,
			"errorForAvailabilitySet":        testAccAzureRMKubernetesClusterNodePool_errorForAvailabilitySet,
			"multiplePools":                  testAccAzureRMKubernetesClusterNodePool_multiplePools,
			"manualScale":                    testAccAzureRMKubernetesClusterNodePool_manualScale,
			"manualScaleMultiplePools":       testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePools,
			"manualScaleMultiplePoolsUpdate": testAccAzureRMKubernetesClusterNodePool_manualScaleMultiplePoolsUpdate,
			"manualScaleUpdate":              testAccAzureRMKubernetesClusterNodePool_manualScaleUpdate,
			"manualScaleVMSku":               testAccAzureRMKubernetesClusterNodePool_manualScaleVMSku,
			"nodePublicIP":                   testAccAzureRMKubernetesClusterNodePool_nodePublicIP,
			"nodeTaints":                     testAccAzureRMKubernetesClusterNodePool_nodeTaints,
			"requiresImport":                 testAccAzureRMKubernetesClusterNodePool_requiresImport,
			"osDiskSizeGB":                   testAccAzureRMKubernetesClusterNodePool_osDiskSizeGB,
			"virtualNetworkAutomatic":        testAccAzureRMKubernetesClusterNodePool_virtualNetworkAutomatic,
			"virtualNetworkManual":           testAccAzureRMKubernetesClusterNodePool_virtualNetworkManual,
			"windows":                        testAccAzureRMKubernetesClusterNodePool_windows,
			"windowsAndLinux":                testAccAzureRMKubernetesClusterNodePool_windowsAndLinux,
		},
		"other": {
			"basicAvailabilitySet":   testAccAzureRMKubernetesCluster_basicAvailabilitySet,
			"basicVMSS":              testAccAzureRMKubernetesCluster_basicVMSS,
			"requiresImport":         testAccAzureRMKubernetesCluster_requiresImport,
			"linuxProfile":           testAccAzureRMKubernetesCluster_linuxProfile,
			"nodeTaints":             testAccAzureRMKubernetesCluster_nodeTaints,
			"nodeResourceGroup":      testAccAzureRMKubernetesCluster_nodeResourceGroup,
			"upgradeConfig":          testAccAzureRMKubernetesCluster_upgrade,
			"tags":                   testAccAzureRMKubernetesCluster_tags,
			"windowsProfile":         testAccAzureRMKubernetesCluster_windowsProfile,
			"privateLinkOn":          testAccAzureRMKubernetesCluster_privateLinkOn,
			"updatePublicRangesOff":  testAccAzureRMKubernetesCluster_privateLinkOff,
			"managedClusterIdentity": testAccAzureRMKubernetesCluster_managedClusterIdentiy,
		},
		"scaling": {
			"addAgent":                         testAccAzureRMKubernetesCluster_addAgent,
			"removeAgent":                      testAccAzureRMKubernetesCluster_removeAgent,
			"autoScalingNodeCountUnset":        testAccAzureRMKubernetesCluster_autoScalingNodeCountUnset,
			"autoScalingNoAvailabilityZones":   testAccAzureRMKubernetesCluster_autoScalingNoAvailabilityZones,
			"autoScalingWithAvailabilityZones": testAccAzureRMKubernetesCluster_autoScalingWithAvailabilityZones,
		},
		"datasource": {
			"basic":                                       testAccDataSourceAzureRMKubernetesCluster_basic,
			"roleBasedAccessControl":                      testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl,
			"roleBasedAccessControlAAD":                   testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD,
			"internalNetwork":                             testAccDataSourceAzureRMKubernetesCluster_internalNetwork,
			"advancedNetworkingAzure":                     testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure,
			"advancedNetworkingAzureCalicoPolicy":         testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy,
			"advancedNetworkingAzureNPMPolicy":            testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy,
			"advancedNetworkingAzureComplete":             testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete,
			"advancedNetworkingAzureCalicoPolicyComplete": testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete,
			"advancedNetworkingAzureNPMPolicyComplete":    testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete,
			"advancedNetworkingKubenet":                   testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet,
			"advancedNetworkingKubenetComplete":           testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete,
			"addOnProfileOMS":                             testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS,
			"addOnProfileKubeDashboard":                   testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboard,
			"addOnProfileAzurePolicy":                     testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicy,
			"addOnProfileRouting":                         testAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting,
			"autoscalingNoAvailabilityZones":              testAccDataSourceAzureRMKubernetesCluster_autoscalingNoAvailabilityZones,
			"autoscalingWithAvailabilityZones":            testAccDataSourceAzureRMKubernetesCluster_autoscalingWithAvailabilityZones,
			"nodeTaints":                                  testAccDataSourceAzureRMKubernetesCluster_nodeTaints,
			"enableNodePublicIP":                          testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIP,
			"privateLink":                                 testAccDataSourceAzureRMKubernetesCluster_privateLink,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testCheckAzureRMKubernetesClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Managed Kubernetes Cluster: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.KubernetesClustersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		aks, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on kubernetesClustersClient: %+v", err)
		}

		if aks.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Managed Kubernetes Cluster %q (Resource Group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMKubernetesClusterDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Containers.KubernetesClustersClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_kubernetes_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Managed Kubernetes Cluster still exists:\n%#v", resp)
		}
	}

	return nil
}
