package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type DatabaseSystemDataSource struct{}

func TestDatabaseSystemDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.DatabaseSystemDataSource{}.ResourceType(), "test")
	r := DatabaseSystemDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("network_anchor_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
	})
}

func (d DatabaseSystemDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_database_system" "test" {
  name                = azurerm_oracle_database_system.test.name
  resource_group_name = azurerm_oracle_database_system.test.resource_group_name
}
`, DatabaseSystemResource{}.basic(data))
}
