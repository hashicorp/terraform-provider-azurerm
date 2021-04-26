package iothub_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type IoTHubSharedAccessPolicyDataSource struct {
}

func TestAccDataSourceIotHubSharedAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_iothub_shared_access_policy", "test")
	r := IoTHubSharedAccessPolicyDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
			),
		},
	})
}

func (IoTHubSharedAccessPolicyDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_iothub_shared_access_policy" "test" {
  name                = azurerm_iothub_shared_access_policy.test.name
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
}
`, IoTHubSharedAccessPolicyResource{}.basic(data))
}
