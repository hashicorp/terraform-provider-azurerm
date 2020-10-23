package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

var kubernetesDataSourceTests = map[string]func(t *testing.T){
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
	"nodeLabels":                                  testAccDataSourceAzureRMKubernetesCluster_nodeLabels,
	"enableNodePublicIP":                          testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIP,
	"privateCluster":                              testAccDataSourceAzureRMKubernetesCluster_privateCluster,
}

func TestAccDataSourceAzureRMKubernetesCluster_basic(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_basic(t)
}

func testAccDataSourceAzureRMKubernetesCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "false"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config_raw", ""),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kubelet_identity.0.object_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kubelet_identity.0.client_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kubelet_identity.0.user_assigned_identity_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "identity.0.type", "SystemAssigned"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_privateCluster(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_privateCluster(t)
}

func testAccDataSourceAzureRMKubernetesCluster_privateCluster(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_privateClusterConfig(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_fqdn"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_cluster_enabled", "true"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(t)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config_raw", ""),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(t)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantId := os.Getenv("ARM_TENANT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, clientId, clientSecret, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.client_app_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.server_app_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "kube_admin_config_raw"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_internalNetwork(t)
}

func testAccDataSourceAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_internalNetworkConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(t)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "calico"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_policy"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(t)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "azure"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_policy"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCompleteConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyCompleteConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "calico"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_policy"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyCompleteConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "azure"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_policy"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(t)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "kubenet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetCompleteConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "kubenet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(t)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMSConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.oms_agent.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.oms_agent.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "addon_profile.0.oms_agent.0.log_analytics_workspace_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "addon_profile.0.oms_agent.0.oms_agent_identity.0.client_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "addon_profile.0.oms_agent.0.oms_agent_identity.0.object_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "addon_profile.0.oms_agent.0.oms_agent_identity.0.user_assigned_identity_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboard(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboard(t)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboard(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboardConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.kube_dashboard.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.kube_dashboard.0.enabled", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicy(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicy(t)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicyConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.azure_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.azure_policy.0.enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.azure_policy.0.version", "v2"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(t)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_addOnProfileRoutingConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.http_application_routing.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.http_application_routing.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "addon_profile.0.http_application_routing.0.http_application_routing_zone_name"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_autoscalingNoAvailabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_autoscalingNoAvailabilityZones(t)
}

func testAccDataSourceAzureRMKubernetesCluster_autoscalingNoAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_autoScalingNoAvailabilityZonesConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.min_count", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.max_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.enable_auto_scaling", "true"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "agent_pool_profile.0.availability_zones"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_autoscalingWithAvailabilityZones(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_autoscalingWithAvailabilityZones(t)
}

func testAccDataSourceAzureRMKubernetesCluster_autoscalingWithAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_autoScalingWithAvailabilityZonesConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.min_count", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.max_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.enable_auto_scaling", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.availability_zones.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.availability_zones.0", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.availability_zones.1", "2"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_nodeLabels(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_nodeLabels(t)
}

func testAccDataSourceAzureRMKubernetesCluster_nodeLabels(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	labels := map[string]string{"key": "value"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_nodeLabelsConfig(data, labels),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.node_labels.key", "value"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_enableNodePublicIP(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIP(t)
}

func testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIPConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.enable_node_public_ip", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_basicConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_basicVMSSConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlConfig(data acceptance.TestData) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControlConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data acceptance.TestData, clientId, clientSecret, tenantId string) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, clientId, clientSecret, tenantId)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_internalNetworkConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_internalNetworkConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingConfig(data, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyConfig(data, "azure", "calico")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyConfig(data, "azure", "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCompleteConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingCompleteConfig(data, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyCompleteConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyCompleteConfig(data, "azure", "calico")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyCompleteConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyCompleteConfig(data, "azure", "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingConfig(data, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetCompleteConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingCompleteConfig(data, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMSConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_addonProfileOMSConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboardConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_addonProfileKubeDashboardConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicyConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_addonProfileAzurePolicyConfig(data, true)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileRoutingConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_addonProfileRoutingConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_autoScalingNoAvailabilityZonesConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZonesConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_autoScalingWithAvailabilityZonesConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZonesConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_nodeLabelsConfig(data acceptance.TestData, labels map[string]string) string {
	r := testAccAzureRMKubernetesCluster_nodeLabelsConfig(data, labels)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIPConfig(data acceptance.TestData) string {
	r := testAccAzureRMKubernetesCluster_enableNodePublicIPConfig(data, true)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, r)
}
