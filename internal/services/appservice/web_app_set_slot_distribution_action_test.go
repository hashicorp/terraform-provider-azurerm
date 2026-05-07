package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type WebAppSetSlotDistributionAction struct{}

func TestAccWebAppSetSlotDistributionAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_set_slot_distribution", "test")
	a := WebAppSetSlotDistributionAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.twoSlots(data), // first set up the app + slots
			},
			{
				Config: a.basicSingleRule(data), // apply the action with a rule
			},
			{
				Config: a.removeAllRules(data), // apply action removing rules, required for `destroy` success
			},
		},
	})
}

func TestAccWebAppSetSlotDistributionAction_multiRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_web_app_set_slot_distribution", "test")
	a := WebAppSetSlotDistributionAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.twoSlots(data), // first set up the app + slots
			},
			{
				Config: a.changeStepMultiRule(data), // apply the action with multiple rules
			},
			{
				Config: a.removeAllRules(data), // apply action removing rules, required for `destroy` success
			},
		},
	})
}

func (a WebAppSetSlotDistributionAction) removeAllRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "terraform_data" "action_remove_trigger" {
  input = "only_to_trigger_remove_action"
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_web_app_set_slot_distribution.remove]
    }
  }
}

action "azurerm_web_app_set_slot_distribution" "remove" {
  config {
    app_service_id = azurerm_linux_web_app.test.id
  }
}
`, a.twoSlots(data))
}

func (a WebAppSetSlotDistributionAction) basicSingleRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "terraform_data" "action_set_trigger" {
  input = "only_to_trigger_set_action"
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_web_app_set_slot_distribution.test]
    }
  }
}

action "azurerm_web_app_set_slot_distribution" "test" {
  config {
    app_service_id = azurerm_linux_web_app.test.id
    slot_rule {
      hostname           = azurerm_linux_web_app_slot.test1.default_hostname
      rule_name          = azurerm_linux_web_app_slot.test1.name
      reroute_percentage = 10
    }
  }
}
`, a.twoSlots(data))
}

func (a WebAppSetSlotDistributionAction) changeStepMultiRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "terraform_data" "action_set_trigger" {
  input = "only_to_trigger_set_action"
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_web_app_set_slot_distribution.test]
    }
  }
}

action "azurerm_web_app_set_slot_distribution" "test" {
  config {
    app_service_id = azurerm_linux_web_app.test.id
    slot_rule {
      hostname                = azurerm_linux_web_app_slot.test1.default_hostname
      rule_name               = azurerm_linux_web_app_slot.test1.name
      reroute_percentage      = 10
      change_step             = 1
      change_interval_minutes = 1
      min_reroute_percentage  = 5
      max_reroute_percentage  = 20
    }
    slot_rule {
      hostname           = azurerm_linux_web_app_slot.test2.default_hostname
      rule_name          = azurerm_linux_web_app_slot.test2.name
      reroute_percentage = 2.2
    }
  }
}
`, a.twoSlots(data))
}

func (a WebAppSetSlotDistributionAction) twoSlots(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test1" {
  name           = "acctestWAS-1-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}
}

resource "azurerm_linux_web_app_slot" "test2" {
  name           = "acctestWAS-2-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}
}
`, a.template(data), data.RandomInteger)
}

func (WebAppSetSlotDistributionAction) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-slotdist-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-WAS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary)
}
