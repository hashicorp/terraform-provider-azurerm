package synapse_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccSynapseSparkPool_V0ToV1_530 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.3.0 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccSynapseSparkPool_V0ToV1_530(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_spark_pool", "test")
	r := SynapseSparkPoolResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-synapse-%[2]d/providers/Microsoft.Synapse/workspaces/acctestsw%[2]d/bigDataPools/acctestSSP%[3]s", data.Subscriptions.Primary, data.RandomInteger, data.RandomString)),
			),
		},
		{
			Config: r.basicV0(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-synapse-%[2]d/providers/microsoft.synapse/workspaces/acctestsw%[2]d/bigDataPools/acctestSSP%[3]s", data.Subscriptions.Primary, data.RandomInteger, data.RandomString)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-synapse-%[2]d/providers/Microsoft.Synapse/workspaces/acctestsw%[2]d/bigDataPools/acctestSSP%[3]s", data.Subscriptions.Primary, data.RandomInteger, data.RandomString)),
			),
		},
	}, "5.3.0")
}

func (r SynapseSparkPoolResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_spark_pool" "test-import" {
  name                 = "acctestSSP%[4]s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  node_count           = 3
  spark_version        = "3.5"
}

removed {
  from = azurerm_synapse_spark_pool.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_synapse_spark_pool.test-import
  id = "/subscriptions/%[2]s/resourcegroups/acctestRG-synapse-%[3]d/providers/microsoft.synapse/workspaces/acctestsw%[3]d/bigDataPools/acctestSSP%[4]s"
}
`, r.template(data, data.Locations.Primary), data.Subscriptions.Primary, data.RandomInteger, data.RandomString)
}

func (r SynapseSparkPoolResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_spark_pool" "test-import" {
  name                 = "acctestSSP%[2]s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  node_size_family     = "MemoryOptimized"
  node_size            = "Small"
  node_count           = 3
  spark_version        = "3.5"
}
`, r.template(data, data.Locations.Primary), data.RandomString)
}
