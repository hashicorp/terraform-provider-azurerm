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

type StackHCIMarketplaceGalleryImageResource struct{}

func TestAccStackHCIMarketplaceGalleryImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_marketplace_gallery_image", "test")
	r := StackHCIMarketplaceGalleryImageResource{}

	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
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
	r := StackHCIMarketplaceGalleryImageResource{}

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
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_marketplace_gallery_image", "test")
	r := StackHCIMarketplaceGalleryImageResource{}

	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

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
		{
			Config: r.basic(data),
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
	r := StackHCIMarketplaceGalleryImageResource{}

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
  name                = "acctest-mgi-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q
  os_type = "Windows"
  version = "20348.2113.231109"
  identifier {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition"
  }

  depends_on = [azurerm_role_assignment.test]
}
`, template, os.Getenv(customLocationIdEnv))
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
}
`, config)
}

func (r StackHCIMarketplaceGalleryImageResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_marketplace_gallery_image" "test" {
  name                = "acctest-ln-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q

  tags = {
    foo = "bar"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCIMarketplaceGalleryImageResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_stack_hci_marketplace_gallery_image" "test" {
  name                = "acctest-ln-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = %q

  tags = {
    foo = "bar"
    env = "test"
  }
}
`, template, os.Getenv(customLocationIdEnv))
}

func (r StackHCIMarketplaceGalleryImageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}

variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-mgi-${var.random_string}"
  location = var.primary_location
}

data "azurerm_client_config" "test" {}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = "Azure Connected Machine Resource Manager"
  principal_id         = data.azurerm_client_config.test.object_id
}
`, data.Locations.Primary, data.RandomString)
}
