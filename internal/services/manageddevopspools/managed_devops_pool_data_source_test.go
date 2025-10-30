package manageddevopspools_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagedDevOpsPoolDataSource struct{}

func TestAccManagedDevOpsPoolDataSource_basic(t *testing.T) {
	if os.Getenv("ARM_MANAGED_DEVOPS_ORG_URL") == "" {
		t.Skip("Skipping as `ARM_MANAGED_DEVOPS_ORG_URL` is not specified")
	}

	data := acceptance.BuildTestData(t, "data.azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
			),
		},
	})
}

func (ManagedDevOpsPoolDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_managed_devops_pool" "test" {
  name                = azurerm_managed_devops_pool.test.name
  resource_group_name = azurerm_managed_devops_pool.test.resource_group_name
}
`, ManagedDevOpsPoolResource{}.basic(data))
}
