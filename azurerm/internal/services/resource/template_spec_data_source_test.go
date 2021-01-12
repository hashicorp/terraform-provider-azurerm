package resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type TemplateSpecDataSource struct {
}

func TestAccTemplateSpecDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_template_spec", "test")
	r := TemplateSpecDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (TemplateSpecDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_template_spec" "test" {
  name                = azurerm_template_spec.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, TemplateSpecResource{}.basic(data))
}
