package digitaltwins_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DigitalTwinsInstanceDataSource struct {
}

func TestAccDigitalTwinsInstanceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
	})
}

func (DigitalTwinsInstanceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_digital_twins_instance" "test" {
  name                = azurerm_digital_twins_instance.test.name
  resource_group_name = azurerm_digital_twins_instance.test.resource_group_name
}
`, DigitalTwinsInstanceResource{}.basic(data))
}
