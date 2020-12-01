package devspace_test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devspace/parse"
)

func TestAccDevSpaceController_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_devspace_controller", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevSpaceControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevSpaceController_basic(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevSpaceControllerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDevSpaceController_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_devspace_controller", "test")
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDevSpaceControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDevSpaceController_basic(data, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckDevSpaceControllerExists(data.ResourceName),
				),
			},
			{
				Config:      testAccDevSpaceController_requiresImport(data, clientId, clientSecret),
				ExpectError: acceptance.RequiresImportError("azurerm_devspace_controller"),
			},
		},
	})
}

func testCheckDevSpaceControllerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DevSpace.ControllersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ControllerID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)

		if err == nil {
			return nil
		}

		if result.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevSpace Controller %q (Resource Group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("Bad: Get devSpaceControllerClient: %+v", err)
	}
}

func testCheckDevSpaceControllerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DevSpace.ControllersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_devspace_controller" {
			continue
		}

		log.Printf("[WARN] azurerm_devspace_controller still exists in state file.")

		id, err := parse.ControllerID(rs.Primary.ID)
		if err != nil {
			return err
		}

		result, err := client.Get(ctx, id.ResourceGroup, id.Name)

		if err == nil {
			return fmt.Errorf("DevSpace Controller still exists:\n%#v", result)
		}

		if result.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return nil
}

func testAccDevSpaceController_basic(data acceptance.TestData, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_prefix          = "acctestaks1"

  linux_profile {
    admin_username = "acctestuser1"

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
}

resource "azurerm_devspace_controller" "test" {
  name                                     = "acctestdsc%d"
  location                                 = azurerm_resource_group.test.location
  resource_group_name                      = azurerm_resource_group.test.name
  target_container_host_resource_id        = azurerm_kubernetes_cluster.test.id
  target_container_host_credentials_base64 = base64encode(azurerm_kubernetes_cluster.test.kube_config_raw)
  sku_name                                 = "S1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, clientId, clientSecret, data.RandomInteger)
}

func testAccDevSpaceController_requiresImport(data acceptance.TestData, clientId string, clientSecret string) string {
	template := testAccDevSpaceController_basic(data, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_devspace_controller" "import" {
  name                                     = azurerm_devspace_controller.test.name
  location                                 = azurerm_devspace_controller.test.location
  resource_group_name                      = azurerm_devspace_controller.test.resource_group_name
  target_container_host_resource_id        = azurerm_devspace_controller.test.target_container_host_resource_id
  target_container_host_credentials_base64 = base64encode(azurerm_kubernetes_cluster.test.kube_config_raw)
  sku_name                                 = azurerm_devspace_controller.test.sku_name
}
`, template)
}
