package devcenter_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"testing"
)

type DevCenterGalleryResource struct{}

func (r DevCenterGalleryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	//TODO implement me
	panic("implement me")
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

func (r DevCenterGalleryResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

provider "azurerm" {
  features {}
}

resource "azurerm_dev_center_gallery" "test" {
  name = "acctestgallery-%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  dev_center_name = azurerm_dev_center.test.name
  gallery_resource_id = azurerm_dev_center_gallery_resource.test.id
}
`, r.template(data), data.RandomString)
}

func (r DevCenterGalleryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name    = "acctestrg-%[1]s"
  location = "%[2]s"
}

resource "azurerm_dev_center" "test" {
  location            = azurerm_resource_group.test.location
  name = "acctestdc-%[1]s"
  resource_group_name = azurerm_resource_group.test.name
}


`)
}
