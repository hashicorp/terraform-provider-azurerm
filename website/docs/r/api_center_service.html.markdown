---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_service"
subcategory: "ApiCenter"
description: |-
  Documentation for azurerm_api_center_service.
---

# azurerm_api_center_service

Documentation for azurerm_api_center_service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_api_center_service" "example" {
  name                = "example-service"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the API Center service. Changing this value forces a new resource to be created.

* `resource_group_name` - (Required) No description provided. Changing this value forces a new resource to be created.

* `location` - (Required) No description provided. Changing this value forces a new resource to be created.

* `identity` - (Optional) (see [nested schema](#schema-6e6573746564-696e707574-6964656e74697479))

* `sku_name` - (Optional) The SKU name for the API Center service. Possible values include `Free` and `Standard`. Defaults to `Free`.

* `tags` - (Optional) No description provided.

---

<a id="schema-6e6573746564-696e707574-6964656e74697479"></a>
The `identity` block supports the following:

* `type` - (Required) No description provided. Possible values include `UserAssigned`, `SystemAssigned`, and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) No description provided.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Resource

* `identity` - (see [nested schema](#schema-6e6573746564-6f7574707574-6964656e74697479))

---

<a id="schema-6e6573746564-6f7574707574-6964656e74697479"></a>
The `identity` block exports the following:

* `principal_id` - No description provided.

* `tenant_id` - No description provided.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

API Center Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_center_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiCenter/services/example-service
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiCenter` - 2024-03-01
