// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
}

type AutomationRunbookDataSource struct{}

func TestAccAutomationRunbookDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_runbook", "test")
	r := AutomationRunbookDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
	})
}

func (AutomationRunbookDataSource) basic(data acceptance.TestData) string {
	template := AutomationRunbookResource{}.withDraft(data)
	return fmt.Sprintf(`
%s

data "azurerm_automation_runbook" "test" {
  name                = azurerm_automation_runbook.test.name
  automation_account_name = azurerm_automation_account.test.name
  resource_group_name = azurerm_resource_group.test.resource_group_name
}
`, template)
}