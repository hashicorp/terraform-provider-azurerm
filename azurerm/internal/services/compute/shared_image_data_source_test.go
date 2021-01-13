package compute

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SharedImageDataSource struct {
}

func TestAccDataSourceAzureRMSharedImage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data, ""),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceAzureRMSharedImage_basic_hyperVGeneration_V2(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data, "V2"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("hyper_v_generation").HasValue("V2"),
			),
		},
	})
}

func TestAccDataSourceAzureRMSharedImage_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_shared_image", "test")
	r := SharedImageDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.complete(data, "V1"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("hyper_v_generation").HasValue("V1"),
			),
		},
	})
}

func (SharedImageDataSource) basic(data acceptance.TestData, hyperVGen string) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.basic(data, hyperVGen))
}

func (SharedImageDataSource) complete(data acceptance.TestData, hyperVGen string) string {
	return fmt.Sprintf(`
%s

data "azurerm_shared_image" "test" {
  name                = azurerm_shared_image.test.name
  gallery_name        = azurerm_shared_image.test.gallery_name
  resource_group_name = azurerm_shared_image.test.resource_group_name
}
`, SharedImageResource{}.complete(data, hyperVGen))
}
