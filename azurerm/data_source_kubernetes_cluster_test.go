package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func testAccDataSourceAzureRMKubernetesCluster_basic(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_basicConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.password"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config_raw", ""),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlConfig(ri, location, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.azure_active_directory.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config_raw", ""),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantId := os.Getenv("ARM_TENANT_ID")
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(ri, location, clientId, clientSecret, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "role_based_access_control.0.azure_active_directory.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role_based_access_control.0.azure_active_directory.0.client_app_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role_based_access_control.0.azure_active_directory.0.server_app_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role_based_access_control.0.azure_active_directory.0.tenant_id"),
					resource.TestCheckResourceAttr(dataSourceName, "kube_admin_config.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_admin_config_raw"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_internalNetworkConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}
func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_policy", "calico"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_policy", "azure"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCompleteConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyCompleteConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_policy", "calico"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyCompleteConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_policy", "azure"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "kubenet"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetCompleteConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
					resource.TestCheckResourceAttr(dataSourceName, "network_profile.0.network_plugin", "kubenet"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.network_plugin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.dns_service_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.docker_bridge_cidr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_profile.0.service_cidr"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMSConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.oms_agent.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.oms_agent.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "addon_profile.0.oms_agent.0.log_analytics_workspace_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboard(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboardConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.kube_dashboard.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.kube_dashboard.0.enabled", "false"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicy(t *testing.T) {
	resourceName := "azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	config := testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicyConfig(ri, clientId, clientSecret, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "addon_profile.0.azure_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "addon_profile.0.azure_policy.0.enabled", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_addOnProfileRoutingConfig(ri, clientId, clientSecret, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.http_application_routing.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "addon_profile.0.http_application_routing.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "addon_profile.0.http_application_routing.0.http_application_routing_zone_name"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_autoscalingNoAvailabilityZones(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	config := testAccDataSourceAzureRMKubernetesCluster_autoScalingNoAvailabilityZonesConfig(ri, clientId, clientSecret, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.min_count", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.max_count", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.enable_auto_scaling", "true"),
					resource.TestCheckNoResourceAttr(dataSourceName, "agent_pool_profile.0.availability_zones"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_autoscalingWithAvailabilityZones(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	config := testAccDataSourceAzureRMKubernetesCluster_autoScalingWithAvailabilityZonesConfig(ri, clientId, clientSecret, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.min_count", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.max_count", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.type", "VirtualMachineScaleSets"),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.enable_auto_scaling", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.availability_zones.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.availability_zones.0", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.availability_zones.1", "2"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_nodeTaints(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	config := testAccDataSourceAzureRMKubernetesCluster_nodeTaintsConfig(ri, clientId, clientSecret, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.1.node_taints.0", "key=value:NoSchedule"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIP(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	config := testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIPConfig(ri, clientId, clientSecret, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "agent_pool_profile.0.enable_node_public_ip", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_basicConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_basicVMSSConfig(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlConfig(rInt int, location, clientId, clientSecret string) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControlConfig(rInt, location, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(rInt int, location, clientId, clientSecret, tenantId string) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControlAADConfig(rInt, location, clientId, clientSecret, tenantId)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_internalNetworkConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_internalNetworkConfig(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingConfig(rInt, clientId, clientSecret, location, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyConfig(rInt, clientId, clientSecret, location, "azure", "calico")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyConfig(rInt, clientId, clientSecret, location, "azure", "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCompleteConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingCompleteConfig(rInt, clientId, clientSecret, location, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyCompleteConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyCompleteConfig(rInt, clientId, clientSecret, location, "azure", "calico")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyCompleteConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyCompleteConfig(rInt, clientId, clientSecret, location, "azure", "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingConfig(rInt, clientId, clientSecret, location, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetCompleteConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingCompleteConfig(rInt, clientId, clientSecret, location, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMSConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileOMSConfig(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboardConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileKubeDashboardConfig(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileAzurePolicyConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileAzurePolicyConfig(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileRoutingConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileRoutingConfig(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_autoScalingNoAvailabilityZonesConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZonesConfig(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_autoScalingWithAvailabilityZonesConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZonesConfig(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_nodeTaintsConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_nodeTaintsConfig(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_enableNodePublicIPConfig(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_enableNodePublicIPConfig(rInt, clientId, clientSecret, location, true)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}
