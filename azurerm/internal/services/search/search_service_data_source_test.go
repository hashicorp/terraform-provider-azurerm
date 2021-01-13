package search_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SearchServiceDataSource struct {
}

func TestAccDataSourceSearchService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_search_service", "test")
	r := SearchServiceDataSource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("replica_count").Exists(),
				check.That(data.ResourceName).Key("partition_count").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("identity.type").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (SearchServiceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_search_service" "test" {
  name                = azurerm_search_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, SearchServiceResource{}.basic(data))
}
