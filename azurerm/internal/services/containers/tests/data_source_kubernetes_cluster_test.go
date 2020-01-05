package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func testAccDataSourceAzureRMKubernetesCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_basicConfig(data, clientId, clientSecret),
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
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_privateLink(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	privateIpAddressCdir := "10.0.0.0/8"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_privateLinkConfig(data, clientId, clientSecret, data.Locations.Primary, privateIpAddressCdir, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_fqdn"),
					resource.TestCheckResourceAttr(data.ResourceName, "api_server_authorized_ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_link_enabled", "true"),
				),
			},
			data.ImportStep("service_principal.0.client_secret"),
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantId := os.Getenv("ARM_TENANT_ID")

	resource.ParallelTest(t, resource.TestCase{
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

func testAccDataSourceAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_internalNetworkConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "agent_pool_profile.0.vnet_subnet_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureConfig(data, clientId, clientSecret),
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
func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCompleteConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyCompleteConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyCompleteConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetCompleteConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMSConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.oms_agent.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.oms_agent.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "addon_profile.0.oms_agent.0.log_analytics_workspace_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboard(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboardConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicyConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.azure_policy.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "addon_profile.0.azure_policy.0.enabled", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_addOnProfileRoutingConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_autoscalingNoAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_autoScalingNoAvailabilityZonesConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_autoscalingWithAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_autoScalingWithAvailabilityZonesConfig(data, clientId, clientSecret),
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

func testAccDataSourceAzureRMKubernetesCluster_nodeTaints(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_nodeTaintsConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.1.node_taints.0", "key=value:NoSchedule"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIPConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "agent_pool_profile.0.enable_node_public_ip", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_basicConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_basicVMSSConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlConfig(data acceptance.TestData, clientId, clientSecret string) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControlConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data acceptance.TestData, clientId, clientSecret, tenantId string) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(data, clientId, clientSecret, tenantId)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_internalNetworkConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_internalNetworkConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingConfig(data, clientId, clientSecret, data.Locations.Primary, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyConfig(data, clientId, clientSecret, data.Locations.Primary, "azure", "calico")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyConfig(data, clientId, clientSecret, data.Locations.Primary, "azure", "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCompleteConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingCompleteConfig(data, clientId, clientSecret, data.Locations.Primary, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyCompleteConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyCompleteConfig(data, clientId, clientSecret, data.Locations.Primary, "azure", "calico")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyCompleteConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyCompleteConfig(data, clientId, clientSecret, data.Locations.Primary, "azure", "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingConfig(data, clientId, clientSecret, data.Locations.Primary, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetCompleteConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingCompleteConfig(data, clientId, clientSecret, data.Locations.Primary, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMSConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileOMSConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboardConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileKubeDashboardConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicyConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileAzurePolicyConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileRoutingConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileRoutingConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_autoScalingNoAvailabilityZonesConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZonesConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_autoScalingWithAvailabilityZonesConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZonesConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_nodeTaintsConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_nodeTaintsConfig(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIPConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	r := testAccAzureRMKubernetesCluster_enableNodePublicIPConfig(data, clientId, clientSecret, data.Locations.Primary, true)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}
