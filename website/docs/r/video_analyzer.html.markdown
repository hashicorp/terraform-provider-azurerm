---
subcategory: "Video Analyzer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_video_analyzer"
description: |-
  Manages a Video Analyzer.
---

# azurerm_video_analyzer

Manages a Video Analyzer.

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
    azurerm_role_assignment.contributor,
    azurerm_role_assignment.reader,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Video Analyzer. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Video Analyzer. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `storage_account` - (Required) An `storage_account` blocks is defined below.

* `identity` - (Required) An `identity` block is defined below.

* `tags` - (Optional) A mapping of tags assigned to the resource.

---

A `storage_account` block supports the following:

* `id` - (Required) Specifies the ID of the Storage Account that will be associated with the Video Analyzer instance.

* `user_assigned_identity_id` - (Required) Specifies the User Assigned Identity ID which should be assigned to a access this Storage Account.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Identity that should be configured on this Video Analyzer. At this time the only possible value is `UserAssigned`.

* `identity_ids` - (Required) A list of User Assigned Identity ID's to be assigned to this Video Analyzer.

---


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Video Analyzer.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Video Analyzer.
* `update` - (Defaults to 30 minutes) Used when updating the Video Analyzer.
* `read` - (Defaults to 5 minutes) Used when retrieving the Video Analyzer.
* `delete` - (Defaults to 30 minutes) Used when deleting the Video Analyzer.

## Import

Video Analyzer can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_video_analyzer.analyzer /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Media/videoAnalyzers/analyzer1
```
