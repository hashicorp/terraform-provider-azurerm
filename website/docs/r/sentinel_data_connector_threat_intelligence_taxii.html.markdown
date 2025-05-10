---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_threat_intelligence_taxii"
description: |-
  Manages an Threat Intelligence TAXII Data Connector.
---

# azurerm_sentinel_data_connector_threat_intelligence_taxii

Manages an Threat Intelligence TAXII Data Connector.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "example" {
  workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_sentinel_data_connector_threat_intelligence_taxii" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.example.workspace_id
  display_name               = "example"
  api_root_url               = "https://foo/taxii2/api2/"
  collection_id              = "someid"
}
```

## Arguments Reference

The following arguments are supported:

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this Threat Intelligence TAXII Data Connector resides in. Changing this forces a new Threat Intelligence TAXII Data Connector to be created.

* `name` - (Required) The name which should be used for this Threat Intelligence TAXII Data Connector. Changing this forces a new Threat Intelligence TAXII Data Connector to be created.

* `display_name` - (Required) The friendly name which should be used for this Threat Intelligence TAXII Data Connector.

* `api_root_url` - (Required) The API root URI of the TAXII server.

* `collection_id` - (Required) The collection ID of the TAXII server.

* `user_name` - (Optional) The user name for the TAXII server.

* `password` - (Optional) The password for the TAXII server.

* `polling_frequency` - (Optional) The polling frequency for the TAXII server. Possible values are `OnceAMinute`, `OnceAnHour` and `OnceADay`. Defaults to `OnceAnHour`.

* `lookback_date` - (Optional) The lookback date for the TAXII server in RFC3339. Defaults to `1970-01-01T00:00:00Z`.

---

* `tenant_id` - (Optional) The ID of the tenant that this Threat Intelligence TAXII Data Connector connects to. Changing this forces a new Threat Intelligence TAXII Data Connector to be created.

-> **Note:** Currently, only the same tenant as the running account is allowed. Cross-tenant scenario is not supported yet.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Threat Intelligence TAXII Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Threat Intelligence TAXII Data Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Threat Intelligence TAXII Data Connector.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel Data Connector Threat Intelligence Taxii.
* `delete` - (Defaults to 30 minutes) Used when deleting the Threat Intelligence TAXII Data Connector.

## Import

Threat Intelligence TAXII Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_threat_intelligence_taxii.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
