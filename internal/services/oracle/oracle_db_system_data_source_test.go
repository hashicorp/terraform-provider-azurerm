package oracle_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
)

type DbSystemDataSource struct{}

func TestDbSystemDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, oracle.DbSystemDataSource{}.ResourceType(), "test")
	r := DbSystemDataSource{}

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

func (d DbSystemDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_oracle_db_system" "test" {
  name                = azurerm_oracle_db_system.test.name
  resource_group_name = azurerm_oracle_db_system.test.resource_group_name
}
`, DbSystemResource{}.basic(data))
}
