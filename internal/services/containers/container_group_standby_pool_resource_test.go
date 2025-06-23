package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbycontainergrouppools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerGroupStandbyPoolResource struct{}

func TestAccContianerGroupStandbyPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group_standby_pool", "test")
	r := ContainerGroupStandbyPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContianerGroupStandbyPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group_standby_pool", "test")
	r := ContainerGroupStandbyPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContianerGroupStandbyPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group_standby_pool", "test")
	r := ContainerGroupStandbyPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContianerGroupStandbyPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_group_standby_pool", "test")
	r := ContainerGroupStandbyPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (ContainerGroupStandbyPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := standbycontainergrouppools.ParseStandbyContainerGroupPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Containers.StandbyContainerGroupPoolsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ContainerGroupStandbyPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_group_standby_pool" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  container_gorup_id  = azurerm_container_group.test.id
  max_ready_capacity  = 2

}
`, r.template(data), data.RandomInteger)
}

func (r ContainerGroupStandbyPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_group_standby_pool" "import" {
  name                = azurerm_container_group_standby_pool.test.name
  resource_group_name = azurerm_container_group_standby_pool.test.resource_group_name
  container_gorup_id  = azurerm_container_group_standby_pool.test.contianer_group_id
  max_ready_capacity  = azurerm_container_group_standby_pool.test.max_ready_capacity
}

`, r.basic(data))
}

func (r ContainerGroupStandbyPoolResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_group_standby_pool" "test" {
  name                     = "acctest-%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  container_gorup_id       = azurerm_container_group.test.id
  container_group_revision = 1
  subnet_ids               = [azurerm_subnet.test.id]
  refill_policy            = "always"
  max_ready_capacity       = 2

  tags = {
    acc = "test"
  }

}
`, r.template(data), data.RandomInteger)
}

func (ContainerGroupStandbyPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "container-contributor" {
  name = "Container Instance Contributor"
}

data "azurerm_role_definition" "nw-contributor" {
  name = "Network Contributor"
}

data "azurerm_role_definition" "mi-contributor" {
  name = "Managed Identity Contributor"
}

data "azurerm_role_definition" "mi-operator" {
  name = "Managed Identity Operator"
}

data "azuread_service_principal" "test" {
  display_name = "Standby Pool Resource Provider"
}

resource "azurerm_role_assignment" "container-contributor" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.container-contributor.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "nw-contributor" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.nw-contributor.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "mi-contributor" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.mi-contributor.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "mi-operator" {
  scope              = azurerm_resource_group.test.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.mi-operator.id}"
  principal_id       = data.azuread_service_principal.test.object_id
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "testsubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  ip_address_type     = "Private"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports {
      port = 80
    }
  }
  dns_config {
    nameservers    = ["reddog.microsoft.com", "somecompany.somedomain"]
    options        = ["one:option", "two:option", "red:option", "blue:option"]
    search_domains = ["default.svc.cluster.local."]
  }

  subnet_ids = [azurerm_subnet.test.id]

  tags = {
    environment = "Testing"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
