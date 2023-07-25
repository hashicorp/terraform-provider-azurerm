---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_private_link_scoped_service"
description: |-
  Manages an Azure Monitor Private Link Scoped Service
---

# azurerm_monitor_private_link_scoped_service

Manages an Azure Monitor Private Link Scoped Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_monitor_private_link_scope" "example" {
  name                = "example-ampls"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_monitor_private_link_scoped_service" "example" {
  name                = "example-amplsservice"
  resource_group_name = azurerm_resource_group.example.name
  scope_name          = azurerm_monitor_private_link_scope.example.name
  linked_resource_id  = azurerm_application_insights.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Monitor Private Link Scoped Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Monitor Private Link Scoped Service should exist. Changing this forces a new resource to be created.

* `scope_name` - (Required) The name of the Azure Monitor Private Link Scope. Changing this forces a new resource to be created.

* `linked_resource_id` - (Required) The ID of the linked resource. It must be the Log Analytics workspace or the Application Insights component or the Data Collection endpoint. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Monitor Private Link Scoped Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Monitor Private Link Scope.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Monitor Private Link Scope.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Monitor Private Link Scope.
* `update` - (Defaults to 30 minutes) Used when updating the Monitor Private Link Scoped Service.

## Import

Azure Monitor Private Link Scoped Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_private_link_scoped_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/privateLinkScopes/pls1/scopedResources/sr1
```
