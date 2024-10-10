package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/marketplacegalleryimages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StackHCIMarketplaceGalleryImageResource struct {
	// az vm image list --all --output table --sku 2022-datacenter-azure-edition-core
	imageVersion string
}

func TestAccStackHCIMarketplaceGalleryImage_basic(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_marketplace_gallery_image", "test")
	r := StackHCIMarketplaceGalleryImageResource{
		imageVersion: "20348.2402.240607",
	}

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

func TestAccStackHCIMarketplaceGalleryImage_complete(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_marketplace_gallery_image", "test")
	r := StackHCIMarketplaceGalleryImageResource{
		imageVersion: "20348.2582.240703",
	}

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

func TestAccStackHCIMarketplaceGalleryImage_update(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_marketplace_gallery_image", "test")
	r := StackHCIMarketplaceGalleryImageResource{
		imageVersion: "20348.2655.240810",
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data),
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
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStackHCIMarketplaceGalleryImage_requiresImport(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "azurerm_stack_hci_marketplace_gallery_image", "test")
	r := StackHCIMarketplaceGalleryImageResource{
		imageVersion: "20348.2655.240905",
	}

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

func (r StackHCIMarketplaceGalleryImageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.MarketplaceGalleryImages
	id, err := marketplacegalleryimages.ParseMarketplaceGalleryImageID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StackHCIMarketplaceGalleryImageResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_marketplace_gallery_image" "test" {
  name                = "acctest-mgi-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  hyperv_generation   = "V2"
  os_type             = "Windows"
  version             = "%s"
  identifier {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition-core"
  }

  depends_on = [azurerm_role_assignment.test]
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv), r.imageVersion)
}

func (r StackHCIMarketplaceGalleryImageResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_marketplace_gallery_image" "import" {
  name                = azurerm_stack_hci_marketplace_gallery_image.test.name
  resource_group_name = azurerm_stack_hci_marketplace_gallery_image.test.resource_group_name
  location            = azurerm_stack_hci_marketplace_gallery_image.test.location
  custom_location_id  = azurerm_stack_hci_marketplace_gallery_image.test.custom_location_id
  hyperv_generation   = azurerm_stack_hci_marketplace_gallery_image.test.hyperv_generation
  os_type             = azurerm_stack_hci_marketplace_gallery_image.test.os_type
  version             = azurerm_stack_hci_marketplace_gallery_image.test.version
  identifier {
    publisher = azurerm_stack_hci_marketplace_gallery_image.test.identifier.0.publisher
    offer     = azurerm_stack_hci_marketplace_gallery_image.test.identifier.0.offer
    sku       = azurerm_stack_hci_marketplace_gallery_image.test.identifier.0.sku
  }
}
`, config)
}

func (r StackHCIMarketplaceGalleryImageResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_storage_path" "test" {
  name                = "acctest-sp-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  path                = "C:\\ClusterStorage\\UserStorage_2\\sp-mgi-%[2]s"
}

resource "azurerm_stack_hci_marketplace_gallery_image" "test" {
  name                = "acctest-mgi-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  hyperv_generation   = "V2"
  os_type             = "Windows"
  version             = "%s"
  storage_path_id     = azurerm_stack_hci_storage_path.test.id
  identifier {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition-core"
  }
  tags = {
    foo = "bar"
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv), r.imageVersion)
}

func (r StackHCIMarketplaceGalleryImageResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_storage_path" "test" {
  name                = "acctest-sp-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  path                = "C:\\ClusterStorage\\UserStorage_2\\sp-mgi-%[2]s"
}

resource "azurerm_stack_hci_marketplace_gallery_image" "test" {
  name                = "acctest-mgi-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %[3]q
  hyperv_generation   = "V2"
  os_type             = "Windows"
  version             = "%s"
  storage_path_id     = azurerm_stack_hci_storage_path.test.id
  identifier {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition-core"
  }
  tags = {
    foo = "bar"
    env = "test"
  }
}
`, template, data.RandomString, os.Getenv(customLocationIdEnv), r.imageVersion)
}

func (r StackHCIMarketplaceGalleryImageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-mgi-%s"
  location = "%s"
}

data "azurerm_client_config" "test" {}

// service principal of 'Microsoft.AzureStackHCI Resource Provider'
data "azuread_service_principal" "hciRp" {
  client_id = "1412d89f-b8a8-4111-b4fd-e82905cbd85d"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Azure Connected Machine Resource Manager"
  principal_id         = data.azuread_service_principal.hciRp.object_id
}
`, data.RandomString, data.Locations.Primary)
}
