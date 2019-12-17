package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMContainerRegistryWebhook_basic(t *testing.T) {
	resourceName := "azurerm_container_registry_webhook.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryWebhook_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_withTags(t *testing.T) {
	resourceName := "azurerm_container_registry_webhook.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMContainerRegistryWebhook_withTags(ri, acceptance.Location())
	postConfig := testAccAzureRMContainerRegistryWebhook_withTagsUpdate(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.label", "test"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.label", "test1"),
					resource.TestCheckResourceAttr(resourceName, "tags.ENV", "prod"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_actions(t *testing.T) {
	resourceName := "azurerm_container_registry_webhook.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMContainerRegistryWebhook_actions(ri, acceptance.Location())
	postConfig := testAccAzureRMContainerRegistryWebhook_actionsUpdate(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: postConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_status(t *testing.T) {
	resourceName := "azurerm_container_registry_webhook.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMContainerRegistryWebhook_status(ri, acceptance.Location())
	postConfig := testAccAzureRMContainerRegistryWebhook_statusUpdate(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_serviceUri(t *testing.T) {
	resourceName := "azurerm_container_registry_webhook.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMContainerRegistryWebhook_serviceUri(ri, acceptance.Location())
	postConfig := testAccAzureRMContainerRegistryWebhook_serviceUriUpdate(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service_uri", "https://mywebhookreceiver.example/mytag"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service_uri", "https://my.webhookreceiver.example/mytag/2"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_scope(t *testing.T) {
	resourceName := "azurerm_container_registry_webhook.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMContainerRegistryWebhook_scope(ri, acceptance.Location())
	postConfig := testAccAzureRMContainerRegistryWebhook_scopeUpdate(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "scope", "mytag:*"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "scope", "mytag:4"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerRegistryWebhook_customHeaders(t *testing.T) {
	resourceName := "azurerm_container_registry_webhook.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMContainerRegistryWebhook_customHeaders(ri, acceptance.Location())
	postConfig := testAccAzureRMContainerRegistryWebhook_customHeadersUpdate(ri, acceptance.Location())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "custom_headers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers.Content-Type", "application/json"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryWebhookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "custom_headers.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers.Content-Type", "application/xml"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers.Accept-Charset", "utf-8"),
				),
			},
		},
	})
}

func testAccAzureRMContainerRegistryWebhook_basic(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_withTags(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_withTagsUpdate(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_actions(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_actionsUpdate(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_status(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_statusUpdate(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_serviceUri(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_serviceUriUpdate(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_scope(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_scopeUpdate(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_customHeaders(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
}

func testAccAzureRMContainerRegistryWebhook_customHeadersUpdate(rInt int, location string) string {
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
`, rInt, location, rInt, location, rInt, location)
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
