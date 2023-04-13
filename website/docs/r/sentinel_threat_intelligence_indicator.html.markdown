---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_threat_intelligence_indicator"
description: |-
  Manages a Sentinel Threat Intelligence Indicator.
---

# azurerm_sentinel_threat_intelligence_indicator

Manages a Sentinel Threat Intelligence Indicator.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "east us"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-law"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "example" {
  resource_group_name = azurerm_resource_group.example.name
  workspace_name      = azurerm_log_analytics_workspace.example.name
}

resource "azurerm_sentinel_threat_intelligence_indicator" "example" {
  workspace_id      = azurerm_log_analytics_workspace.example.id
  pattern_type      = "domain-name"
  pattern           = "http://example.com"
  source            = "Microsoft Sentinel"
  validate_from_utc = "2022-12-14T16:00:00Z"
  display_name      = "example-indicator"

  depends_on = [azurerm_sentinel_log_analytics_workspace_onboarding.test]
}
```

## Arguments Reference

The following arguments are supported:

* `display_name` - (Required) The display name of the Threat Intelligence Indicator.

* `pattern_type` - (Required) The type of pattern used by the Threat Intelligence Indicator. Possible values are `domain-name`, `file`, `ipv4-addr`, `ipv6-addr` and `url`.

* `pattern` - (Required) The pattern used by the Threat Intelligence Indicator. When `pattern_type` set to `file`, `pattern` must be specified with `<HashName>:<Value>` format, such as `MD5:78ecc5c05cd8b79af480df2f8fba0b9d`.

* `source` - (Required) Source of the Threat Intelligence Indicator.

* `validate_from_utc` - (Required) The start of validate date in RFC3339.

* `workspace_id` - (Required) The ID of the Log Analytics Workspace. Changing this forces a new Sentinel Threat Intelligence Indicator to be created.

---

* `confidence` - (Optional) Confidence levels of the Threat Intelligence Indicator.

* `created_by` - (Optional) The creator of the Threat Intelligence Indicator.

* `description` - (Optional) The description of the Threat Intelligence Indicator.

* `extension` - (Optional) The extension config of the Threat Intelligence Indicator in JSON format.

* `external_reference` - (Optional) One or more `external_reference` blocks as defined below.

* `granular_marking` - (Optional) One or more `granular_marking` blocks as defined below.

* `kill_chain_phase` - (Optional) One or more `kill_chain_phase` blocks as defined below.

* `tags` - (Optional) Specifies a list of tags of the Threat Intelligence Indicator.

* `language` - (Optional) The language of the Threat Intelligence Indicator.

* `modified_by` - (Optional) The user or service principal who modified the Threat Intelligence Indicator.

* `object_marking_refs` - (Optional) Specifies a list of Threat Intelligence marking references.

* `pattern_version` - (Optional) The version of a Threat Intelligence entity.

* `revoked` - (Optional) Whether the Threat Intelligence entity revoked.

* `threat_types` - (Optional) Specifies a list of threat types of this Threat Intelligence Indicator.

* `validate_until_utc` - (Optional) The end of validate date of the Threat Intelligence Indicator in RFC3339 format.

---

A `external_reference` block supports the following:

* `description` - (Optional) The description of the external reference of the Threat Intelligence Indicator.

* `hashes` - (Optional) The list of hashes of the external reference of the Threat Intelligence Indicator.

* `source_name` - (Optional) The source name of the external reference of the Threat Intelligence Indicator.

* `url` - (Optional) The url of the external reference of the Threat Intelligence Indicator.

---

A `granular_marking` block supports the following:

* `language` - (Optional) The language of granular marking of the Threat Intelligence Indicator.

* `marking_ref` - (Optional) The reference of the granular marking of the Threat Intelligence Indicator.

* `selectors` - (Optional) A list of selectors of the granular marking of the Threat Intelligence Indicator.

---

A `kill_chain_phase` block supports the following:

* `name` - (Optional) The name which should be used for the Lockheed Martin cyber kill chain phase.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Sentinel Threat Intelligence Indicator.

* `created_on` - The date of this Threat Intelligence Indicator created.

* `defanged` - Whether the Threat Intelligence entity is defanged?

* `external_id` - The external ID of the Threat Intelligence Indicator.

* `external_last_updated_time_utc` - the External last updated time in UTC.

* `indicator_types` - A list of indicator types of this Threat Intelligence Indicator.

* `last_updated_time_utc` - The last updated time of the Threat Intelligence Indicator in UTC.

* `guid` - The guid of this Sentinel Threat Intelligence Indicator.

* `parsed_pattern` - A `parsed_pattern` block as defined below.

---

A `parsed_pattern` block exports the following:

* `pattern_type_key` - The type key of parsed pattern.

* `pattern_type_values` - A `pattern_type_values` block as defined below.

---

A `pattern_type_values` block exports the following:

* `value` - The value of the parsed pattern type.

* `value_type` - The type of the value of the parsed pattern type value.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel Threat Intelligence Indicator.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Threat Intelligence Indicator.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel Threat Intelligence Indicator.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel Threat Intelligence Indicator.

## Import

Sentinel Threat Intelligence Indicators can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_threat_intelligence_indicator.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/threatIntelligence/main/indicators/indicator1
```
