// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccSynapseIntegrationRuntimeSelfHosted_V1ToV2_530 tests the state migration path from an `id` with
// lowercased static segments to their canonicalized format. It uses v5.3.0 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccSynapseIntegrationRuntimeSelfHosted_V1ToV2_530(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_integration_runtime_self_hosted", "test")
	r := SynapseIntegrationRuntimeSelfHostedResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-synapse-%[2]d/providers/Microsoft.Synapse/workspaces/acctestdf%[2]d/integrationRuntimes/acctestSIR%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicV1(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Ensure import is a no-op to prevent the resource from normalizing the ID due to a combined CreateUpdate func
					plancheck.ExpectResourceAction(importedResourceName, plancheck.ResourceActionNoop),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-synapse-%[2]d/providers/microsoft.synapse/workspaces/acctestdf%[2]d/integrationRuntimes/acctestSIR%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-synapse-%[2]d/providers/Microsoft.Synapse/workspaces/acctestdf%[2]d/integrationRuntimes/acctestSIR%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.3.0")
}

func (r SynapseIntegrationRuntimeSelfHostedResource) basicV1(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_integration_runtime_self_hosted" "test-import" {
  name                 = "acctestSIR%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  description          = "test"
}

removed {
  from = azurerm_synapse_integration_runtime_self_hosted.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_synapse_integration_runtime_self_hosted.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acctestRG-synapse-%[2]d/providers/microsoft.synapse/workspaces/acctestdf%[2]d/integrationRuntimes/acctestSIR%[2]d"
}
`, r.template(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (r SynapseIntegrationRuntimeSelfHostedResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_synapse_integration_runtime_self_hosted" "test-import" {
  name                 = "acctestSIR%[2]d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  description          = "test"
}
`, r.template(data), data.RandomInteger)
}
