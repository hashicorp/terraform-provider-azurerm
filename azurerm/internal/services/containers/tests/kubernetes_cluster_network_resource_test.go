package tests

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

var kubernetesNetworkAuthTests = map[string]func(t *testing.T){
	"enableNodePublicIP":                            testAccAzureRMKubernetesCluster_enableNodePublicIP,
	"internalNetwork":                               testAccAzureRMKubernetesCluster_internalNetwork,
	"changingLoadBalancerProfile":                   testAccAzureRMKubernetesCluster_changingLoadBalancerProfile,
	"prefixedLoadBalancerProfile":                   testAccAzureRMKubernetesCluster_prefixedLoadBalancerProfile,
	"standardLoadBalancer":                          testAccAzureRMKubernetesCluster_standardLoadBalancer,
	"standardLoadBalancerComplete":                  testAccAzureRMKubernetesCluster_standardLoadBalancerComplete,
	"standardLoadBalancerProfile":                   testAccAzureRMKubernetesCluster_standardLoadBalancerProfile,
	"standardLoadBalancerProfileComplete":           testAccAzureRMKubernetesCluster_standardLoadBalancerProfileComplete,
	"advancedNetworkingKubenet":                     testAccAzureRMKubernetesCluster_advancedNetworkingKubenet,
	"advancedNetworkingKubenetComplete":             testAccAzureRMKubernetesCluster_advancedNetworkingKubenetComplete,
	"advancedNetworkingAzure":                       testAccAzureRMKubernetesCluster_advancedNetworkingAzure,
	"advancedNetworkingAzureComplete":               testAccAzureRMKubernetesCluster_advancedNetworkingAzureComplete,
	"advancedNetworkingAzureCalicoPolicy":           testAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy,
	"advancedNetworkingAzureCalicoPolicyComplete":   testAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete,
	"advancedNetworkingAzureNPMPolicy":              testAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy,
	"advancedNetworkingAzureNPMPolicyComplete":      testAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete,
	"basicLoadBalancerProfile":                      testAccAzureRMKubernetesCluster_basicLoadBalancerProfile,
	"standardLoadBalancerProfileWithPortAndTimeout": testAccAzureRMKubernetesCluster_standardLoadBalancerProfileWithPortAndTimeout,
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingKubenet(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_advancedNetworkingKubenet(t)
}

func testAccAzureRMKubernetesCluster_advancedNetworkingKubenet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingConfig(data, "kubenet"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "kubenet"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t)
}

func testAccAzureRMKubernetesCluster_advancedNetworkingKubenetComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingCompleteConfig(data, "kubenet"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "kubenet"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingAzure(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_advancedNetworkingAzure(t)
}

func testAccAzureRMKubernetesCluster_advancedNetworkingAzure(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingConfig(data, "azure"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t)
}

func testAccAzureRMKubernetesCluster_advancedNetworkingAzureComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingCompleteConfig(data, "azure"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t)
}

func testAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyConfig(data, "azure", "calico"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "calico"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t)
}

func testAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyCompleteConfig(data, "azure", "calico"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "calico"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyNetworkModeTransparent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyNetworkMode(data, "azure", "calico", "transparent"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "calico"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_mode", "transparent"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingAzureCalicoPolicyNetworkModeBridge(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyNetworkMode(data, "azure", "calico", "bridge"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "calico"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_mode", "bridge"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(t)
}

func testAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyConfig(data, "azure", "azure"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "azure"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t)
}

func testAccAzureRMKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyCompleteConfig(data, "azure", "azure"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_plugin", "azure"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.network_policy", "azure"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_enableNodePublicIP(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_enableNodePublicIP(t)
}

func testAccAzureRMKubernetesCluster_enableNodePublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				// Enabled
				Config: testAccAzureRMKubernetesCluster_enableNodePublicIPConfig(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.enable_node_public_ip", "true"),
				),
			},
			data.ImportStep(),
			{
				// Disabled
				Config: testAccAzureRMKubernetesCluster_enableNodePublicIPConfig(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.enable_node_public_ip", "false"),
				),
			},
			data.ImportStep(),
			{
				// Enabled
				Config: testAccAzureRMKubernetesCluster_enableNodePublicIPConfig(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.enable_node_public_ip", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_internalNetwork(t)
}

func testAccAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_internalNetworkConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_node_pool.0.max_pods", "60"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_outboundTypeLoadBalancer(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_outboundTypeLoadBalancer(t)
}

func testAccAzureRMKubernetesCluster_outboundTypeLoadBalancer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_outboundTypeLoadBalancerConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_privateClusterOn(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_privateClusterOn(t)
}

func testAccAzureRMKubernetesCluster_privateClusterOn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
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
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_privateClusterOff(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_privateClusterOff(t)
}

func testAccAzureRMKubernetesCluster_privateClusterOff(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_privateClusterConfig(data, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "private_cluster_enabled", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_standardLoadBalancer(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_standardLoadBalancer(t)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_standardLoadBalancerConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_standardLoadBalancerComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_standardLoadBalancerComplete(t)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancerComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_standardLoadBalancerCompleteConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_standardLoadBalancerProfile(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_standardLoadBalancerProfile(t)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancerProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_standardLoadBalancerProfileConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.effective_outbound_ips.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.idle_timeout_in_minutes", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.outbound_ports_allocated", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_standardLoadBalancerProfileComplete(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_standardLoadBalancerProfileComplete(t)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancerProfileComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.PreCheck(t) },
		Providers:                 acceptance.SupportedProviders,
		CheckDestroy:              testCheckAzureRMKubernetesClusterDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_standardLoadBalancerProfileCompleteConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.effective_outbound_ips.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_standardLoadBalancerProfileWithPortAndTimeout(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_standardLoadBalancerProfileWithPortAndTimeout(t)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancerProfileWithPortAndTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.PreCheck(t) },
		Providers:                 acceptance.SupportedProviders,
		CheckDestroy:              testCheckAzureRMKubernetesClusterDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_standardLoadBalancerProfileWithPortAndTimeoutConfig(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.outbound_ports_allocated", "8000"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.idle_timeout_in_minutes", "10"),
				),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_basicLoadBalancerProfile(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_basicLoadBalancerProfile(t)
}

func testAccAzureRMKubernetesCluster_basicLoadBalancerProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMKubernetesCluster_basicLoadBalancerProfileConfig(data),
				ExpectError: regexp.MustCompile("only load balancer SKU 'Standard' supports load balancer profiles. Provided load balancer type: basic"),
			},
		},
	})
}

func TestAccAzureRMKubernetesCluster_prefixedLoadBalancerProfile(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_prefixedLoadBalancerProfile(t)
}

func testAccAzureRMKubernetesCluster_prefixedLoadBalancerProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_prefixedLoadBalancerProfileConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.effective_outbound_ips.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMKubernetesCluster_changingLoadBalancerProfile(t *testing.T) {
	checkIfShouldRunTestsIndividually(t)
	testAccAzureRMKubernetesCluster_changingLoadBalancerProfile(t)
}

func testAccAzureRMKubernetesCluster_changingLoadBalancerProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMKubernetesCluster_changingLoadBalancerProfileConfigIPPrefix(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.effective_outbound_ips.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_changingLoadBalancerProfileConfigManagedIPs(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.effective_outbound_ips.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_changingLoadBalancerProfileConfigIPIds(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_sku", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.effective_outbound_ips.#", "1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMKubernetesCluster_changingLoadBalancerProfileConfigIPPrefix(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_profile.0.load_balancer_profile.0.effective_outbound_ips.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMKubernetesCluster_advancedNetworkingConfig(data acceptance.TestData, networkPlugin string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, networkPlugin)
}

func testAccAzureRMKubernetesCluster_advancedNetworkingCompleteConfig(data acceptance.TestData, networkPlugin string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin     = "%s"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, networkPlugin)
}

// nolint unparam
func testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyConfig(data acceptance.TestData, networkPlugin string, networkPolicy string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin = "%s"
    network_policy = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, networkPlugin, networkPolicy)
}

// nolint unparam
func testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyCompleteConfig(data acceptance.TestData, networkPlugin string, networkPolicy string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "akc-routetable-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name                   = "akc-route-%d"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin     = "%s"
    network_policy     = "%s"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, networkPlugin, networkPolicy)
}

// nolint unparam
func testAccAzureRMKubernetesCluster_advancedNetworkingWithPolicyNetworkMode(data acceptance.TestData, networkPlugin string, networkPolicy string, networkMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "akc-routetable-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name                   = "akc-route-%d"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_mode       = "%s"
    network_plugin     = "%s"
    network_policy     = "%s"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, networkMode, networkPlugin, networkPolicy)
}

func testAccAzureRMKubernetesCluster_enableNodePublicIPConfig(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name                  = "default"
    node_count            = 1
    vm_size               = "Standard_DS2_v2"
    enable_node_public_ip = %t
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, enabled)
}

func testAccAzureRMKubernetesCluster_internalNetworkConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["172.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "172.0.2.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
    max_pods       = 60
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_outboundTypeLoadBalancerConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 2
    vm_size    = "Standard_DS2_v2"
    max_pods   = 60
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin     = "kubenet"
    load_balancer_sku  = "Basic"
    pod_cidr           = "10.244.0.0/16"
    service_cidr       = "10.0.0.0/16"
    dns_service_ip     = "10.0.0.10"
    docker_bridge_cidr = "172.17.0.1/16"
    outbound_type      = "loadBalancer"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_privateClusterConfig(data acceptance.TestData, enablePrivateCluster bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                    = "acctestaks%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  dns_prefix              = "acctestaks%d"
  private_cluster_enabled = %t

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, enablePrivateCluster, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancerConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesVersion, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancerCompleteConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "akc-routetable-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name                   = "akc-route-%d"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin     = "azure"
    dns_service_ip     = "10.10.0.10"
    docker_bridge_cidr = "172.18.0.1/16"
    service_cidr       = "10.10.0.0/16"
    load_balancer_sku  = "Standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancerProfileConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "Standard"
    load_balancer_profile {
      managed_outbound_ip_count = 3
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesVersion, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancerProfileCompleteConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "Standard"
    load_balancer_profile {
      outbound_ip_address_ids = [azurerm_public_ip.test.id]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesVersion, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_standardLoadBalancerProfileWithPortAndTimeoutConfig(data acceptance.TestData, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }

  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "Standard"
    load_balancer_profile {
      managed_outbound_ip_count = 2
      outbound_ports_allocated  = 8000
      idle_timeout_in_minutes   = 10
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, currentKubernetesVersion, data.RandomInteger, clientId, clientSecret)
}

func testAccAzureRMKubernetesCluster_basicLoadBalancerProfileConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "basic"
    load_balancer_profile {
      managed_outbound_ip_count = 3
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesVersion, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_prefixedLoadBalancerProfileConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_public_ip_prefix" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestipprefix%d"
  prefix_length       = 31
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "Standard"
    load_balancer_profile {
      outbound_ip_prefix_ids = [azurerm_public_ip_prefix.test.id]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesVersion, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_changingLoadBalancerProfileConfigIPPrefix(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_public_ip_prefix" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestipprefix%d"
  prefix_length       = 31
}

resource "azurerm_public_ip" "test" {
  name                = "acctestipone%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"
    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "standard"
    load_balancer_profile {
      outbound_ip_prefix_ids = [azurerm_public_ip_prefix.test.id]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesVersion, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_changingLoadBalancerProfileConfigManagedIPs(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_public_ip_prefix" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestipprefix%d"
  prefix_length       = 31
}

resource "azurerm_public_ip" "test" {
  name                = "acctestipone%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"
    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "standard"
    load_balancer_profile {
      managed_outbound_ip_count = "1"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesVersion, data.RandomInteger)
}

func testAccAzureRMKubernetesCluster_changingLoadBalancerProfileConfigIPIds(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_public_ip_prefix" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestipprefix%d"
  prefix_length       = 31
}

resource "azurerm_public_ip" "test" {
  name                = "acctestipone%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"
    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vm_size        = "Standard_DS2_v2"
    vnet_subnet_id = azurerm_subnet.test.id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "standard"
    load_balancer_profile {
      outbound_ip_address_ids = [azurerm_public_ip.test.id]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesVersion, data.RandomInteger)
}
