package devspace_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/devspace/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// temporarily works around the unused test, since the tests are skipped
var _ interface{} = DevSpaceControllerResource{}
var _ interface{} = DevSpaceControllerResource{}.basic(acceptance.TestData{}, "", "")
var _ interface{} = DevSpaceControllerResource{}.requiresImport(acceptance.TestData{}, "", "")

type DevSpaceControllerResource struct {
}

func TestAccDevSpaceController_basic(t *testing.T) {
	t.Skip("A breaking API change has means new DevSpace Controllers cannot be provisioned, so skipping..")

	data := acceptance.BuildTestData(t, "azurerm_devspace_controller", "test")
	r := DevSpaceControllerResource{}
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, clientId, clientSecret),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDevSpaceController_requiresImport(t *testing.T) {
	t.Skip("A breaking API change has means new DevSpace Controllers cannot be provisioned, so skipping..")

	data := acceptance.BuildTestData(t, "azurerm_devspace_controller", "test")
	r := DevSpaceControllerResource{}
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, clientId, clientSecret),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, clientId, clientSecret),
			ExpectError: acceptance.RequiresImportError("azurerm_devspace_controller"),
		},
	})
}

func (t DevSpaceControllerResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ControllerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevSpace.ControllersClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving DevSpace Controller %q (Resource Group: %q): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ControllerProperties != nil), nil
}

func (DevSpaceControllerResource) basic(data acceptance.TestData, clientId string, clientSecret string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-devspace-%d"
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

func (r DevSpaceControllerResource) requiresImport(data acceptance.TestData, clientId string, clientSecret string) string {
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
`, r.basic(data, clientId, clientSecret))
}
