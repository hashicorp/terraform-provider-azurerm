// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KubernetesClusterDataSource struct{}

func TestAccDataSourceKubernetesCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kube_config.0.client_key").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.client_certificate").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.cluster_ca_certificate").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.host").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.username").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.password").Exists(),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").HasValue(""),
				check.That(data.ResourceName).Key("kubelet_identity.0.object_id").Exists(),
				check.That(data.ResourceName).Key("kubelet_identity.0.client_id").Exists(),
				check.That(data.ResourceName).Key("kubelet_identity.0.user_assigned_identity_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_privateCluster(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: KubernetesClusterResource{}.privateClusterConfig(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("private_fqdn").Exists(),
				check.That(data.ResourceName).Key("private_cluster_enabled").HasValue("true"),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccDataSourceKubernetesCluster_roleBasedAccessControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.roleBasedAccessControlConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("role_based_access_control_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_roleBasedAccessControlAAD_VOneDotTwoFourDotNine(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.roleBasedAccessControlAADManagedConfigVOneDotTwoFourDotNine(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kube_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_config.0.host").IsSet(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_roleBasedAccessControlAAD(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantId := os.Getenv("ARM_TENANT_ID")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.roleBasedAccessControlAADConfig(data, clientId, clientSecret, tenantId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("role_based_access_control_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.0.client_app_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.0.server_app_id").Exists(),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_localAccountDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}
	clientData := data.Client()

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.localAccountDisabled(data, clientData.TenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("role_based_access_control_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.0.managed").HasValue("true"),
				check.That(data.ResourceName).Key("kube_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_config_raw").Exists(),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_internalNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.internalNetworkConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_advancedNetworkingAzure(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingAzureConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").HasValue("azure"),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_advancedNetworkingNone(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingNoneConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").HasValue("none"),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_advancedNetworkingAzureCalicoPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingAzureCalicoPolicyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").HasValue("azure"),
				check.That(data.ResourceName).Key("network_profile.0.network_policy").HasValue("calico"),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_policy").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_advancedNetworkingAzureNPMPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingAzureNPMPolicyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").HasValue("azure"),
				check.That(data.ResourceName).Key("network_profile.0.network_policy").HasValue("azure"),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_policy").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_advancedNetworkingAzureComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingAzureCompleteConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").HasValue("azure"),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_advancedNetworkingAzureCalicoPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingAzureCalicoPolicyCompleteConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").HasValue("azure"),
				check.That(data.ResourceName).Key("network_profile.0.network_policy").HasValue("calico"),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_policy").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_advancedNetworkingAzureNPMPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingAzureNPMPolicyCompleteConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").HasValue("azure"),
				check.That(data.ResourceName).Key("network_profile.0.network_policy").HasValue("azure"),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_policy").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_advancedNetworkingKubenet(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingKubenetConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").HasValue("kubenet"),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_advancedNetworkingKubenetComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingKubenetCompleteConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").HasValue("kubenet"),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_addOnProfileOMS(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileOMSConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("oms_agent.#").HasValue("1"),
				check.That(data.ResourceName).Key("oms_agent.0.log_analytics_workspace_id").Exists(),
				check.That(data.ResourceName).Key("oms_agent.0.oms_agent_identity.0.client_id").Exists(),
				check.That(data.ResourceName).Key("oms_agent.0.oms_agent_identity.0.object_id").Exists(),
				check.That(data.ResourceName).Key("oms_agent.0.oms_agent_identity.0.user_assigned_identity_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_addOnProfileAzurePolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileAzurePolicyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("azure_policy_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_addOnProfileRouting(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileRoutingConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("http_application_routing_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("http_application_routing_zone_name").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_addOnProfileIngressApplicationGatewayAppGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileIngressApplicationGatewayAppGatewayConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("ingress_application_gateway.#").HasValue("1"),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.effective_gateway_id").MatchesOtherKey(
					check.That(data.ResourceName).Key("ingress_application_gateway.0.gateway_id"),
				),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.subnet_cidr").IsEmpty(),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.subnet_id").IsEmpty(),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.ingress_application_gateway_identity.0.client_id").Exists(),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.ingress_application_gateway_identity.0.object_id").Exists(),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.ingress_application_gateway_identity.0.user_assigned_identity_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_addOnProfileIngressApplicationGatewaySubnetCIDR(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileIngressApplicationGatewaySubnetCIDRConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("ingress_application_gateway.#").HasValue("1"),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.gateway_id").IsEmpty(),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.subnet_cidr").HasValue(addOnAppGatewaySubnetCIDR),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.subnet_id").IsEmpty(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_addOnProfileIngressApplicationGatewaySubnetId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileIngressApplicationGatewaySubnetIdConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("ingress_application_gateway.#").HasValue("1"),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.gateway_id").IsEmpty(),
				check.That(data.ResourceName).Key("ingress_application_gateway.0.subnet_cidr").IsEmpty(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_addOnProfileOpenServiceMesh(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileOpenServiceMeshConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("open_service_mesh_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_addOnProfileAzureKeyvaultSecretsProvider(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileAzureKeyvaultSecretsProviderConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_vault_secrets_provider.0.secret_rotation_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("key_vault_secrets_provider.0.secret_rotation_interval").HasValue("2m"),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_autoscalingNoAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.autoScalingNoAvailabilityZonesConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.min_count").HasValue("1"),
				check.That(data.ResourceName).Key("agent_pool_profile.0.max_count").HasValue("2"),
				check.That(data.ResourceName).Key("agent_pool_profile.0.type").HasValue("VirtualMachineScaleSets"),
				check.That(data.ResourceName).Key("agent_pool_profile.0.enable_auto_scaling").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_autoscalingWithAvailabilityZones(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.autoScalingWithAvailabilityZonesConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.min_count").HasValue("1"),
				check.That(data.ResourceName).Key("agent_pool_profile.0.max_count").HasValue("2"),
				check.That(data.ResourceName).Key("agent_pool_profile.0.type").HasValue("VirtualMachineScaleSets"),
				check.That(data.ResourceName).Key("agent_pool_profile.0.enable_auto_scaling").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_nodeLabels(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}
	labels := map[string]string{"key": "value"}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.nodeLabelsConfig(data, labels),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.node_labels.key").HasValue("value"),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_nodePublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.nodePublicIPConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.enable_node_public_ip").HasValue("true"),
				check.That(data.ResourceName).Key("agent_pool_profile.0.node_public_ip_prefix_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_oidcIssuerEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.oidcIssuer(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("oidc_issuer_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("oidc_issuer_url").IsSet(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_oidcIssuerDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.oidcIssuer(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("oidc_issuer_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("oidc_issuer_url").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_microsoftDefender(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.microsoftDefender(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("microsoft_defender.0.log_analytics_workspace_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_customCaTrustCerts(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	fakeCertList := []string{
		"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURjRENDQWxpZ0F3SUJBZ0lFU1QwSUhEQU5CZ2txaGtpRzl3MEJBUXNGQURCUk1Rc3dDUVlEVlFRR0V3SlEKVERFTk1Bc0dBMVVFQXd3RVZHVnpkREVWTUJNR0ExVUVCd3dNUkdWbVlYVnNkQ0JEYVhSNU1Sd3dHZ1lEVlFRSwpEQk5FWldaaGRXeDBJRU52YlhCaGJua2dUSFJrTUI0WERUSXpNRFV5T0RFeE1qY3dNMW9YRFRNek1EVXlOVEV4Ck1qY3dNMW93VVRFTE1Ba0dBMVVFQmhNQ1VFd3hEVEFMQmdOVkJBTU1CRlJsYzNReEZUQVRCZ05WQkFjTURFUmwKWm1GMWJIUWdRMmwwZVRFY01Cb0dBMVVFQ2d3VFJHVm1ZWFZzZENCRGIyMXdZVzU1SUV4MFpEQ0NBU0l3RFFZSgpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFLN2JIYWtxSkdRMWVBOUFHUmlhNGl2anNDRXlGMDhDCjNpSzJZeWthNkREeldmTk1tRWpOUjJiQVZOMEhlLy9pWTd1VjJ2dXl6V1UxMzZGVkdMZkdyeTZGOHNQQUZaSzYKSE4vcWk1QVp6MUpoOGdWSTRwS1pjZEFxQS81clF3VVlvWVN3Q245dGVOYytsbU1ZUk5OcTVwdlV2NjcrNEM3MgpPc3BOSUxSclhBbWNUb1YveVRZVzFKWDBOeEJJSHZZaFZXUE9LQXpRZDQ5UEpSeFpqMUgydCszMEFsazgzTDFwClFzTGx2SzV3MjJpeXdkYVpRN1lmV0xXd1hPQzVPWXdRTUw1R3BHUFNQaEdxdjhqSUhpcHBVeTdrRDlNWFFZOFoKdDl2QkczMzVWSEdlUjI2QnNQQXRFbTJjR05ocjA5cmRvdWJGd2tDR05OYXNVamFoVW9CKzhPY0NBd0VBQWFOUQpNRTR3SFFZRFZSME9CQllFRk9CNmNpTGtUL21Cc2xXSm5Na2phQzZqbjd4ek1COEdBMVVkSXdRWU1CYUFGT0I2CmNpTGtUL21Cc2xXSm5Na2phQzZqbjd4ek1Bd0dBMVVkRXdRRk1BTUJBZjh3RFFZSktvWklodmNOQVFFTEJRQUQKZ2dFQkFKTklHdHJpeFlCRUc1Yy9iQWdOMHlMOEJvOW9nN29ha0hVMUc5TjBxOUNWWXhjOVhma2ZUaEhYOVBUeApMbVNGcHJEQlAyYnVGTzVIUDFpbnNFT1E2N1lGanAvRjVJWGdaQ2twZUpGdDBTL0R3N2ZRbFJJN2RCNGQzNmIzCmE1R2txU0M4aFlZemxLUm9DRGNhalp4QmdoVUFxK0tnTnV4RmNsM1Fnd1Uyam1QbkU4a1A4TmgyM3hlVUJ3WEkKL3pqbU1rdjV4SFhKdHBpdlpzTlpSSUttQW56RU9TWGlRK2JMTStTdlhtSkhYd29YYTZyTXg4YmkySzV4WkhIRwpkUHA1TnQ3L2dxOUdXcm95SkVjSFpEclBiSnR2WGFibTZYUXpxTTFYUzA3SDlaSFBXc0dENGlBM1k0T3JUUlRCClZ5blRPUDl5U3cwbklaVEk4YjZuR2RHTzBOOD0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==",
		"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURmakNDQW1hZ0F3SUJBZ0lFZnlWdk56QU5CZ2txaGtpRzl3MEJBUXNGQURCWU1Rc3dDUVlEVlFRR0V3SlEKVERFVU1CSUdBMVVFQXd3TFJtRnJaU0JEWlhKMElESXhGVEFUQmdOVkJBY01ERVJsWm1GMWJIUWdRMmwwZVRFYwpNQm9HQTFVRUNnd1RSR1ZtWVhWc2RDQkRiMjF3WVc1NUlFeDBaREFlRncweU16QTJNRFF3TnpJME1qZGFGdzB5Ck5UQTJNRE13TnpJME1qZGFNRmd4Q3pBSkJnTlZCQVlUQWxCTU1SUXdFZ1lEVlFRRERBdEdZV3RsSUVObGNuUWcKTWpFVk1CTUdBMVVFQnd3TVJHVm1ZWFZzZENCRGFYUjVNUnd3R2dZRFZRUUtEQk5FWldaaGRXeDBJRU52YlhCaApibmtnVEhSa01JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBMENTdVdUaGNjSG5MCkhFdjk4SUVNc2JLY3h4YVh4YTZiRXl1Yy9sUjRackpVN2p6eVlWNGVscTV5WTgwdDFCM0MyV3E2SXFoajErSGYKYW0xaStsU1FTejM1eWNnTWlwSWp2cUxKOVIzMVF0Wi9TRURkdGV2b2JqbytEa1dCOE55cG9Ia0pVbEIyQnR6ZgpOK09KeVFSdXU1b1cya2c5OE5Bd3JuTGpmQ0lremVWcFh5d0l4Tkx2ZmFrVGxpNWpYdG9WWG5pOTU5bmtINWVwClkrRnVoSEQwaU5CS25XYVkxR2QwVGhhSHNwTERmNFUycmo2WE5SZHd6QVZoVkdhUm02cndvSHRZeDVrYys1ZWMKQ0F4UEdRWFRzTzJUTHVrQzJ2YXI0M3RUM0ZjSC9taDRST2JaaThZS2xSQ3Fldm1QU1RmZ293RUFkTjlvSmxyRApXN2lzN2NnQjhRSURBUUFCbzFBd1RqQWRCZ05WSFE0RUZnUVVuRkRqN0pBQW9WZ2NzQkgyNzdMOHZlM0Q4U293Ckh3WURWUjBqQkJnd0ZvQVVuRkRqN0pBQW9WZ2NzQkgyNzdMOHZlM0Q4U293REFZRFZSMFRCQVV3QXdFQi96QU4KQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBT0diT0Zyek4rN2YxbzhJSDNtMXZxT3IyTUtvNEZMWExGRjBVbEhkNApwZXRhL05aQjArUmQ3TnUrOCtnUnlUbEJWZU9EZjN5SXU0TlFCUU92MlNqdS9Jakd0MUtmaUF3WkUwT1RUQXc3CnhIWStsMVBJWEFFVWNqNk00cjFKQzc4ZVZrc2pycTZoV1RPZ0RrSVZuRjY3bXlReXduR25EY1k0d0Fqc2pUajgKKzR4NTIrRi9QaVNQVGtjUFNuN0s2UjQzaEt5QUs2Z0poOHE5cVNhME5RQ2U2czhwTGU2SVY5SElWVVFFVERVOQpsM1VWWHNBMGx4dlB0blU1TXo2QWQ5cDA5L2w4d3o0cUdBdGFCUEd3K0R2cTNlaHdTd2VZZ3VHSktDQjhjb01JCjJRVUo0Zi9mNkFNVWtMeWxYZ3RSUEt1QjA3d3YwTmk1eWI5MjlFY1FJQ0l2dFE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t",
	}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.customCaTrustCertificates(data, fakeCertList),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("custom_ca_trust_certificates_base64.0").Exists(),
				check.That(data.ResourceName).Key("custom_ca_trust_certificates_base64.1").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesCluster_serviceMesh(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_cluster", "test")
	r := KubernetesClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.serviceMesh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("true"),
			),
		},
	})
}

func (KubernetesClusterDataSource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.basicVMSSConfig(data))
}

func (KubernetesClusterDataSource) roleBasedAccessControlConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.roleBasedAccessControlConfig(data))
}

func (KubernetesClusterDataSource) roleBasedAccessControlAADManagedConfigVOneDotTwoFourDotNine(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.roleBasedAccessControlAADManagedConfigVOneDotTwoFourDotNine(data, ""))
}

func (KubernetesClusterDataSource) localAccountDisabled(data acceptance.TestData, tenantId string) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.roleBasedAccessControlAADManagedConfigWithLocalAccountDisabled(data, tenantId))
}

func (KubernetesClusterDataSource) roleBasedAccessControlAADConfig(data acceptance.TestData, clientId, clientSecret, tenantId string) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.roleBasedAccessControlAADConfig(data, clientId, clientSecret, tenantId))
}

func (KubernetesClusterDataSource) internalNetworkConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.internalNetworkConfig(data))
}

func (KubernetesClusterDataSource) advancedNetworkingAzureConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.advancedNetworkingConfig(data, "azure"))
}

func (KubernetesClusterDataSource) advancedNetworkingNoneConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.advancedNetworkingConfig(data, "none"))
}

func (KubernetesClusterDataSource) advancedNetworkingAzureCalicoPolicyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.advancedNetworkingWithPolicyConfig(data, "azure", "calico"))
}

func (KubernetesClusterDataSource) advancedNetworkingAzureNPMPolicyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.advancedNetworkingWithPolicyConfig(data, "azure", "azure"))
}

func (KubernetesClusterDataSource) advancedNetworkingAzureCompleteConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.advancedNetworkingCompleteConfig(data, "azure", true))
}

func (KubernetesClusterDataSource) advancedNetworkingAzureCalicoPolicyCompleteConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.advancedNetworkingWithPolicyCompleteConfig(data, "azure", "calico"))
}

func (KubernetesClusterDataSource) advancedNetworkingAzureNPMPolicyCompleteConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.advancedNetworkingWithPolicyCompleteConfig(data, "azure", "azure"))
}

func (KubernetesClusterDataSource) advancedNetworkingKubenetConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.advancedNetworkingConfig(data, "kubenet"))
}

func (KubernetesClusterDataSource) advancedNetworkingKubenetCompleteConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.advancedNetworkingCompleteConfig(data, "kubenet", true))
}

func (KubernetesClusterDataSource) addOnProfileOMSConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.addonProfileOMSConfig(data))
}

func (KubernetesClusterDataSource) addOnProfileAzurePolicyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.addonProfileAzurePolicyConfig(data, true))
}

func (KubernetesClusterDataSource) addOnProfileRoutingConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.addonProfileRoutingConfig(data, true))
}

func (KubernetesClusterDataSource) addOnProfileIngressApplicationGatewayAppGatewayConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.addonProfileIngressApplicationGatewayAppGatewayConfig(data))
}

func (KubernetesClusterDataSource) addOnProfileIngressApplicationGatewaySubnetCIDRConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.addonProfileIngressApplicationGatewaySubnetCIDRConfig(data))
}

func (KubernetesClusterDataSource) addOnProfileIngressApplicationGatewaySubnetIdConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.addonProfileIngressApplicationGatewaySubnetIdConfig(data))
}

func (KubernetesClusterDataSource) addOnProfileOpenServiceMeshConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.addonProfileOpenServiceMeshConfig(data, true))
}

func (KubernetesClusterDataSource) addOnProfileAzureKeyvaultSecretsProviderConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.addonProfileAzureKeyVaultSecretsProviderConfig(data, true, "2m"))
}

func (KubernetesClusterDataSource) autoScalingNoAvailabilityZonesConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.autoscaleNoAvailabilityZonesConfig(data))
}

func (KubernetesClusterDataSource) autoScalingWithAvailabilityZonesConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.autoscaleWithAvailabilityZonesConfig(data))
}

func (KubernetesClusterDataSource) nodeLabelsConfig(data acceptance.TestData, labels map[string]string) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.nodeLabelsConfig(data, labels))
}

func (KubernetesClusterDataSource) nodePublicIPConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.nodePublicIPPrefixConfig(data))
}

func (KubernetesClusterDataSource) oidcIssuer(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
%s
data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.oidcIssuer(data, enabled))
}

func (KubernetesClusterDataSource) microsoftDefender(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.microsoftDefender(data))
}

func (KubernetesClusterDataSource) customCaTrustCertificates(data acceptance.TestData, fakeCertsList []string) string {
	return fmt.Sprintf(`
%s
data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.customCATrustCertificates(data, fakeCertsList))
}

func (KubernetesClusterDataSource) serviceMesh(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_kubernetes_cluster" "test" {
  name                = azurerm_kubernetes_cluster.test.name
  resource_group_name = azurerm_kubernetes_cluster.test.resource_group_name
}
`, KubernetesClusterResource{}.serviceMeshProfile(data, true, true))
}
