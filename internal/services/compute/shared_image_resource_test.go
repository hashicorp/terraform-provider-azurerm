package compute_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SharedImageResource struct{}

func TestAccSharedImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithHyperVGen(data, ""),
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
			Config: r.basicWithHyperVGen(data, "V2"),
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
			Config: r.basicWithHyperVGen(data, ""),
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
			Config: r.completeWithHyperVGen(data, "V1"),
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

func TestAccSharedImage_withTrustedLaunchEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTrustedLaunchEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImage_withAcceleratedNetworkSupportEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAcceleratedNetworkSupportEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImage_description(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.description(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.descriptionUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImage_releaseNoteURI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.releaseNoteURI(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.releaseNoteURIUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImage_disallowedDiskTypes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithDiskTypesNotAllowed(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithDiskTypesNotAllowedUpdated(data),
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
		{
			Config: r.basicWithDiskTypesNotAllowed(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImage_endOfLifeDate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}

	endOfLifeDate := time.Now().Add(time.Hour * 10).Format(time.RFC3339)
	endOfLifeDateUpdated := time.Now().Add(time.Hour * 20).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.endOfLifeDate(data, endOfLifeDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.endOfLifeDate(data, endOfLifeDateUpdated),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSharedImage_recommended(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_shared_image", "test")
	r := SharedImageResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithRecommended(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithRecommendedUpdated(data),
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
		{
			Config: r.basicWithRecommended(data),
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

func (SharedImageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}
resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (SharedImageResource) basicWithHyperVGen(data acceptance.TestData, hyperVGen string) string {
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
`, r.basicWithHyperVGen(data, ""), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (SharedImageResource) completeWithHyperVGen(data acceptance.TestData, hyperVGen string) string {
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

func (SharedImageResource) withTrustedLaunchEnabled(data acceptance.TestData) string {
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
  name                   = "acctestimg%d"
  gallery_name           = azurerm_shared_image_gallery.test.name
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  os_type                = "Linux"
  hyper_v_generation     = "V2"
  trusted_launch_enabled = true

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (SharedImageResource) withAcceleratedNetworkSupportEnabled(data acceptance.TestData) string {
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
  name                = "acctestimg%d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  accelerated_network_support_enabled = true

  identifier {
    publisher = "AccTesPublisher%d"
    offer     = "AccTesOffer%d"
    sku       = "AccTesSku%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (SharedImageResource) description(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  description         = "image description"

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (SharedImageResource) descriptionUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  description         = "image description updated"

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (SharedImageResource) releaseNoteURI(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  release_note_uri    = "https://test.com/changelog.md"

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (SharedImageResource) releaseNoteURIUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"
  release_note_uri    = "https://test.com/changelog2.md"

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (SharedImageResource) basicWithDiskTypesNotAllowed(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  disk_types_not_allowed = [
    "Standard_LRS",
  ]

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (SharedImageResource) basicWithDiskTypesNotAllowedUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  disk_types_not_allowed = [
    "Standard_LRS",
    "Premium_LRS",
  ]

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (SharedImageResource) endOfLifeDate(data acceptance.TestData, endOfLifeDate string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  end_of_life_date = "%[3]s"

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger, endOfLifeDate)
}

func (SharedImageResource) basicWithRecommended(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  max_recommended_vcpu_count   = 8
  min_recommended_vcpu_count   = 7
  max_recommended_memory_in_gb = 6
  min_recommended_memory_in_gb = 5

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (SharedImageResource) basicWithRecommendedUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_shared_image" "test" {
  name                = "acctestimg%[2]d"
  gallery_name        = azurerm_shared_image_gallery.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  os_type             = "Linux"

  max_recommended_vcpu_count   = 4
  min_recommended_vcpu_count   = 3
  max_recommended_memory_in_gb = 2
  min_recommended_memory_in_gb = 1

  identifier {
    publisher = "AccTesPublisher%[2]d"
    offer     = "AccTesOffer%[2]d"
    sku       = "AccTesSku%[2]d"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}
