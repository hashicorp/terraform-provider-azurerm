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

type GalleryApplicationVersionResource struct{}

func TestAccGalleryApplicationVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_health_check").HasValue("false"),
				check.That(data.ResourceName).Key("exclude_from_latest").HasValue("false"),
				check.That(data.ResourceName).Key("target_region.0.storage_account_type").HasValue("Standard_LRS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplicationVersion_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}
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

func TestAccGalleryApplicationVersion_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}
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

func TestAccGalleryApplicationVersion_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}
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

func TestAccGalleryApplicationVersion_enableHealthCheck(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.enableHealthCheck(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.enableHealthCheckUpdated(data),
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
			Config: r.enableHealthCheck(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplicationVersion_endOfLifeDate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}

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

func TestAccGalleryApplicationVersion_excludeFromLatest(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.excludeFromLatest(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.excludeFromLatestUpdated(data),
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
			Config: r.excludeFromLatest(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplicationVersion_targetRegion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.targetRegion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.targetRegionUpdated(data),
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
			Config: r.targetRegion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGalleryApplicationVersion_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}
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

func (r GalleryApplicationVersionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.GalleryApplicationVersionID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Compute.GalleryApplicationVersionsClient.Get(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, id.VersionName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.ID != nil), nil
}

func (r GalleryApplicationVersionResource) template(data acceptance.TestData) string {
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

resource "azurerm_gallery_application" "test" {
  name              = "acctest-app-%[2]d"
  gallery_id        = azurerm_shared_image_gallery.test.id
  location          = azurerm_resource_group.test.location
  supported_os_type = "Linux"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc%[3]s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test" {
  name                   = "scripts"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_content         = "[scripts file content]"
}

`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r GalleryApplicationVersionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }
}
`, template)
}

func (r GalleryApplicationVersionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "import" {
  name                   = azurerm_gallery_application_version.test.name
  gallery_application_id = azurerm_gallery_application_version.test.gallery_application_id
  location               = azurerm_gallery_application_version.test.location

  manage_action {
    install = azurerm_gallery_application_version.test.manage_action.0.install
    remove  = azurerm_gallery_application_version.test.manage_action.0.remove
  }

  source {
    media_link = azurerm_gallery_application_version.test.source.0.media_link
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }
}
`, config)
}

func (r GalleryApplicationVersionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  enable_health_check = true
  end_of_life_date    = "%s"
  exclude_from_latest = true

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
    update  = "[update command]"
  }

  source {
    media_link                 = azurerm_storage_blob.test.id
    default_configuration_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
    storage_account_type   = "Premium_LRS"
  }

  tags = {
    ENV = "Test"
  }
}
`, template, time.Now().Add(time.Hour*10).Format(time.RFC3339))
}

func (r GalleryApplicationVersionResource) enableHealthCheck(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  enable_health_check = true

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }
}
`, template)
}

func (r GalleryApplicationVersionResource) enableHealthCheckUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  enable_health_check = false

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }
}
`, template)
}

func (r GalleryApplicationVersionResource) endOfLifeDate(data acceptance.TestData, endOfLifeDate string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  end_of_life_date = "%s"

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }
}
`, template, endOfLifeDate)
}

func (r GalleryApplicationVersionResource) excludeFromLatest(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  exclude_from_latest = true

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }
}
`, template)
}

func (r GalleryApplicationVersionResource) excludeFromLatestUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  exclude_from_latest = false

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }
}
`, template)
}

func (r GalleryApplicationVersionResource) targetRegion(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }

  target_region {
    name                   = "%s"
    regional_replica_count = 2
    storage_account_type   = "Premium_LRS"
  }
}
`, template, data.Locations.Secondary)
}

func (r GalleryApplicationVersionResource) targetRegionUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }

  target_region {
    name                   = "%s"
    regional_replica_count = 3
    storage_account_type   = "Premium_LRS"
  }
}
`, template, data.Locations.Secondary)
}

func (r GalleryApplicationVersionResource) tags(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }

  tags = {
    ENV = "Test"
  }
}
`, template)
}

func (r GalleryApplicationVersionResource) tagsUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_gallery_application_version" "test" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.test.id
  location               = azurerm_gallery_application.test.location

  manage_action {
    install = "[install command]"
    remove  = "[remove command]"
  }

  source {
    media_link = azurerm_storage_blob.test.id
  }

  target_region {
    name                   = azurerm_gallery_application.test.location
    regional_replica_count = 1
  }

  tags = {
    ENV = "Test2"
  }
}
`, template)
}
