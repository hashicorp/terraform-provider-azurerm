---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_credential_user_managed_identity"
description: |-
  Manage a Data Factory User Assigned Managed Identity credential resource
---

# azurerm_data_factory_credential_user_managed_identity

Manage a Data Factory User Assigned Managed Identity credential resource. These resources are used by Data Factory to access data sources.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "westus"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
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

* `name` - (Required) Specifies the name of the Credential. Changing this forces a new resource to be created.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Credential with. Changing this forces a new resource.

* `identity_id` - (Required) The Resouce ID of an existing User Assigned Managed Identity. This can be changed without recreating the resource. Changing this forces a new resource to be created.

~> **Note:** Attempting to create a Credential resource without first assigning the identity to the parent Data Factory will result in an Azure API error.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Credential.

~> **Note:** Manually altering a Credential resource will cause annotations to be lost, resulting in a change being detected on the next run.

* `description` - (Optional) The description for the Data Factory Credential.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Credential.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Data Factory Credential.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Credential.
* `update` - (Defaults to 5 minutes) Used when updating the Data Factory Credential.
* `delete` - (Defaults to 5 minutes) Used when deleting the Data Factory Credential.

## Import

Data Factory Credentials can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_credential_user_managed_identity.example /subscriptions/1f3d6e58-feed-4bb6-87e5-a52305ad3375/resourceGroups/example-resources/providers/Microsoft.DataFactory/factories/example/credentials/credential1
```
