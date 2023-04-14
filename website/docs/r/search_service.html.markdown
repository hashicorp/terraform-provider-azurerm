---
subcategory: "Search"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_search_service"
description: |-
  Manages a Search Service.
---

# azurerm_search_service

Manages a Search Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_search_service" "example" {
  name                = "example-resource"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "standard"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Search Service should exist. Changing this forces a new Search Service to be created.

* `name` - (Required) The Name which should be used for this Search Service. Changing this forces a new Search Service to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Search Service should exist. Changing this forces a new Search Service to be created.

* `sku` - (Required) The SKU which should be used for this Search Service. Possible values include `basic`, `free`, `standard`, `standard2`, `standard3`, `storage_optimized_l1` and `storage_optimized_l2`. Changing this forces a new Search Service to be created.

-> The `basic` and `free` SKUs provision the Search Service in a Shared Cluster - the `standard` SKUs use a Dedicated Cluster.

~> **NOTE:** The SKUs `standard2`, `standard3`, `storage_optimized_l1` and `storage_optimized_l2` are only available by submitting a quota increase request to Microsoft. Please see the [product documentation](https://learn.microsoft.com/azure/azure-resource-manager/troubleshooting/error-resource-quota?tabs=azure-cli) on how to submit a quota increase request.

---

* `authentication_failure_mode` - (Optional) Describes what response the Search Service should return for requests that fail authentication. Possible values include `http401WithBearerChallenge` or `http403`.

-> **NOTE:** `authentication_failure_mode` cannot be defined if the `local_authentication_disabled` is set to `true`.

* `customer_managed_key_enforcement_enabled` - (Optional) Should the Search Service enforce having non customer encrypted resources? Possible values include `true` or `false`. If `true` the Search Service will be marked as `non-compliant` if there are one or more non customer encrypted resources, if `false` no enforcement will be made and the Search Service can contain one or more non customer encrypted resources. Defaults to `false`.

* `hosting_mode` - (Optional) Enable high density partitions that allow for up to a 1000 indexes. Possible values are `highDensity` or `default`. Defaults to `default`. Changing this forces a new Search Service to be created.

-> **NOTE:** When the Search Service is in `highDensity` mode the maximum number of partitions allowed is `3`, to enable `hosting_mode` you must use a `standard3` SKU.

* `local_authentication_disabled` - (Optional) Should tha Search Service *not* be allowed to use API keys for authentication? Possible values include `true` or `false`. Defaults to `false`.

-> **NOTE:** If the `local_authentication_disabled` field is `false` and the `authentication_failure_mode` has not been defined the Search Service will be in `API Keys Only` mode. If the `local_authentication_disabled` field is `false` and the `authentication_failure_mode` has also been set to `http401WithBearerChallenge` or `http403` the Search Service will be in `Role-based access contol and API Keys` mode (e.g. `Both`). If the `local_authentication_disabled` field is `true` the Search Service will be in `Role-based access contol Only` mode. When the `local_authentication_disabled` field is `true` the `authentication_failure_mode` cannot be defined.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this resource. Defaults to `true`.

* `partition_count` - (Optional) The number of partitions which should be created. Possible values include `1`, `2`, `3`, `4`, `6`, or `12`. Defaults to `1`.

-> **NOTE:** `partition_count` cannot be configured when using a `free` or `basic` SKU. For more information please to the [product documentation](https://learn.microsoft.com/azure/search/search-sku-tier).

* `replica_count` - (Optional) The number of replica's which should be created.

-> **NOTE:** `replica_count` cannot be configured when using a `free` SKU. For more information please to the [product documentation](https://learn.microsoft.com/azure/search/search-sku-tier).

* `allowed_ips` - (Optional) A list of inbound IPv4 or CIDRs that are allowed to access the Search Service. If the incoming IP request is from an IP address which is not included in the `allowed_ips` it will be blocked by the Search Services firewall.

-> **NOTE:** The `allowed_ips` are only applied if the `public_network_access_enabled` field has been set to `true`, else all traffic over the public interface will be rejected, even if the `allowed_ips` field has been defined. When the `public_network_access_enabled` field has been set to `false` the private endpoint connections are the only allowed access point to the Search Service.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Search Service.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Search Service. The only possible value is `SystemAssigned`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Search Service.

* `customer_managed_key_enforcement_compliance` - Describes whether the Search Service is `compliant` or not with respect to having non customer encrypted resources.

* `primary_key` - The Primary Key used for Search Service Administration.

* `query_keys` - A `query_keys` block as defined below.

* `secondary_key` - The Secondary Key used for Search Service Administration.

---

A `query_keys` block exports the following:

* `key` - The value of this Query Key.

* `name` - The name of this Query Key.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Search Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Search Service.
* `update` - (Defaults to 60 minutes) Used when updating the Search Service.
* `delete` - (Defaults to 60 minutes) Used when deleting the Search Service.

## Import

Search Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_search_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Search/searchServices/service1
```
