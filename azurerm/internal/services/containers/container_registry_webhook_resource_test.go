package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ContainerRegistryWebhookResource struct {
}

func TestAccContainerRegistryWebhook_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")
	r := ContainerRegistryWebhookResource{}

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

func TestAccContainerRegistryWebhook_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")
	r := ContainerRegistryWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("prod"),
			),
		},
	})
}

func TestAccContainerRegistryWebhook_actions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")
	r := ContainerRegistryWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.actions(data),
		},
		data.ImportStep(),
		{
			Config: r.actionsUpdate(data),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistryWebhook_status(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")
	r := ContainerRegistryWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.status(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("enabled"),
			),
		},
		{
			Config: r.statusUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("status").HasValue("disabled"),
			),
		},
	})
}

func TestAccContainerRegistryWebhook_serviceUri(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")
	r := ContainerRegistryWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceUri(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_uri").HasValue("https://mywebhookreceiver.example/mytag"),
			),
		},
		{
			Config: r.serviceUriUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_uri").HasValue("https://my.webhookreceiver.example/mytag/2"),
			),
		},
	})
}

func TestAccContainerRegistryWebhook_scope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")
	r := ContainerRegistryWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scope(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue("mytag:*"),
			),
		},
		{
			Config: r.scopeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue("mytag:4"),
			),
		},
	})
}

func TestAccContainerRegistryWebhook_customHeaders(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")
	r := ContainerRegistryWebhookResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customHeaders(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_headers.%").HasValue("1"),
				check.That(data.ResourceName).Key("custom_headers.Content-Type").HasValue("application/json"),
			),
		},
		{
			Config: r.customHeadersUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_headers.%").HasValue("2"),
				check.That(data.ResourceName).Key("custom_headers.Content-Type").HasValue("application/xml"),
				check.That(data.ResourceName).Key("custom_headers.Accept-Charset").HasValue("utf-8"),
			),
		},
	})
}

func (ContainerRegistryWebhookResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  actions = [
    "push"
  ]

  tags = {
    label = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  actions = [
    "push"
  ]

  tags = {
    label = "test1"
    ENV   = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) actions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) actionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  actions = [
    "push",
    "delete"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) status(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  status = "enabled"

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) statusUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  status = "disabled"

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) serviceUri(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) serviceUriUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://my.webhookreceiver.example/mytag/2"

  status = "disabled"

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) scope(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  scope = "mytag:*"

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) scopeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  scope = "mytag:4"

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) customHeaders(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  custom_headers = {
    "Content-Type" = "application/json"
  }

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (ContainerRegistryWebhookResource) customHeadersUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrwebhooktest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Standard"
}

resource "azurerm_container_registry_webhook" "test" {
  name                = "testwebhook%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  service_uri = "https://mywebhookreceiver.example/mytag"

  custom_headers = {
    "Content-Type"   = "application/xml"
    "Accept-Charset" = "utf-8"
  }

  actions = [
    "push"
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (t ContainerRegistryWebhookResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	registryName := id.Path["registries"]
	name := id.Path["webhooks"]

	resp, err := clients.Containers.WebhooksClient.Get(ctx, resourceGroup, registryName, name)
	if err != nil {
		return nil, fmt.Errorf("reading Container Registry Webhook (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
