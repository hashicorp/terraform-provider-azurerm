package machinelearning_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type WorkspaceDataSource struct{}

func TestAccMachineLearningWorkspaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_machine_learning_workspace", "test")
	d := WorkspaceDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: d.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
	})
}

func ( WorkspaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_machine_learning_workspace" "test" {
  name                = azurerm_machine_learning_workspace.test.name
  resource_group_name = azurerm_machine_learning_workspace.test.resource_group_name
}
`, WorkspaceResource{}.complete(data))
}
