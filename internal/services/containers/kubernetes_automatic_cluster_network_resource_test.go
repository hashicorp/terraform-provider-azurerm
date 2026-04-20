// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccKubernetesAutomaticCluster_advancedNetworkingBlock(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingBlock(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.advancedNetworkingBlockUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.advancedNetworkingBlock(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.advancedNetworkingBlockRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.advancedNetworkingBlock(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_advancedNetworkingNetworkDataplane(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingBlock(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_serviceMeshProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceMeshProfile(data, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceMeshProfile(data, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_serviceMeshProfileLifeCycle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceMeshProfileDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").DoesNotExist(),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").DoesNotExist(),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceMeshProfile(data, true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceMeshProfileDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").DoesNotExist(),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").DoesNotExist(),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_advancedNetworkingIPVersionsIPv4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingConfigWithIPVersions(data, []string{"IPv4"}),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_advancedNetworkingAzure(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_advancedNetworkingAzureComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingCompleteConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("network_profile.0.docker_bridge_cidr"),
	})
}

func TestAccKubernetesAutomaticCluster_advancedNetworkingAzureWithoutDockerBridgeCidr(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingCompleteConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_advancedNetworkingAzureCiliumPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingWithCiliumPolicyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_advancedNetworkingAzureCiliumPolicyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingWithOverlayConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.advancedNetworkingWithCiliumPolicyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_enableNodePublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enableNodePublicIPConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_public_ip_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_internalNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.internalNetworkConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.max_pods").HasValue("60"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_nodePublicIPPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodePublicIPPrefixConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_node_pool.0.node_public_ip_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("default_node_pool.0.node_public_ip_prefix_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_natGatewayProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.natGatewayProfileConfig(data, 3, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.nat_gateway_profile.0.managed_outbound_ip_count").HasValue("3"),
				check.That(data.ResourceName).Key("network_profile.0.nat_gateway_profile.0.effective_outbound_ips.#").HasValue("3"),
				check.That(data.ResourceName).Key("network_profile.0.nat_gateway_profile.0.idle_timeout_in_minutes").HasValue("10"),
			),
		},
		data.ImportStep(),

		{
			Config: r.natGatewayProfileConfig(data, 4, 5),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.nat_gateway_profile.0.managed_outbound_ip_count").HasValue("4"),
				check.That(data.ResourceName).Key("network_profile.0.nat_gateway_profile.0.idle_timeout_in_minutes").HasValue("5"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_managedNatGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedNatGatewayConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.nat_gateway_profile.0.idle_timeout_in_minutes").HasValue("4"),
				check.That(data.ResourceName).Key("network_profile.0.nat_gateway_profile.0.managed_outbound_ip_count").HasValue("1"),
				check.That(data.ResourceName).Key("network_profile.0.nat_gateway_profile.0.effective_outbound_ips.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_userAssignedNatGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedNatGatewayConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_privateClusterOn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateClusterConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_fqdn").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_privateClusterOnWithPrivateDNSZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateClusterWithPrivateDNSZoneConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_cluster_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_privateClusterOnWithPrivateDNSZoneSubDomain(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateClusterWithPrivateDNSZoneSubDomain(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccKubernetesAutomaticCluster_privateClusterOnWithPrivateDNSZoneSystem(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateClusterWithPrivateDNSZoneSystemConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_privateClusterOff(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateClusterConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// check.That(data.ResourceName).Key("private_cluster_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_podCidrs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.podCidrs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_serviceCidrs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceCidrs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_standardLoadBalancer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardLoadBalancerConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_sku").HasValue("standard"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_standardLoadBalancerComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardLoadBalancerCompleteConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_sku").HasValue("standard"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_standardLoadBalancerProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardLoadBalancerProfileConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_sku").HasValue("standard"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.managed_outbound_ip_count").HasValue("3"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.effective_outbound_ips.#").HasValue("3"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.idle_timeout_in_minutes").HasValue("30"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.outbound_ports_allocated").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_standardLoadBalancerProfileComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardLoadBalancerProfileCompleteConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_sku").HasValue("standard"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.effective_outbound_ips.#").HasValue("1"),
			),
			PreventPostDestroyRefresh: true,
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_standardLoadBalancerProfileWithPortAndTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardLoadBalancerProfileWithPortAndTimeoutConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.outbound_ports_allocated").HasValue("8000"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.idle_timeout_in_minutes").HasValue("10"),
			),
			PreventPostDestroyRefresh: true,
		},
	})
}

func TestAccKubernetesAutomaticCluster_basicLoadBalancerProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.basicLoadBalancerProfileConfig(data),
			ExpectError: regexp.MustCompile("only load balancer SKU 'Standard' supports load balancer profiles. Provided load balancer type: basic"),
		},
	})
}

func TestAccKubernetesAutomaticCluster_changingLoadBalancerProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.changingLoadBalancerProfileConfigIPPrefix(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_sku").HasValue("standard"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.effective_outbound_ips.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.changingLoadBalancerProfileConfigManagedIPs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_sku").HasValue("standard"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.managed_outbound_ip_count").HasValue("1"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.effective_outbound_ips.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.unsetLoadBalancerProfileConfigIPIds(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.changingLoadBalancerProfileConfigIPIds(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_sku").HasValue("standard"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.outbound_ip_address_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.effective_outbound_ips.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.changingLoadBalancerProfileConfigIPPrefix(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("network_profile.0.load_balancer_profile.0.effective_outbound_ips.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.unsetPrefixedLoadBalancerProfileConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_httpProxyConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}
	noProxy := "\"localhost\", \"127.0.0.1\", \"mcr.microsoft.com\""
	newNoProxy := "\"localhost\", \"127.0.0.1\", \"mcr.microsoft.com\", \"monitor.azure.com\""
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpProxyConfig(data, noProxy),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("http_proxy_config.0.http_proxy").IsSet(),
				check.That(data.ResourceName).Key("http_proxy_config.0.https_proxy").IsSet(),
				check.That(data.ResourceName).Key("http_proxy_config.0.no_proxy.#").HasValue("3"),
			),
			ExpectNonEmptyPlan: true,
		},
		data.ImportStep(),
		{
			Config: r.httpProxyConfig(data, newNoProxy),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("http_proxy_config.0.http_proxy").IsSet(),
				check.That(data.ResourceName).Key("http_proxy_config.0.https_proxy").IsSet(),
				check.That(data.ResourceName).Key("http_proxy_config.0.no_proxy.#").HasValue("4"),
			),
			ExpectNonEmptyPlan: true,
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_httpProxyConfigUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}
	noProxy := "\"localhost\", \"127.0.0.1\", \"mcr.microsoft.com\""
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpProxyConfig(data, noProxy),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectNonEmptyPlan: true,
		},
		data.ImportStep(),
		{
			Config: r.httpProxyConfigUpdate(data, noProxy),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectNonEmptyPlan: true,
		},
		data.ImportStep(),
		{
			Config: r.httpProxyConfig(data, noProxy),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
			ExpectNonEmptyPlan: true,
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_httpProxyConfigWithTrustedCa(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpProxyConfigWithTrustedCa(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("http_proxy_config.0.http_proxy").IsSet(),
				check.That(data.ResourceName).Key("http_proxy_config.0.https_proxy").IsSet(),
				check.That(data.ResourceName).Key("http_proxy_config.0.trusted_ca").IsSet(),
				check.That(data.ResourceName).Key("http_proxy_config.0.no_proxy.0").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_httpProxyConfigWithSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.httpProxyConfigWithSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("http_proxy_config.0.http_proxy").IsSet(),
				check.That(data.ResourceName).Key("http_proxy_config.0.https_proxy").IsSet(),
				check.That(data.ResourceName).Key("http_proxy_config.0.no_proxy.0").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_networkPluginMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkPluginMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_networkPluginModeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkPluginBase(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.networkPluginMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_networkDataPlane(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.networkDataPlane(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_clusterPoolNodePublicIPTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clusterPoolNodePublicIPTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_clusterPoolNetworkProfileComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clusterPoolNetworkProfileComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_clusterPoolNetworkProfileUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.clusterPoolNetworkProfileComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.clusterPoolNodePublicIPTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_apiServerVnetIntegration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiServerVnetIntegrationConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_apiServerVnetIntegrationManagedVnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiServerVnetIntegrationManagedVnetConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (KubernetesAutomaticClusterResource) advancedNetworkingConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  network_profile {
    outbound_type = "loadBalancer"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) advancedNetworkingBlock(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  network_profile {
    outbound_type = "loadBalancer"

    advanced_networking {
      observability_enabled = true
      security_enabled      = true
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) advancedNetworkingBlockUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  network_profile {
    outbound_type = "loadBalancer"

    advanced_networking {
      observability_enabled = true
      security_enabled      = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) advancedNetworkingBlockRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  network_profile {
    outbound_type = "loadBalancer"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) serviceMeshProfile(data acceptance.TestData, internalIngressEnabled bool, externalIngressEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.1.0.0/16", "fd00:db8:deca::/48"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24", "fd00:db8:deca:deed::/64"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24", "fd00:db8:deca:deee::/64"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_subnet.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  linux_profile {
    admin_username = "acctestuser%[1]d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type  = "loadBalancer"
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
  }

  service_mesh_profile {
    mode                             = "Istio"
    internal_ingress_gateway_enabled = %[3]t
    external_ingress_gateway_enabled = %[4]t
    revisions                        = ["asm-1-27"]
  }

}
`, data.RandomInteger, data.Locations.Primary, internalIngressEnabled, externalIngressEnabled)
}

func (KubernetesAutomaticClusterResource) serviceMeshProfileDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.1.0.0/16", "fd00:db8:deca::/48"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24", "fd00:db8:deca:deed::/64"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24", "fd00:db8:deca:deee::/64"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  linux_profile {
    admin_username = "acctestuser%[1]d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type  = "loadBalancer"
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (KubernetesAutomaticClusterResource) advancedNetworkingConfigWithIPVersions(data acceptance.TestData, ipVersions []string) string {
	temp := make([]string, 0)
	for _, v := range ipVersions {
		temp = append(temp, fmt.Sprintf(`"%s"`, v))
	}
	ipv := strings.Join(temp, ",")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.1.0.0/16", "fd00:db8:deca::/48"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24", "fd00:db8:deca:deed::/64"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24", "fd00:db8:deca:deee::/64"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  linux_profile {
    admin_username = "acctestuser%[1]d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 2
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type  = "loadBalancer"
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
    ip_versions    = [%[3]s]
  }
}
`, data.RandomInteger, data.Locations.Primary, ipv)
}

func (KubernetesAutomaticClusterResource) advancedNetworkingCompleteConfig(data acceptance.TestData) string {
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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type  = "loadBalancer"
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

// nolint unparam
func (KubernetesAutomaticClusterResource) advancedNetworkingWithPolicyConfig(data acceptance.TestData, networkPolicy string) string {
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
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}


resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type  = "loadBalancer"
    network_policy = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, networkPolicy)
}

func (KubernetesAutomaticClusterResource) advancedNetworkingWithOverlayConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name           = "default"
    node_count     = 2
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type       = "loadBalancer"
    network_plugin_mode = "overlay"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) advancedNetworkingWithCiliumPolicyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name           = "default"
    node_count     = 2
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }


  network_profile {
    outbound_type       = "loadBalancer"
    network_policy      = "cilium"
    network_plugin_mode = "overlay"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

// nolint unparam
func (KubernetesAutomaticClusterResource) advancedNetworkingWithPolicyCompleteConfig(data acceptance.TestData, networkPolicy string) string {
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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}


resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type  = "loadBalancer"
    network_policy = "%s"
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, networkPolicy)
}

// nolint unparam
func (KubernetesAutomaticClusterResource) advancedNetworkingWithPolicyNetworkMode(data acceptance.TestData, networkPolicy string, networkMode string) string {
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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_mode   = "%s"
    network_policy = "%s"
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, networkMode, networkPolicy)
}

func (KubernetesAutomaticClusterResource) enableNodePublicIPConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name                   = "default"
    node_count             = 1
    node_public_ip_enabled = true
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) internalNetworkConfig(data acceptance.TestData) string {
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
  address_prefixes     = ["172.0.2.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.0.3.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    max_pods       = 60
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type = "loadBalancer"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) nodePublicIPPrefixConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpipprefix%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  prefix_length       = 31
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name                     = "default"
    node_count               = 1
    node_public_ip_enabled   = true
    node_public_ip_prefix_id = azurerm_public_ip_prefix.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  network_profile {
    outbound_type = "managedNATGateway"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) managedNatGatewayConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 2
    max_pods   = 60
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin_mode = "overlay"
    load_balancer_sku   = "standard"
    pod_cidr            = "10.244.0.0/16"
    service_cidr        = "10.0.0.0/16"
    dns_service_ip      = "10.0.0.10"
    outbound_type       = "managedNATGateway"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) natGatewayProfileConfig(data acceptance.TestData, ipCount int, idleTimeOut int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 2
    max_pods   = 60
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin_mode = "overlay"
    load_balancer_sku   = "standard"
    pod_cidr            = "10.244.0.0/16"
    service_cidr        = "10.0.0.0/16"
    dns_service_ip      = "10.0.0.10"
    outbound_type       = "managedNATGateway"
    nat_gateway_profile {
      managed_outbound_ip_count = %d
      idle_timeout_in_minutes   = %d
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, ipCount, idleTimeOut)
}

func (KubernetesAutomaticClusterResource) userAssignedNatGatewayConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["172.16.0.0/20"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-PIP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_nat_gateway_public_ip_association" "test" {
  nat_gateway_id       = azurerm_nat_gateway.test.id
  public_ip_address_id = azurerm_public_ip.test.id
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.16.0.0/22"]
}

resource "azurerm_subnet_nat_gateway_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  nat_gateway_id = azurerm_nat_gateway.test.id
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.16.4.0/22"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}


resource "azurerm_kubernetes_automatic_cluster" "test" {
  depends_on          = [azurerm_nat_gateway_public_ip_association.test]
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name           = "default"
    node_count     = 1
    max_pods       = 60
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    network_plugin_mode = "overlay"
    load_balancer_sku   = "standard"
    pod_cidr            = "10.244.0.0/16"
    service_cidr        = "10.0.0.0/16"
    dns_service_ip      = "10.0.0.10"
    outbound_type       = "userAssignedNATGateway"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) privateClusterConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) privateClusterWithPrivateDNSZoneConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "privatelink.%s.azmk8s.io"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestRG-aks-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_private_dns_zone.test.id
  role_definition_name = "Private DNS Zone Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "vnet" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  private_dns_zone_id = azurerm_private_dns_zone.test.id

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 1
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
  }

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) privateClusterWithPrivateDNSZoneSubDomain(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "privatelink.%s.azmk8s.io"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestRG-aks-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_private_dns_zone.test.id
  role_definition_name = "Private DNS Zone Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}


resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test1" {
  scope                = azurerm_subnet.test1.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}


resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                       = "acctestaks%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  private_dns_zone_id        = azurerm_private_dns_zone.test.id
  dns_prefix_private_cluster = "prefix"

  linux_profile {
    admin_username = "acctestuser%d"
    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    load_balancer_sku = "standard"
    outbound_type     = "loadBalancer"
  }

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) privateClusterWithPrivateDNSZoneSystemConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  private_dns_zone_id = "System"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) podCidrs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  network_profile {
    network_plugin_mode = "overlay"
    pod_cidrs           = ["10.1.1.0/24"]
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) serviceCidrs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  network_profile {
    network_plugin_mode = "overlay"
    dns_service_ip      = "10.1.1.10"
    service_cidrs       = ["10.1.1.0/24"]
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) standardLoadBalancerConfig(data acceptance.TestData) string {
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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) standardLoadBalancerCompleteConfig(data acceptance.TestData) string {
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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    dns_service_ip    = "10.10.0.10"
    service_cidr      = "10.10.0.0/16"
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) standardLoadBalancerProfileConfig(data acceptance.TestData) string {
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
  address_prefixes     = ["10.1.0.0/24"]
}
resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
    load_balancer_profile {
      managed_outbound_ip_count = 3
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) standardLoadBalancerProfileCompleteConfig(data acceptance.TestData) string {
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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}


resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    load_balancer_sku = "standard"
    outbound_type     = "loadBalancer"
    load_balancer_profile {
      outbound_ip_address_ids = [azurerm_public_ip.test.id]
      backend_pool_type       = "NodeIP"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) standardLoadBalancerProfileWithPortAndTimeoutConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azuread_application" "test" {
  display_name = "acctestspa-%[1]d"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"
  kubernetes_version  = "%[3]s"

  linux_profile {
    admin_username = "acctestuser%[1]d"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name           = "default"
    node_count     = 1
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
    load_balancer_profile {
      managed_outbound_ip_count = 2
      outbound_ports_allocated  = 8000
      idle_timeout_in_minutes   = 10
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, currentKubernetesAutomaticVersion)
}

func (KubernetesAutomaticClusterResource) basicLoadBalancerProfileConfig(data acceptance.TestData) string {
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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "basic"
    load_balancer_profile {
      managed_outbound_ip_count = 3
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) unsetPrefixedLoadBalancerProfileConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_public_ip_prefix" "test" {
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  name                = "acctestipprefix%d"
  prefix_length       = 31
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
    load_balancer_profile {
      outbound_ip_prefix_ids = []
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) changingLoadBalancerProfileConfigIPPrefix(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
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

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
    load_balancer_profile {
      outbound_ip_prefix_ids = [azurerm_public_ip_prefix.test.id]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) changingLoadBalancerProfileConfigManagedIPs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
  address_prefixes     = ["10.1.0.0/24"]
}
resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
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

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }
  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
    load_balancer_profile {
      managed_outbound_ip_count = "1"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) changingLoadBalancerProfileConfigIPIds(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
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

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }


  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }


  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
    load_balancer_profile {
      outbound_ip_address_ids = [azurerm_public_ip.test.id]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) unsetLoadBalancerProfileConfigIPIds(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
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

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
    load_balancer_profile {
      outbound_ip_address_ids = []
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) httpProxyConfigUpdate(data acceptance.TestData, noProxy string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_public_ip" "test_proxy" {
  name                = "acceptanceTestPublicIp1"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test_proxy2" {
  name                = "acceptanceTestPublicIp3"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test_aks" {
  name                = "acceptanceTestPublicIp2"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "AllowProxyAccessOn8888"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "8888"
    source_address_prefix      = "${azurerm_public_ip.test_aks.ip_address}/32"
    destination_address_prefix = "*"
  }
}

resource "azurerm_network_interface" "test" {
  name                = "test-nic%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test_proxy2.id
  }
}

resource "azurerm_network_interface_security_group_association" "test" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

locals {
  custom_data = <<CUSTOM_DATA
  #!/bin/sh
  echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections
  sudo apt-get update
  sudo apt-get install tinyproxy -y
  sudo echo "Allow ${azurerm_public_ip.test_aks.ip_address}/32" >> /etc/tinyproxy/tinyproxy.conf
  systemctl restart tinyproxy
  CUSTOM_DATA
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "vm-test-proxy%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_B1s"
  admin_username                  = "adminuser"
  admin_password                  = "P@ssW0RD1234"
  custom_data                     = base64encode(local.custom_data)
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-focal"
    sku       = "20_04-lts-gen2"
    version   = "latest"
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
    load_balancer_profile {
      outbound_ip_address_ids = [azurerm_public_ip.test_aks.id]
    }
  }

  http_proxy_config {
    http_proxy  = "http://${azurerm_public_ip.test_proxy2.ip_address}:8888/"
    https_proxy = "http://${azurerm_public_ip.test_proxy2.ip_address}:8888/"
    no_proxy    = [%s]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger, noProxy)
}

func (KubernetesAutomaticClusterResource) httpProxyConfig(data acceptance.TestData, noProxy string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

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
  address_prefixes     = ["10.1.0.0/24"]
}
resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_public_ip" "test_proxy" {
  name                = "acceptanceTestPublicIp1"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test_proxy2" {
  name                = "acceptanceTestPublicIp3"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test_aks" {
  name                = "acceptanceTestPublicIp2"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "AllowProxyAccessOn8888"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "8888"
    source_address_prefix      = "${azurerm_public_ip.test_aks.ip_address}/32"
    destination_address_prefix = "*"
  }
}

resource "azurerm_network_interface" "test" {
  name                = "test-nic%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test_proxy.id
  }
}

resource "azurerm_network_interface_security_group_association" "test" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

locals {
  custom_data = <<CUSTOM_DATA
  #!/bin/sh
  echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections
  sudo apt-get update
  sudo apt-get install tinyproxy -y
  sudo echo "Allow ${azurerm_public_ip.test_aks.ip_address}/32" >> /etc/tinyproxy/tinyproxy.conf
  systemctl restart tinyproxy
  CUSTOM_DATA
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "vm-test-proxy%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_B1s"
  admin_username                  = "adminuser"
  admin_password                  = "P@ssW0RD1234"
  custom_data                     = base64encode(local.custom_data)
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts-gen2"
    version   = "latest"
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    outbound_type     = "loadBalancer"
    load_balancer_sku = "standard"
    load_balancer_profile {
      outbound_ip_address_ids = [azurerm_public_ip.test_aks.id]
    }
  }

  http_proxy_config {
    http_proxy  = "http://${azurerm_public_ip.test_proxy.ip_address}:8888/"
    https_proxy = "http://${azurerm_public_ip.test_proxy.ip_address}:8888/"
    no_proxy    = [%s]
  }
}








`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger, noProxy)
}

func (KubernetesAutomaticClusterResource) httpProxyConfigWithTrustedCa(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_public_ip" "test_proxy" {
  name                = "acceptanceTestPublicIp1"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test_aks" {
  name                = "acceptanceTestPublicIp2"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "AllowProxyAccessOn8888"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "8888"
    source_address_prefix      = "${azurerm_public_ip.test_aks.ip_address}/32"
    destination_address_prefix = "*"
  }
}

resource "azurerm_network_interface" "test" {
  name                = "test-nic%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test_proxy.id
  }
}

resource "azurerm_network_interface_security_group_association" "test" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

locals {
  custom_data = <<CUSTOM_DATA
  #!/bin/sh
  echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections
  sudo apt-get update
  sudo apt-get install tinyproxy -y
  sudo echo "Allow ${azurerm_public_ip.test_aks.ip_address}/32" >> /etc/tinyproxy/tinyproxy.conf
  systemctl restart tinyproxy
  CUSTOM_DATA
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "vm-test-proxy%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_B1s"
  admin_username                  = "adminuser"
  admin_password                  = "P@ssW0RD1234"
  custom_data                     = base64encode(local.custom_data)
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts-gen2"
    version   = "latest"
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    node_count = 2
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    load_balancer_sku = "standard"
    load_balancer_profile {
      outbound_ip_address_ids = [azurerm_public_ip.test_aks.id]
    }
  }

  http_proxy_config {
    http_proxy  = "http://${azurerm_public_ip.test_proxy.ip_address}:8888/"
    https_proxy = "http://${azurerm_public_ip.test_proxy.ip_address}:8888/"
    no_proxy = [
      "localhost",
      "127.0.0.1",
      "mcr.microsoft.com"
    ]
    trusted_ca = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ6RENDQVMyZ0F3SUJBZ0lCQVRBS0JnZ3Foa2pPUFFRREJEQVFNUTR3REFZRFZRUUtFd1ZGVGtOUFRUQWUKRncwd09URXhNVEF5TXpBd01EQmFGdzB4TURBMU1Ea3lNekF3TURCYU1CQXhEakFNQmdOVkJBb1RCVVZPUTA5TgpNSUdiTUJBR0J5cUdTTTQ5QWdFR0JTdUJCQUFqQTRHR0FBUUJBcUN1Um94NU4zTVRVOHdUdUllSUJYRjdpTW5oCm50cW1HVktRMGhmUUZEUUd2K0x5ZHVvN0pQcUZwL1kyamxYU2ROckFkejVXeGJyWStrRHhJcGtCUXRJQWtJREQKWlZtVHVlcTNaREFmY0dkRU5uek5KVkNhUGxIWEpMdkVFSU5jb0prVU8rK2NWeXl3ZHJlVkpjNjd2aE54MVRkWApWM3BwN2YrUmJPbU5LYm5WUkJ5ak5UQXpNQTRHQTFVZER3RUIvd1FFQXdJSGdEQVRCZ05WSFNVRUREQUtCZ2dyCkJnRUZCUWNEQVRBTUJnTlZIUk1CQWY4RUFqQUFNQW9HQ0NxR1NNNDlCQU1FQTRHTUFEQ0JpQUpDQWJiYjdzdkkKNXR1aEN5QTNqUVRTZ0E4enB2azBZV05Ya1owN3h6ZFY4amRNTXVtQ2FXOXljRUlxSjVLU3F1dVBoVXc5b2VregpCNTFkYXliVjFWUVhWVmRWQWtJQStrTU1TSnp3dHpIcU5BVVRtaVpQY2c3SDh2MUFTbDR0UjZscEtUcFVQWTJYCmxYT0N0MllmNGRzRnNpanV2emJKQmR4NzVkNEVmNVRSSFBjZytQSE5aZ2c9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
  }

  lifecycle {
    ignore_changes = [
      http_proxy_config.0.no_proxy
    ]
  }
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) httpProxyConfigWithSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "test-nic%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "vm-test-proxy%d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_B1s"
  admin_username                  = "adminuser"
  admin_password                  = "P@ssW0RD1234"
  custom_data                     = base64encode("#!/bin/bash\nsudo apt-get update\nsudo apt-get install tinyproxy -y\nsudo echo \"Allow 10.0.0.0/8\" \u003e\u003e /etc/tinyproxy/tinyproxy.conf\nsystemctl restart tinyproxy")
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts-gen2"
    version   = "latest"
  }
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
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
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  http_proxy_config {
    http_proxy  = "http://${azurerm_network_interface.test.private_ip_address}:8888/"
    https_proxy = "http://${azurerm_network_interface.test.private_ip_address}:8888/"
    no_proxy = [
      "10.1.0.0/24"
    ]
  }

  lifecycle {
    ignore_changes = [
      http_proxy_config.0.no_proxy
    ]
  }
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, currentKubernetesAutomaticVersion, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) networkDataPlane(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestRG-vnet-%[2]d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestRG-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.10.0.0/16"]

}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.11.2.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}


resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name           = "default"
    node_count     = 1
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }

  network_profile {
    pod_cidr            = "192.168.0.0/16"
    network_plugin_mode = "overlay"
    outbound_type       = "loadBalancer"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) networkPluginBase(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestRG-vnet-%[2]d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestRG-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.10.0.0/16"]

}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.11.0.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name           = "default"
    node_count     = 1
    vm_size        = "Standard_DS3_v2"
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }
  network_profile {
    outbound_type = "loadBalancer"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) networkPluginMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestRG-vnet-%[2]d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestRG-subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.10.0.0/16"]

}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.11.0.0/24"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_role_assignment" "network" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name           = "default"
    node_count     = 1
    vm_size        = "Standard_DS3_v2"
    vnet_subnet_id = azurerm_subnet.test.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test1.id
  }
  network_profile {
    pod_cidr            = "192.168.0.0/16"
    outbound_type       = "loadBalancer"
    network_plugin_mode = "overlay"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) clusterPoolNodePublicIPTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctestasg-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name                   = "default"
    node_count             = 1
    node_public_ip_enabled = true
    node_network_profile {
      node_public_ip_tags = {
        RoutingPreference = "Internet"
      }
    }
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) clusterPoolNetworkProfileComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_application_security_group" "test" {
  name                = "acctestasg-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  default_node_pool {
    name                   = "default"
    node_count             = 1
    node_public_ip_enabled = true
    node_network_profile {
      allowed_host_ports {
        port_start = 8001
        port_end   = 8002
        protocol   = "UDP"
      }
      application_security_group_ids = [azurerm_application_security_group.test.id]
      node_public_ip_tags = {
        RoutingPreference = "Internet"
      }
    }
    upgrade_settings {
      max_surge = "10%%"
    }
  }
  identity {
    type = "SystemAssigned"
  }
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) apiServerVnetIntegrationConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.0.0.0/8"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/16"]

  delegation {
    name = "aks-delegation"

    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.ContainerService/managedClusters"
    }
  }
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.2.0.0/16"]
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestRG-aks-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_subnet.test.id
  role_definition_name = "Network Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name           = "default"
    node_count     = 1
    vnet_subnet_id = azurerm_subnet.test1.id
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  network_profile {
    load_balancer_sku = "standard"
    outbound_type     = "loadBalancer"
  }

  api_server_access_profile {
    subnet_id = azurerm_subnet.test.id
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) apiServerVnetIntegrationManagedVnetConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    load_balancer_sku = "standard"
    outbound_type     = "managedNATGateway"
  }

  api_server_access_profile {
  }
}
`, data.Locations.Primary, data.RandomInteger)
}
