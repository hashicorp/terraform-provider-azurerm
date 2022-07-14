---
subcategory: "Video Analyzer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_video_analyzer_edge_module"
description: |-
  Manages a Video Analyzer Edge Module.
---

# azurerm_video_analyzer_edge_module

Manages a Video Analyzer Edge Module.

!> Video Analyzer (Preview) is now Deprecated and will be Retired on 2022-11-30 - as such the `azurerm_video_analyzer_edge_module` resource is deprecated and will be removed in v4.0 of the AzureRM Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "video-analyzer-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "exampleidentity"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_role_assignment" "contributor" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_role_assignment" "reader" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_video_analyzer" "example" {
  name                = "exampleanalyzer"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  storage_account {
    id                        = azurerm_storage_account.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id
    ]
  }

  tags = {
    environment = "staging"
  }

  depends_on = [
    azurerm_user_assigned_identity.example,
    azurerm_role_assignment.contributor,
    azurerm_role_assignment.reader,
  ]
}

resource "azurerm_video_analyzer_edge_module" "example" {
  name                = "example-edge-module"
  resource_group_name = azurerm_resource_group.example.name
  video_analyzer_name = azurerm_video_analyzer.example.name
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Video Analyzer Edge Module. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Video Analyzer Edge Module. Changing this forces a new resource to be created.

* `video_analyzer_name` - (Required) The name of the Video Analyzer in which to create the Edge Module. Changing this forces a new resource to be created.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Video Analyzer Edge Module.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Video Analyzer Edge Module.
* `read` - (Defaults to 5 minutes) Used when retrieving the Video Analyzer Edge Module.
* `delete` - (Defaults to 30 minutes) Used when deleting the Video Analyzer Edge Module.

## Import

Video Analyzer Edge Module can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_video_analyzer_edge_module.edge /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Media/videoAnalyzers/analyzer1/edgeModules/edge1
```
