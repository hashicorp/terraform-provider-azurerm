---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_app_set_slot_distribution"
description: |-
  Configures routing traffic to Web App deployment slots.
---

# Action: azurerm_web_app_set_slot_distribution

Allows configuration of an existing Web App site to distribute traffic amongst one or more staging slots. The configuration rules can be a fixed percent, or a ramp-up profile that will adjust over time.

## Example Usage

```terraform
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_plan" "example" {
  name                = "example-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_web_app" "example" {
  name                = "example-linux-web-app"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_service_plan.example.location
  service_plan_id     = azurerm_service_plan.example.id

  site_config {}
}

resource "azurerm_linux_web_app_slot" "example_slot_1" {
  name           = "example-slot-1"
  app_service_id = azurerm_linux_web_app.example.id

  site_config {}
}

resource "azurerm_linux_web_app_slot" "example_slot_2" {
  name           = "example-slot-2"
  app_service_id = azurerm_linux_web_app.example.id

  site_config {}
}

# can be invoked from resource lifecycle trigger or `terraform -invoke` CLI
action "azurerm_web_app_set_slot_distribution" "enable_distribution" {
  config {
    app_service_id = azurerm_linux_web_app.example.id
    slot_rule {
      hostname                = azurerm_linux_web_app_slot.example_slot_1.default_hostname
      rule_name               = azurerm_linux_web_app_slot.example_slot_1.name
      reroute_percentage      = 10
    }
    slot_rule {
      hostname           = azurerm_linux_web_app_slot.example_slot_2.default_hostname
      rule_name          = azurerm_linux_web_app_slot.example_slot_2.name
      reroute_percentage = 5
    }
  }
}

# empty rules will remove all slot distribution and set production back to 100%
action "azurerm_web_app_set_slot_distribution" "remove_distribution" {
  config {
    app_service_id = azurerm_linux_web_app.example.id
  }
}
```

## Argument Reference

This action supports the following arguments:

* `app_service_id` - (Required) The ID of the Web App (Linux or Windows) this distribution will be configured for.

* `slot_rule` - (Optional) One or more `slot_rule` blocks as defined below.

* `timeout` - (Optional) Timeout duration for the action to complete. Defaults to `15m`.

---

A `slot_rule` block supports the following:

* `rule_name` - (Required) Name of the routing rule, unique within this action.

* `hostname` - (Required) Hostname of a slot to which the traffic will be redirected, unique within this action.

* `reroute_percentage` - (Required) Percentage of the traffic which will be redirected to `hostname`.

~> **Note:** The remaining optional attributes support a more advanced auto ramp-up scenario and will dynamically scale the distribution over time.

* `change_step` - (Optional) The step amount to add/remove from `reroute_percentage` until it reaches `min_reroute_percentage` or `max_reroute_percentage`.

* `change_interval_minutes` - (Optional) Specifies interval in minutes to reevaluate `reroute_percentage`, required for setting a `change_step`.

* `min_reroute_percentage` - (Optional) Specifies lower boundary above which `reroute_percentage` will stay.

* `max_reroute_percentage` - (Optional) Specifies upper boundary below which `reroute_percentage` will stay.

* `change_decision_callback_url` - (Optional) URL of a custom decision algorithm provided in TiPCallback site extension.
