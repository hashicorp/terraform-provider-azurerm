package extendedlocation_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagerDataSource struct{}

const (
	customLocationName              = "ARM_TEST_CUSTOM_LOCATION_NAME"
	customLocationResourceGroupName = "ARM_TEST_CUSTOM_LOCATION_RESOURCE_GROUP_NAME"
)

func TestAccCustomLocationDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_custom_location", "test")
	d := ManagerDataSource{}

	if os.Getenv(customLocationName) == "" || os.Getenv(customLocationResourceGroupName) == "" {
		t.Skipf("Skipping test due to missing environment variables: %s, %s", customLocationName, customLocationResourceGroupName)
	}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").IsNotEmpty(),
				check.That(data.ResourceName).Key("host_resource_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("display_name").IsNotEmpty(),
				check.That(data.ResourceName).Key("namespace").IsNotEmpty(),
				check.That(data.ResourceName).Key("cluster_extension_ids.#").HasValue("2"),
			),
		},
	})
}

func (d ManagerDataSource) basic() string {
	return fmt.Sprintf(`
provider azurerm {
  features {}
}

data "azurerm_custom_location" "test" {
  name                = %q
  resource_group_name = %q
}
`, os.Getenv(customLocationName), os.Getenv(customLocationResourceGroupName))
}
