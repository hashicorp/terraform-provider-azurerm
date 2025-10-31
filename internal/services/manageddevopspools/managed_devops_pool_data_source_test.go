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
				check.That(data.ResourceName).Key("maximum_concurrency").HasValue("1"),
				check.That(data.ResourceName).Key("dev_center_project_resource_id").Exists(),
				check.That(data.ResourceName).Key("azure_devops_organization_profile.0.organization.0.url").Exists(),
				check.That(data.ResourceName).Key("vmss_fabric_profile.0.sku_name").HasValue("Standard_B1s"),
				check.That(data.ResourceName).Key("vmss_fabric_profile.0.image.0.well_known_image_name").HasValue("ubuntu-24.04"),
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
