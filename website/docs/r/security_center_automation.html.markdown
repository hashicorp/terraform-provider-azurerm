---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_automation"
description: |-
  Manages Security Center Automation and Continuous Export.
---

# azurerm_security_center_automation

Manages Security Center Automation and Continuous Export. This resource supports three types of destination in the `action`, Logic Apps, Log Analytics and Event Hubs

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
  capacity            = 2
}

resource "azurerm_eventhub" "example" {
  name                = "acceptanceTestEventHub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 2
}

resource "azurerm_eventhub_authorization_rule" "example" {
  name                = "example-rule"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name
  listen              = true
  send                = false
  manage              = false
}

resource "azurerm_security_center_automation" "example" {
  name                = "example-automation"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  action {
    type              = "EventHub"
    resource_id       = azurerm_eventhub.example.id
    connection_string = azurerm_eventhub_authorization_rule.example.primary_connection_string
  }

  source {
    event_source = "Alerts"
    rule_set {
      rule {
        property_path  = "properties.metadata.severity"
        operator       = "Equals"
        expected_value = "High"
        property_type  = "String"
      }
    }
  }

  scopes = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}"]
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Security Center Automation should exist. Changing this forces a new Security Center Automation to be created.

* `name` - (Required) The name which should be used for this Security Center Automation. Changing this forces a new Security Center Automation to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Security Center Automation should exist. Changing this forces a new Security Center Automation to be created.

* `scopes` - (Required) A list of scopes on which the automation logic is applied, at least one is required. Supported scopes are a subscription (in this format `/subscriptions/00000000-0000-0000-0000-000000000000`) or a resource group under that subscription (in the format `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example`). The automation will only apply on defined scopes.

* `source` - (Required) One or more `source` blocks as defined below. A `source` defines what data types will be processed and a set of rules to filter that data.

* `action` - (Required) One or more `action` blocks as defined below. An `action` tells this automation where the data is to be sent to upon being evaluated by the rules in the `source`.

---

* `description` - (Optional) Specifies the description for the Security Center Automation.

* `enabled` - (Optional) Boolean to enable or disable this Security Center Automation.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

A `action` block defines where the data will be exported and sent to, it supports the following:

* `type` - (Required) Type of Azure resource to send data to. Must be set to one of: `LogicApp`, `EventHub` or `LogAnalytics`.

* `resource_id` - (Required) The resource id of the target Logic App, Event Hub namespace or Log Analytics workspace.

* `connection_string` - (Optional, but required when `type` is `EventHub`) A connection string to send data to the target Event Hub namespace, this should include a key with send permissions.

* `trigger_url` - (Optional, but required when `type` is `LogicApp`) The callback URL to trigger the Logic App that will receive and process data sent by this automation. This can be found in the Azure Portal under "See trigger history"

---

A `source` block defines the source data in Security Center to be exported, supports the following:

* `event_source` - (Required) Type of data that will trigger this automation. Must be one of `Alerts`, `Assessments`, `SecureScoreControls`, `SecureScores` or `SubAssessments`. Note. assessments are also referred to as recommendations 

* `rule_set` - (Optional) A set of rules which evaluate upon event and data interception. This is defined in one or more `rule_set` blocks as defined below.
  
~> **NOTE:** When multiple `rule_set` block are provided, a logical 'OR' is applied to the evaluation of them.

---

A `rule_set` block supports the following:

* `rule` - (Required) One or more `rule` blocks as defined below. 

~> **NOTE:** This automation will trigger when all of the `rule`s in this `rule_set` are evaluated as 'true'. This is equivalent to a logical 'AND'.

---

A `rule` block supports the following:

* `expected_value` - (Required) A value that will be compared with the value in `property_path`.

* `operator` - (Required) The comparison operator to use, must be one of: `Contains`, `EndsWith`, `Equals`, `GreaterThan`, `GreaterThanOrEqualTo`, `LesserThan`, `LesserThanOrEqualTo`, `NotEquals`, `StartsWith`

* `property_path` - (Required) The JPath of the entity model property that should be checked.

* `property_type` - (Required) The data type of the compared operands, must be one of: `Integer`, `String`, `Boolean` or `Number`.

~> **NOTE:** The schema for Security Center alerts (when `event_source` is "Alerts") [can be found here](https://docs.microsoft.com/en-us/azure/security-center/alerts-schemas?tabs=schema-continuousexport)


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Security Center Automation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Security Center Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Automation.
* `update` - (Defaults to 30 minutes) Used when updating the Security Center Automation.
* `delete` - (Defaults to 30 minutes) Used when deleting the Security Center Automation.

## Import

Security Center Automations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_automation.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Security/automations/automation1
```
