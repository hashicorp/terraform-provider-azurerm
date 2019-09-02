package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMKubernetesCluster_basic(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_basic(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(t *testing.T) {
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
				Config: testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(ri, location, clientId, clientSecret),
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

func TestAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
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
				Config: testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(ri, location, clientId, clientSecret, tenantId),
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

func TestAccDataSourceAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_internalNetwork(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(ri, clientId, clientSecret, location)

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
func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboard(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboard(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(ri, clientId, clientSecret, location)

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

func TestAccDataSourceAzureRMKubernetesCluster_autoscalingNoAvailabilityZones(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	config := testAccDataSourceAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(ri, clientId, clientSecret, testLocation())

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

func TestAccDataSourceAzureRMKubernetesCluster_autoscalingWithAvailabilityZones(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	config := testAccDataSourceAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(ri, clientId, clientSecret, testLocation())

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

func TestAccDataSourceAzureRMKubernetesCluster_nodeTaints(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := tf.AccRandTimeInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	config := testAccDataSourceAzureRMKubernetesCluster_nodeTaints(ri, clientId, clientSecret, testLocation())

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

func testAccDataSourceAzureRMKubernetesCluster_basic(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_basic(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControl(rInt int, location, clientId, clientSecret string) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControl(rInt, location, clientId, clientSecret)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_roleBasedAccessControlAAD(rInt int, location, clientId, clientSecret, tenantId string) string {
	resource := testAccAzureRMKubernetesCluster_roleBasedAccessControlAAD(rInt, location, clientId, clientSecret, tenantId)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_internalNetwork(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_internalNetwork(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzure(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworking(rInt, clientId, clientSecret, location, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicy(rInt, clientId, clientSecret, location, "azure", "calico")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicy(rInt, clientId, clientSecret, location, "azure", "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureComplete(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingComplete(rInt, clientId, clientSecret, location, "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyComplete(rInt, clientId, clientSecret, location, "azure", "calico")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyComplete(rInt, clientId, clientSecret, location, "azure", "azure")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenet(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworking(rInt, clientId, clientSecret, location, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_advancedNetworkingComplete(rInt, clientId, clientSecret, location, "kubenet")
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileOMS(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileOMS(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileKubeDashboard(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileKubeDashboard(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_addOnProfileRouting(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_addonProfileRouting(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_autoScalingNoAvailabilityZones(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_autoscaleNoAvailabilityZones(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_autoScalingWithAvailabilityZones(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_autoscaleWithAvailabilityZones(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMKubernetesCluster_nodeTaints(rInt int, clientId string, clientSecret string, location string) string {
	r := testAccAzureRMKubernetesCluster_nodeTaints(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, r)
}
