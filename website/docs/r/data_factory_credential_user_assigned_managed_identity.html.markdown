---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_credential_user_managed_identity"
description: |-
  Manage a Data Factory User Assigned Managed Identity resource
---

# azurerm_data_factory_credential_user_managed_identity

Manage a Data Factory User Assigned Managed Identity resource

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "westus"
}

resource "azurerm_data_factory" "example" {
  name                = "bruceharrison-12334"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }
}

resource "azurerm_user_assigned_identity" "example" {
  location            = azurerm_resource_group.example.location
  name                = "my-user"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_credential_user_managed_identity" "test" {
  name            = azurerm_user_assigned_identity.example.name
  description     = "Short description of this credential"
  data_factory_id = azurerm_data_factory.example.id
  identity_id     = azurerm_user_assigned_identity.example.id

  annotations = ["example", "example2"]
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the credential. Changing this forces a new resource to be created.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Credential with. Changing this forces a new resource.

* `identity_id` - (Required) The Resouce ID of an existing User Assigned Managed Identity. This can be changed without recreating the resource.

~> **Note:** The User Assigned Managed Identity must also be assigned to the parent Data Factory. Attempting to create a Credentials resource without
 first doing this will result in an Azure API error.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Dataset.

* `description` - (Optional) The description for the Data Factory Dataset.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Dataset.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Dataset.

## Import

Data Factory Credentials can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_credential_user_managed_identity.example /subscriptions/1f3d6e58-feed-4bb6-87e5-a52305ad3375/resourceGroups/example-resources/providers/Microsoft.DataFactory/factories/bruceharrison-12334/credentials/credential1
```

HCL import blocks can also be used for import

```hcl
import {
  to = azurerm_data_factory_credential_user_managed_identity.test
  id = "/subscriptions/1f3d6e58-feed-4bb6-87e5-a52305ad3375/resourceGroups/example-resources/providers/Microsoft.DataFactory/factories/bruceharrison-12334/credentials/credential1"
}
```
