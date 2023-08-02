---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_registry"
description: |-
  Manages a Machine Learning Registry.
---

# azurerm_machine_learning_registry

Manages a Machine Learning Registry.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "WestEurope"
}


resource "azurerm_machine_learning_registry" "example" {
  location            = "WestEurope"
  resource_group_name = azurerm_resource_group.example.name
  name                = "mlr-example"

  identity {
    type = "SystemAssigned"
  }

  public_network_access_enabled = true

  region_details {
    location = "WestEurope"
    storage_account_details {
      system_created_storage_account {
        storage_account_type = "Standard_LRS"
      }
    }

    acr_details {
      system_created_acr_account {
        acr_account_sku = "Premium"
      }
    }
  }

  region_details {
    location = "EastUS2"
    storage_account_details {
      system_created_storage_account {
        storage_account_type = "Standard_LRS"
      }
    }

    acr_details {
      system_created_acr_account {
        acr_account_sku = "Premium"
      }
    }
  }

  tags = {
    key = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created.

* `name` - (Required) The name which should be used for this Machine Learning Registry. Changing this forces a new Machine Learning Registry to be created.

* `region_details` - (Required) One or more `region_details` blocks as defined below. Changing this forces a new Machine Learning Registry to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created.

---

* `identity` - (Optional) A `identity` block as defined below. Changing this forces a new Machine Learning Registry to be created.

* `public_network_access_enabled` - (Optional) Should the public network access be enabled? Defaults to `false`. Changing this forces a new Machine Learning Registry to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Registry.

---

A `acr_details` block supports the following:

* `system_created_acr_account` - (Optional) One or more `system_created_acr_account` blocks as defined below.

* `user_created_acr_account` - (Optional) One or more `user_created_acr_account` blocks as defined below.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible values are `SystemAssigned`, `UserAssigned`. Changing this forces a new resource to be created.

* `identity_ids` - (Optional) Specifies the list of User Assigned Managed Service Identity IDs which should be assigned to this Machine Learning Registry. Changing this forces a new Machine Learning Registry to be created.

---

A `region_details` block supports the following:

* `location` - (Required) The Azure Region where the Machine Learning Registry should exist.

* `acr_details` - (Optional) One or more `acr_details` blocks as defined above.

* `storage_account_details` - (Optional) One or more `storage_account_details` blocks as defined below.

---

A `storage_account_details` block supports the following:

* `system_created_storage_account` - (Optional) One or more `system_created_storage_account` blocks as defined below.

* `user_created_storage_account` - (Optional) One or more `user_created_storage_account` blocks as defined below.

---

A `system_created_acr_account` block supports the following:

* `acr_account_sku` - (Optional) The SKU of the ACR account. The only possible value is `Premium`.

---

A `system_created_storage_account` block supports the following:

* `blob_public_access_allowed` - (Optional) Whether public blob access allowed. Defaults to `false`.

* `storage_account_hns_enabled` - (Optional) Should the HNS be enabled for storage account? Defaults to `false`.

* `storage_account_type` - (Optional) Specify the storage account type. Possible values are `Standard_LRS`, `Standard_GRS`, `Standard_RAGRS`, `Standard_ZRS`, `Standard_GZRS`, `Standard_RAGZRS`, `Premium_LRS` and `Premium_ZRS`.

---

A `user_created_acr_account` block supports the following:

* `arm_resource_id` - (Optional) The ID of the user created ACR account.

---

A `user_created_storage_account` block supports the following:

* `arm_resource_id` - (Optional) The ID of the user created ACR account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning Registry.

* `discovery_url` - The discovery URL for this Registry.

* `intellectual_property_publisher` - The intellectual property publisher.

* `managed_resource_group` -The resource ID of the managed resource group if the Registry has system created resources.

* `ml_flow_registry_uri` - The MLFlow registry URI for this registry.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Registry.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Registry.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Registry.
* `delete` - (Defaults to 10 minutes) Used when deleting the Machine Learning Registry.

## Import

Machine Learning Registries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_registry.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.MachineLearningServices/registries/reg1
```
