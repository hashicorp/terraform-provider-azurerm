package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDevSpaceController_basic(t *testing.T) {
	resourceName := "azurerm_devspace_controller.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevSpaceControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevSpaceController_basic(rInt, location, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevSpaceControllerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMDevSpaceController_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_devspace_controller.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevSpaceControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevSpaceController_basic(rInt, location, clientId, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevSpaceControllerExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMDevSpaceController_requiresImport(rInt, location, clientId, clientSecret),
				ExpectError: testRequiresImportError("azurerm_devspace_controller"),
			},
		},
	})
}

func testCheckAzureRMDevSpaceControllerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		ctrlName := rs.Primary.Attributes["name"]
		resGroupName, hasReseGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasReseGroup {
			return fmt.Errorf("Bad: no resource group found in state for DevSpace Controller: %s", ctrlName)
		}

		client := testAccProvider.Meta().(*ArmClient).devSpace.ControllersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		result, err := client.Get(ctx, resGroupName, ctrlName)

		if err == nil {
			return nil
		}

		if result.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevSpace Controller %q (Resource Group: %q) does not exist", ctrlName, resGroupName)
		}

		return fmt.Errorf("Bad: Get devSpaceControllerClient: %+v", err)
	}
}

func testCheckAzureRMDevSpaceControllerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).devSpace.ControllersClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_devspace_controller" {
			continue
		}

		log.Printf("[WARN] azurerm_devspace_controller still exists in state file.")

		ctrlName := rs.Primary.Attributes["name"]
		resGroupName := rs.Primary.Attributes["resource_group_name"]

		result, err := client.Get(ctx, resGroupName, ctrlName)

		if err == nil {
			return fmt.Errorf("DevSpace Controller still exists:\n%#v", result)
		}

		if result.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return nil
}

func testAccAzureRMDevSpaceController_basic(rInt int, location string, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
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

  agent_pool_profile {
    name    = "default"
    count   = "1"
    vm_size = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "%s"
    client_secret = "%s"
  }
}

resource "azurerm_devspace_controller" "test" {
  name                                     = "acctestdsc%d"
  location                                 = "${azurerm_resource_group.test.location}"
  resource_group_name                      = "${azurerm_resource_group.test.name}"
  host_suffix                              = "suffix"
  target_container_host_resource_id        = "${azurerm_kubernetes_cluster.test.id}"
  target_container_host_credentials_base64 = "${base64encode(azurerm_kubernetes_cluster.test.kube_config_raw)}"

  sku {
    name = "S1"
    tier = "Standard"
  }
}
`, rInt, location, rInt, clientId, clientSecret, rInt)
}

func testAccAzureRMDevSpaceController_requiresImport(rInt int, location string, clientId string, clientSecret string) string {
	template := testAccAzureRMDevSpaceController_basic(rInt, location, clientId, clientSecret)
	return fmt.Sprintf(`
%s

resource "azurerm_devspace_controller" "import" {
  name                                     = "${azurerm_devspace_controller.test.name}"
  location                                 = "${azurerm_devspace_controller.test.location}"
  resource_group_name                      = "${azurerm_devspace_controller.test.resource_group_name}"
  host_suffix                              = "${azurerm_devspace_controller.test.host_suffix}"
  target_container_host_resource_id        = "${azurerm_kubernetes_cluster.test.id}"
  target_container_host_credentials_base64 = "${base64encode(azurerm_kubernetes_cluster.test.kube_config_raw)}"

  sku {
    name = "S1"
    tier = "Standard"
  }
}
`, template)
}
