package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SharedImageResource struct {
}

func TestAccSharedImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue(""),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImage_basic_hyperVGeneration_V2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "V2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("hyper_v_generation").HasValue("V2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue(""),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSharedImage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "V1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("os_type").HasValue("Linux"),
				check.That(data.ResourceName).Key("hyper_v_generation").HasValue("V1"),
				check.That(data.ResourceName).Key("description").HasValue("Wubba lubba dub dub"),
				check.That(data.ResourceName).Key("eula").HasValue("Do you agree there's infinite Rick's and Infinite Morty's?"),
				check.That(data.ResourceName).Key("privacy_statement_uri").HasValue("https://council.of.ricks/privacy-statement"),
				check.That(data.ResourceName).Key("release_note_uri").HasValue("https://council.of.ricks/changelog.md"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImage_specialized(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.specialized(data, "V1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t SharedImageResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SharedImageID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.GalleryImagesClient.Get(ctx, id.ResourceGroup, id.GalleryName, id.ImageName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Shared Image %q", id.String())
	}

	return utils.Bool(resp.ID != nil), nil
}

func (SharedImageResource) basic(data acceptance.TestData, hyperVGen string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

variable "hyper_v_generation" {
  default = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  hyper_v_generation  = var.hyper_v_generation != "" ? var.hyper_v_generation : null

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, hyperVGen, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (SharedImageResource) specialized(data acceptance.TestData, hyperVGen string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

variable "hyper_v_generation" {
  default = "%s"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  specialized         = true
  hyper_v_generation  = var.hyper_v_generation != "" ? var.hyper_v_generation : null

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, hyperVGen, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r SharedImageResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_shared_image" "import" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
  location            = azurerm_shared_image.test.location
  os_type             = azurerm_shared_image.test.os_type

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, r.basic(data, ""), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (SharedImageResource) complete(data acceptance.TestData, hyperVGen string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                  = "acctestimg%d"
  gallery_name          = azurerm_shared_image_gallery.test.name
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  os_type               = "Linux"
  hyper_v_generation    = "%s"
  description           = "Wubba lubba dub dub"
  eula                  = "Do you agree there's infinite Rick's and Infinite Morty's?"
  privacy_statement_uri = "https://council.of.ricks/privacy-statement"
  release_note_uri      = "https://council.of.ricks/changelog.md"

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }

  purchase_plan {
    name      = "AccTestPlan"
    publisher = "AccTestPlanPublisher"
    product   = "AccTestPlanProduct"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, hyperVGen, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
