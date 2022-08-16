package compute_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type GalleryApplicationResource struct{}

func TestAccGalleryApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplication_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccGalleryApplication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplication_description(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.description(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.descriptionUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.description(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplication_endOfLifeDate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}

	endOfLifeDate := time.Now().Add(time.Hour * 10).Format(time.RFC3339)
	endOfLifeDateUpdated := time.Now().Add(time.Hour * 20).Format(time.RFC3339)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.endOfLifeDate(data, endOfLifeDate),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.endOfLifeDate(data, endOfLifeDateUpdated),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplication_eula(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.eula(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.eulaUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.eula(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplication_privacyStatementURI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.privacyStatementURI(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.privacyStatementURIUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplication_releaseNoteURI(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.releaseNoteURI(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.releaseNoteURIUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplication_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r GalleryApplicationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.GalleryApplicationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Compute.GalleryApplicationsClient.Get(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.ID != nil), nil
}

func (r GalleryApplicationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-compute-%[2]d"
  location = "%[1]s"
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r GalleryApplicationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "import" {
  name              = azurerm_gallery_application.test.name
  gallery_id        = azurerm_gallery_application.test.gallery_id
  location          = azurerm_gallery_application.test.location
  supported_os_type = azurerm_gallery_application.test.supported_os_type
}
`, config)
}

func (r GalleryApplicationResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  description           = "This is the gallery application description."
  end_of_life_date      = "%s"
  eula                  = "https://eula.net"
  privacy_statement_uri = "https://privacy.statement.net"
  release_note_uri      = "https://release.note.net"

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, time.Now().Add(time.Hour*10).Format(time.RFC3339))
}

func (r GalleryApplicationResource) description(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  description = "This is the gallery application description."
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) descriptionUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  description = "This is the gallery application description updated."
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) endOfLifeDate(data acceptance.TestData, endOfLifeDate string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  end_of_life_date = "%s"
}
`, template, data.RandomInteger, endOfLifeDate)
}

func (r GalleryApplicationResource) eula(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  eula = "https://eula.net"
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) eulaUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  eula = "https://eula2.net"
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) privacyStatementURI(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  privacy_statement_uri = "https://privacy.statement.net"
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) privacyStatementURIUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  privacy_statement_uri = "https://privacy.statement2.net"
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) releaseNoteURI(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  release_note_uri = "https://release.note2.net"
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) releaseNoteURIUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  release_note_uri = "https://release.note.net"
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) tags(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r GalleryApplicationResource) tagsUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"

  tags = {
    ENV = "Test2"
  }
}
`, template, data.RandomInteger)
}
