// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KubernetesAutomaticClusterDataSource struct{}

func TestAccDataSourceKubernetesAutomaticCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				// check.That(data.ResourceName).Key("kube_config.0.client_key").Exists(),
				// check.That(data.ResourceName).Key("kube_config.0.client_certificate").Exists(),
				// check.That(data.ResourceName).Key("kube_config.0.cluster_ca_certificate").Exists(),
				// check.That(data.ResourceName).Key("kube_config.0.host").Exists(),
				// check.That(data.ResourceName).Key("kube_config.0.username").Exists(),
				// check.That(data.ResourceName).Key("kube_config.0.password").Exists(),
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

func TestAccDataSourceKubernetesAutomaticCluster_privateCluster(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: KubernetesAutomaticClusterResource{}.privateClusterConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("private_fqdn").Exists(),
			),
		},
		data.ImportStep("service_principal.0.client_secret"),
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_roleBasedAccessControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.roleBasedAccessControlConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("role_based_access_control_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("azure_active_directory_role_based_access_control.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_admin_config.#").HasValue("0"),
				check.That(data.ResourceName).Key("kube_admin_config_raw").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_roleBasedAccessControlAAD_OlderKubernetesVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.roleBasedAccessControlAADManagedConfigOlderKubernetesVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kube_config.#").HasValue("1"),
				check.That(data.ResourceName).Key("kube_config.0.host").IsSet(),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_internalNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.internalNetworkConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_advancedNetworkingAzure(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingAzureConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_advancedNetworkingAzureComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.advancedNetworkingAzureCompleteConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.vnet_subnet_id").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.network_plugin").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.dns_service_ip").Exists(),
				check.That(data.ResourceName).Key("network_profile.0.service_cidr").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_addOnProfileOMS(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

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

func TestAccDataSourceKubernetesAutomaticCluster_addOnProfileIngressApplicationGatewayAppGateway(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

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

func TestAccDataSourceKubernetesAutomaticCluster_addOnProfileIngressApplicationGatewaySubnetCIDR(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

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

func TestAccDataSourceKubernetesAutomaticCluster_addOnProfileIngressApplicationGatewaySubnetId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

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

func TestAccDataSourceKubernetesAutomaticCluster_addOnProfileOpenServiceMesh(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileOpenServiceMeshConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("open_service_mesh_enabled").HasValue("true"),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_addOnProfileAzureKeyvaultSecretsProvider(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.addOnProfileAzureKeyvaultSecretsProviderConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("key_vault_secrets_provider.0.secret_rotation_interval").HasValue("2m"),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_nodeLabels(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}
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

func TestAccDataSourceKubernetesAutomaticCluster_nodePublicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.nodePublicIPConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("agent_pool_profile.0.node_public_ip_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("agent_pool_profile.0.node_public_ip_prefix_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_microsoftDefender(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.microsoftDefender(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("microsoft_defender.0.log_analytics_workspace_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_serviceMesh(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

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

func TestAccDataSourceKubernetesAutomaticCluster_serviceMeshCertificateAuthority(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.serviceMeshCertificateAuthority(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.certificate_authority.0.key_vault_id").Exists(),
				check.That(data.ResourceName).Key("service_mesh_profile.0.certificate_authority.0.root_cert_object_name").Exists(),
				check.That(data.ResourceName).Key("service_mesh_profile.0.certificate_authority.0.cert_chain_object_name").Exists(),
				check.That(data.ResourceName).Key("service_mesh_profile.0.certificate_authority.0.cert_object_name").Exists(),
				check.That(data.ResourceName).Key("service_mesh_profile.0.certificate_authority.0.key_object_name").Exists(),
			),
		},
	})
}

func TestAccDataSourceKubernetesAutomaticCluster_serviceMeshRevisions(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			// create a cluster with an istio revision with revision currently at asm-1-27
			Config: r.serviceMeshRevisions(data, `["asm-1-27"]`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.revisions.0").HasValue("asm-1-27"),
			),
		},
		{
			// start istio revision canary upgrade to asm-1-28
			Config: r.serviceMeshRevisions(data, `["asm-1-27", "asm-1-28"]`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.revisions.0").HasValue("asm-1-27"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.revisions.1").HasValue("asm-1-28"),
			),
		},
		{
			// rollback the istio revision back to asm-1-27
			Config: r.serviceMeshRevisions(data, `["asm-1-27"]`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.revisions.0").HasValue("asm-1-27"),
			),
		},
		{
			// start istio revision canary upgrade to asm-1-28
			Config: r.serviceMeshRevisions(data, `["asm-1-27", "asm-1-28"]`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.revisions.0").HasValue("asm-1-27"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.revisions.1").HasValue("asm-1-28"),
			),
		},
		{
			// complete the istio revision upgrade to asm-1-27
			Config: r.serviceMeshRevisions(data, `["asm-1-27"]`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("service_mesh_profile.0.mode").HasValue("Istio"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.internal_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.external_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh_profile.0.revisions.0").HasValue("asm-1-27"),
			),
		},
	})
}

func (KubernetesAutomaticClusterDataSource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.basicVMSSConfig(data))
}

func (KubernetesAutomaticClusterDataSource) roleBasedAccessControlConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.roleBasedAccessControlConfig(data))
}

func (KubernetesAutomaticClusterDataSource) roleBasedAccessControlAADManagedConfigOlderKubernetesVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.roleBasedAccessControlAADManagedConfigOlderKubernetesVersion(data, ""))
}

func (KubernetesAutomaticClusterDataSource) internalNetworkConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.internalNetworkConfig(data))
}

func (KubernetesAutomaticClusterDataSource) advancedNetworkingAzureConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.advancedNetworkingConfig(data))
}

func (KubernetesAutomaticClusterDataSource) advancedNetworkingAzureCompleteConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.advancedNetworkingCompleteConfig(data))
}

func (KubernetesAutomaticClusterDataSource) addOnProfileOMSConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.addonProfileOMSConfig(data))
}

func (KubernetesAutomaticClusterDataSource) addOnProfileIngressApplicationGatewayAppGatewayConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.addonProfileIngressApplicationGatewayAppGatewayConfig(data))
}

func (KubernetesAutomaticClusterDataSource) addOnProfileIngressApplicationGatewaySubnetCIDRConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.addonProfileIngressApplicationGatewaySubnetCIDRConfig(data))
}

func (KubernetesAutomaticClusterDataSource) addOnProfileIngressApplicationGatewaySubnetIdConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.addonProfileIngressApplicationGatewaySubnetIdConfig(data))
}

func (KubernetesAutomaticClusterDataSource) addOnProfileOpenServiceMeshConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.addonProfileOpenServiceMeshConfig(data, true))
}

func (KubernetesAutomaticClusterDataSource) addOnProfileAzureKeyvaultSecretsProviderConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.addonProfileAzureKeyVaultSecretsProviderConfig(data, "2m"))
}

func (KubernetesAutomaticClusterDataSource) nodeLabelsConfig(data acceptance.TestData, labels map[string]string) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.nodeLabelsConfig(data, labels))
}

func (KubernetesAutomaticClusterDataSource) nodePublicIPConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.nodePublicIPPrefixConfig(data))
}

func (KubernetesAutomaticClusterDataSource) microsoftDefender(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.microsoftDefender(data))
}

func (KubernetesAutomaticClusterDataSource) serviceMesh(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.serviceMeshProfile(data, true, true))
}

func (s KubernetesAutomaticClusterDataSource) serviceMeshCertificateAuthority(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.addonProfileServiceMeshProfileCertificateAuthorityConfig(data))
}

func (KubernetesAutomaticClusterDataSource) serviceMeshRevisions(data acceptance.TestData, revisions string) string {
	return fmt.Sprintf(`
%s
data "azurerm_kubernetes_automatic_cluster" "test" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
}
`, KubernetesAutomaticClusterResource{}.addonProfileServiceMeshProfileRevisionsConfig(data, revisions))
}
