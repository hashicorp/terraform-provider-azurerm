package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMContainerRegistryWebhook_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryWebhook_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryWebhook_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.label", "test"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistryWebhook_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.label", "test1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "prod"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_actions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryWebhook_actions(data),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMContainerRegistryWebhook_actionsUpdate(data),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_status(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryWebhook_status(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "enabled"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistryWebhook_statusUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "status", "disabled"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_serviceUri(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryWebhook_serviceUri(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "service_uri", "https://mywebhookreceiver.example/mytag"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistryWebhook_serviceUriUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "service_uri", "https://my.webhookreceiver.example/mytag/2"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_scope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryWebhook_scope(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", "mytag:*"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistryWebhook_scopeUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", "mytag:4"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_customHeaders(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_webhook", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryWebhook_customHeaders(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_headers.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_headers.Content-Type", "application/json"),
				),
			},
			{
				Config: testAccAzureRMContainerRegistryWebhook_customHeadersUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_headers.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_headers.Content-Type", "application/xml"),
					resource.TestCheckResourceAttr(data.ResourceName, "custom_headers.Accept-Charset", "utf-8"),
				),
			},
		},
	})
}

func testAccAzureRMContainerRegistryWebhook_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_actions(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_actionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_status(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_statusUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_serviceUri(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_serviceUriUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_scope(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_scopeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_customHeaders(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testAccAzureRMContainerRegistryWebhook_customHeadersUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
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

func testCheckAzureRMContainerRegistryWebhookDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.WebhooksClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_registry_webhook" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		registryName := rs.Primary.Attributes["registry_name"]
		name := rs.Primary.Attributes["name"]

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, registryName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMContainerRegistryWebhookExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		webhookName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Container Registry Webhook: %s", webhookName)
		}

		registryName, hasRegistryName := rs.Primary.Attributes["registry_name"]
		if !hasRegistryName {
			return fmt.Errorf("Bad: no registry name found in state for Container Registry Webhook: %s", webhookName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.WebhooksClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, registryName, webhookName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Container Registry Webhook %q (resource group: %q) does not exist", webhookName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on WebhooksClient: %+v", err)
		}

		return nil
	}
}
