---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_threat_intelligence_taxii"
description: |-
  Manages a Threat Intelligence TAXII Data Connector.
---

# azurerm_sentinel_data_connector_threat_intelligence_taxii

Manages a Threat Intelligence TAXII Data Connector.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_data_connector_threat_intelligence_taxii" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
  display_name               = "example"
  taxii_server_api_root      = "https://limo.anomali.com/api/v1/taxii2/feeds"
  taxii_server_collection_id = 107
  taxii_server_username      = "guest"
  taxii_server_password      = "guest"
}
```

## Arguments Reference

The following arguments are supported:

* `display_name` - (Required) The friendly name of this Threat Intelligence TAXII Data Connector.

- `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this Threat Intelligence TAXII Data Connector resides in. Changing this forces a new Threat Intelligence TAXII Data Connector to be created.

* `name` - (Required) The name which should be used for this Threat Intelligence TAXII Data Connector. Changing this forces a new Threat Intelligence TAXII Data Connector to be created.

* `taxii_server_api_root` - (Required) The API root URL of the TAXII server that this Threat Intelligence TAXII Data Connector connects to. Changing this forces a new Threat Intelligence TAXII Data Connector to be created.

* `taxii_server_collection_id` - (Required) The ID of the TAXII server collection that this Threat Intelligence TAXII Data Connector connects to. Changing this forces a new Threat Intelligence TAXII Data Connector to be created.

---

* `taxii_server_password` - (Optional) The password of the TAXII server collection that this Threat Intelligence TAXII Data Connector connects to.

* `taxii_server_username` - (Optional) The username of the TAXII server collection that this Threat Intelligence TAXII Data Connector connects to.

- `tenant_id` - (Optional) The ID of the tenant that this Threat Intelligence TAXII Data Connector connects to. Changing this forces a new Threat Intelligence TAXII Data Connector to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Threat Intelligence TAXII Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Threat Intelligence TAXII Data Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Threat Intelligence TAXII Data Connector.
* `update` - (Defaults to 30 minutes) Used when updating the Threat Intelligence TAXII Data Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Threat Intelligence TAXII Data Connector.

## Import

Threat Intelligence TAXII Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_threat_intelligence_taxii.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
