---
subcategory: "Load Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_load_test"
description: |-
  Manages a Load Test.
---

<!-- Note: This documentation is generated. Any manual changes will be overwritten -->

# azurerm_load_test

Manages a Load Test Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_user_assigned_identity" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
resource "azurerm_load_test" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Load Test should exist. Changing this forces a new Load Test to be created.

* `name` - (Required) Specifies the name of this Load Test. Changing this forces a new Load Test to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Load Test should exist. Changing this forces a new Load Test to be created.

* `description` - (Optional) Description of the resource.

* `identity` - (Optional) An `identity` block as defined below. Specifies the Managed Identity which should be assigned to this Load Test.

* `encryption` - (Optional) An `encryption` block as defined below. Changing this forces a new Load Test to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Load Test.

---

The `identity` block supports the following arguments:

* `type` - (Required) Specifies the type of Managed Identity that should be assigned to this Load Test. Possible values are `SystemAssigned`, `SystemAssigned, UserAssigned` and `UserAssigned`.

* `identity_ids` - (Optional) A list of the User Assigned Identity IDs that should be assigned to this Load Test.

In addition to the arguments defined above, the `identity` block exports the following attributes:

* `principal_id` - The Principal ID for the System-Assigned Managed Identity assigned to this Load Test.
* 
* `tenant_id` - The Tenant ID for the System-Assigned Managed Identity assigned to this Load Test.

---

The `encryption` block supports the following arguments:

* `key_url` - (Required) The URI specifying the Key vault and key to be used to encrypt data in this resource. The URI should include the key version. Changing this forces a new Load Test to be created.

* `identity` - (Required) An `identity` block as defined below. Changing this forces a new Load Test to be created.

---

The `identity` block for `encryption` supports the following arguments:

* `type` - (Required) Specifies the type of Managed Identity that should be assigned to this Load Test Encryption. Possible values are `SystemAssigned` or `UserAssigned`. Changing this forces a new Load Test to be created.

* `identity_id` - (Required) The User Assigned Identity ID that should be assigned to this Load Test Encryption. Changing this forces a new Load Test to be created.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Load Test.

* `data_plane_uri` - Resource data plane URI.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Load Test.
* `read` - (Defaults to 5 minutes) Used when retrieving the Load Test.
* `update` - (Defaults to 30 minutes) Used when updating the Load Test.
* `delete` - (Defaults to 30 minutes) Used when deleting the Load Test.

## Import

An existing Load Test can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_load_test.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.LoadTestService/loadTests/{loadTestName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Load Test exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Load Test exists. For example `example-resource-group`.
* Where `{loadTestName}` is the name of the Load Test. For example `loadTestValue`.
