package devcenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/galleries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DevCenterGalleryResource struct{}

func (r DevCenterGalleryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := galleries.ParseGalleryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.DevCenter.V20230401.Galleries.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s, %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func TestAccDevCenterGallery_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_gallery", "test")
	r := DevCenterGalleryResource{}

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

func TestAccDevCenterGallery_requireImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_center_gallery", "test")
	r := DevCenterGalleryResource{}

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

func (r DevCenterGalleryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_gallery" "test" {
  name                = "acctestgallery%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  dev_center_name     = azurerm_dev_center.test.name
  gallery_resource_id = azurerm_shared_image_gallery.test.id
}
`, r.template(data), data.RandomString)
}

func (r DevCenterGalleryResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_center_gallery" "import" {
  name                = azurerm_dev_center_gallery.test.name
  resource_group_name = azurerm_dev_center_gallery.test.resource_group_name
  dev_center_name     = azurerm_dev_center_gallery.test.dev_center_name
  gallery_resource_id = azurerm_dev_center_gallery.test.gallery_resource_id
}
`, r.basic(data))
}

func (r DevCenterGalleryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%[1]s"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestua-%[1]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_center" "test" {
  location            = azurerm_resource_group.test.location
  name                = "acctestdc-%[1]s"
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_shared_image_gallery" "test" {
  name                = "acctestsig%[1]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_shared_image_gallery.test.id
  role_definition_name = "Owner"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}


`, data.RandomString, data.Locations.Primary)
}
