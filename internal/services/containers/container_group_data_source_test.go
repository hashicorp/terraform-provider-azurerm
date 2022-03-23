package containers_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerGroupDataSource struct{}

func TestAccDataSourceContainerGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_group", "test")
	r := ContainerGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("fqdn").Exists(),
			),
		},
	})
}

func (ContainerGroupDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_group" "test" {
  name                = azurerm_container_group.test.name
  resource_group_name = azurerm_container_group.test.resource_group_name
}
`, ContainerGroupResource{}.linuxComplete(data))
}
