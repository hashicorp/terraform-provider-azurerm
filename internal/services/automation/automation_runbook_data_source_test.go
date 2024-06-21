// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutomationRunbookDataSource struct{}

func TestAccAutomationRunbookDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_runbook", "test")
	r := AutomationRunbookDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("runbook_type").HasValue("PowerShell"),
				check.That(data.ResourceName).Key("content").HasValue("# Some test content\n# for Terraform acceptance test\n"),
			),
		},
	})
}

func (AutomationRunbookDataSource) basic(data acceptance.TestData) string {
	template := AutomationRunbookResource{}.PSWithContent(data)
	return fmt.Sprintf(`
%s

data "azurerm_automation_runbook" "test" {
  name                    = azurerm_automation_runbook.test.name
  automation_account_name = azurerm_automation_account.test.name
  resource_group_name     = azurerm_resource_group.test.name
}
`, template)
}
