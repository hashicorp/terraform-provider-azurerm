// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccKubernetesAutomaticCluster_serviceMeshProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceMeshProfile(data, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh.#").HasValue("1"),
				check.That(data.ResourceName).Key("service_mesh.0.internal_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh.0.external_ingress_gateway_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceMeshProfile(data, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh.#").HasValue("1"),
				check.That(data.ResourceName).Key("service_mesh.0.internal_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh.0.external_ingress_gateway_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_serviceMeshIngressConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceMeshIngressConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh.#").HasValue("1"),
				check.That(data.ResourceName).Key("service_mesh.0.internal_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh.0.external_ingress_gateway_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceMeshProfile(data, false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh.#").HasValue("1"),
				check.That(data.ResourceName).Key("service_mesh.0.internal_ingress_gateway_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("service_mesh.0.external_ingress_gateway_enabled").HasValue("false"),
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
				check.That(data.ResourceName).Key("service_mesh.#").HasValue("0"),
				check.That(data.ResourceName).Key("service_mesh.0.internal_ingress_gateway_enabled").DoesNotExist(),
				check.That(data.ResourceName).Key("service_mesh.0.external_ingress_gateway_enabled").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceMeshProfile(data, true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh.#").HasValue("1"),
				check.That(data.ResourceName).Key("service_mesh.0.internal_ingress_gateway_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("service_mesh.0.external_ingress_gateway_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceMeshProfileDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_mesh.#").HasValue("0"),
				check.That(data.ResourceName).Key("service_mesh.0.internal_ingress_gateway_enabled").DoesNotExist(),
				check.That(data.ResourceName).Key("service_mesh.0.external_ingress_gateway_enabled").DoesNotExist(),
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
			Config: r.privateClusterConfig(data, true),
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
			Config: r.privateClusterConfig(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
				check.That(data.ResourceName).Key("private_cluster.#").HasValue("1"),
			),
		},
		data.ImportStep(),
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

func TestAccKubernetesAutomaticCluster_apiServerAuthorizedIPRanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.apiServerAuthorizedIPRangesConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api_server_access.#").HasValue("1"),
				check.That(data.ResourceName).Key("api_server_access.0.authorized_ip_ranges.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_privateClusterPublicFQDNEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateClusterWithPublicFQDNConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("private_cluster.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_cluster.0.public_fully_qualified_domain_name_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_webAppRoutingIngressDefaultNginx(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webAppRoutingIngressDefaultNginxConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("web_app_routing_ingress.#").HasValue("1"),
				check.That(data.ResourceName).Key("web_app_routing_ingress.0.default_nginx_controller").HasValue("External"),
				check.That(data.ResourceName).Key("web_app_routing_ingress.0.dns_zone_ids.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_webAppRoutingIngressIstioEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webAppRoutingIngressIstioEnabledConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("web_app_routing_ingress.#").HasValue("1"),
				check.That(data.ResourceName).Key("web_app_routing_ingress.0.istio_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_networkingComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.NetworkingConfigComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("hosted_system.#").HasValue("1"),
				check.That(data.ResourceName).Key("api_server_access.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_cluster.#").HasValue("1"),
				check.That(data.ResourceName).Key("web_app_routing_ingress.#").HasValue("1"),
				check.That(data.ResourceName).Key("service_mesh.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_webAppRoutingWithMultipleDnsZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webAppRoutingWithMultipleDnsZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r KubernetesAutomaticClusterResource) serviceMeshProfile(data acceptance.TestData, internalIngressEnabled bool, externalIngressEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name


  hosted_system {
    node_subnet_id        = azurerm_subnet.node.id
    system_node_subnet_id = azurerm_subnet.systemnode.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access {
    subnet_id = azurerm_subnet.api.id
  }

  service_mesh {
    internal_ingress_gateway_enabled = %[4]t
    external_ingress_gateway_enabled = %[5]t
    revisions                        = ["asm-1-28"]
  }
}
`, data.RandomInteger, data.Locations.Primary, r.networkTemplate(data), internalIngressEnabled, externalIngressEnabled)
}

func (r KubernetesAutomaticClusterResource) serviceMeshProfileDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name


  hosted_system {
    node_subnet_id        = azurerm_subnet.node.id
    system_node_subnet_id = azurerm_subnet.systemnode.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access {
    subnet_id = azurerm_subnet.api.id
  }
}
`, data.RandomInteger, data.Locations.Primary, r.networkTemplate(data))
}

func (KubernetesAutomaticClusterResource) privateClusterConfig(data acceptance.TestData, enablePrivateCluster bool) string {
	privateClusterBlock := ""
	if enablePrivateCluster {
		privateClusterBlock = `
	private_cluster {
	}
`
	}

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

%s

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, privateClusterBlock)
}

func (r KubernetesAutomaticClusterResource) privateClusterWithPrivateDNSZoneConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

%s

resource "azurerm_private_dns_zone" "test" {
  name                = "privatelink.%s.azmk8s.io"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_private_dns_zone.test.id
  role_definition_name = "Private DNS Zone Contributor"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  private_cluster {
    private_dns_zone_id = azurerm_private_dns_zone.test.id
  }

  hosted_system {
    node_subnet_id        = azurerm_subnet.node.id
    system_node_subnet_id = azurerm_subnet.systemnode.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access {
    subnet_id = azurerm_subnet.api.id
  }

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, data.RandomInteger, data.Locations.Primary, r.networkTemplate(data), data.Locations.Primary, data.RandomInteger)
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


  private_cluster {
    private_dns_zone_id = "System"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r KubernetesAutomaticClusterResource) apiServerVnetIntegrationConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

%[3]s

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name


  hosted_system {
    node_subnet_id        = azurerm_subnet.node.id
    system_node_subnet_id = azurerm_subnet.systemnode.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access {
    subnet_id = azurerm_subnet.api.id
  }
}
`, data.Locations.Primary, data.RandomInteger, r.networkTemplate(data))
}

func (KubernetesAutomaticClusterResource) apiServerAuthorizedIPRangesConfig(data acceptance.TestData) string {
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


  identity {
    type = "SystemAssigned"
  }

  api_server_access {
    authorized_ip_ranges = [
      "10.10.0.0/24",
      "10.11.0.0/24",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) privateClusterWithPublicFQDNConfig(data acceptance.TestData) string {
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


  private_cluster {
    public_fully_qualified_domain_name_enabled = true
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) webAppRoutingIngressDefaultNginxConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctest%d.example.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }

  web_app_routing_ingress {
    dns_zone_ids             = [azurerm_dns_zone.test.id]
    default_nginx_controller = "External"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r KubernetesAutomaticClusterResource) serviceMeshIngressConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

%s

resource "azurerm_dns_zone" "test" {
  name                = "acctest%d.example.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  hosted_system {
    node_subnet_id        = azurerm_subnet.node.id
    system_node_subnet_id = azurerm_subnet.systemnode.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access {
    subnet_id = azurerm_subnet.api.id
  }

  web_app_routing_ingress {
    istio_enabled = false

  }
  service_mesh {
    internal_ingress_gateway_enabled = true
    external_ingress_gateway_enabled = true
    proxy_redirect_mechanism         = "CNIChaining"
    revisions                        = ["asm-1-28"]
  }

}
`, data.RandomInteger, data.Locations.Primary, r.networkTemplate(data), data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) webAppRoutingIngressIstioEnabledConfig(data acceptance.TestData) string {
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

  identity {
    type = "SystemAssigned"
  }

  web_app_routing_ingress {
    istio_enabled = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r KubernetesAutomaticClusterResource) NetworkingConfigComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

%[3]s

resource "azurerm_dns_zone" "test" {
  name                = "acctest%[1]d.example.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  hosted_system {
    node_subnet_id        = azurerm_subnet.node.id
    system_node_subnet_id = azurerm_subnet.systemnode.id
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  api_server_access {
    subnet_id = azurerm_subnet.api.id
  }

  private_cluster {
    public_fully_qualified_domain_name_enabled = true
    private_dns_zone_id                        = "None"
  }

  web_app_routing_ingress {
    dns_zone_ids             = [azurerm_dns_zone.test.id]
    default_nginx_controller = "External"
  }

  service_mesh {
    internal_ingress_gateway_enabled = true
    external_ingress_gateway_enabled = true
    proxy_redirect_mechanism         = "CNIChaining"
    revisions                        = ["asm-1-28"]
  }

  depends_on = [
    azurerm_role_assignment.network,
  ]
}
`, data.RandomInteger, data.Locations.Primary, r.networkTemplate(data))
}

func (KubernetesAutomaticClusterResource) webAppRoutingWithMultipleDnsZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%[2]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_zone" "test2" {
  name                = "acctestzone2%[2]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }

  web_app_routing_ingress {
    dns_zone_ids             = [azurerm_dns_zone.test.id, azurerm_dns_zone.test2.id]
    default_nginx_controller = "External"
  }
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) networkTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[1]d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "node" {
  name                 = "acctestsubnet%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "api" {
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

resource "azurerm_subnet" "systemnode" {
  name                 = "acctestsubnet2%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.1.0/24"]

  lifecycle {
    ignore_changes = [
      delegation
    ]
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
`, data.RandomInteger)
}
