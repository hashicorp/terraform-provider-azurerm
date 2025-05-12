---
subcategory: "Search"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_search_service"
description: |-
  Manages a Search Service.
---

# azurerm_search_service

Manages a Search Service.

## Example Usage (supporting API Keys)

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

## Example Usage (using both AzureAD and API Keys)

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

  local_authentication_enabled = true
  authentication_failure_mode  = "http403"
}
```

## Example Usage (supporting only AzureAD Authentication)

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

  local_authentication_enabled = false
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Search Service should exist. Changing this forces a new Search Service to be created.

* `name` - (Required) The Name which should be used for this Search Service. Changing this forces a new Search Service to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Search Service should exist. Changing this forces a new Search Service to be created.

* `sku` - (Required) The SKU which should be used for this Search Service. Possible values include `basic`, `free`, `standard`, `standard2`, `standard3`, `storage_optimized_l1` and `storage_optimized_l2`. Changing this forces a new Search Service to be created.

-> **Note:** The `basic` and `free` SKUs provision the Search Service in a Shared Cluster - the `standard` SKUs use a Dedicated Cluster.

~> **Note:** The SKUs `standard2`, `standard3`, `storage_optimized_l1` and `storage_optimized_l2` are only available by submitting a quota increase request to Microsoft. Please see the [product documentation](https://learn.microsoft.com/azure/azure-resource-manager/troubleshooting/error-resource-quota?tabs=azure-cli) on how to submit a quota increase request.

---

* `allowed_ips` - (Optional) Specifies a list of inbound IPv4 or CIDRs that are allowed to access the Search Service. If the incoming IP request is from an IP address which is not included in the `allowed_ips` it will be blocked by the Search Services firewall.

-> **Note:** The `allowed_ips` are only applied if the `public_network_access_enabled` field has been set to `true`, else all traffic over the public interface will be rejected, even if the `allowed_ips` field has been defined. When the `public_network_access_enabled` field has been set to `false` the private endpoint connections are the only allowed access point to the Search Service.

* `authentication_failure_mode` - (Optional) Specifies the response that the Search Service should return for requests that fail authentication. Possible values include `http401WithBearerChallenge` or `http403`.

-> **Note:** `authentication_failure_mode` can only be configured when using `local_authentication_enabled` is set to `true` - which when set together specifies that both API Keys and AzureAD Authentication should be supported.

* `customer_managed_key_enforcement_enabled` - (Optional) Specifies whether the Search Service should enforce that non-customer resources are encrypted. Defaults to `false`.

* `hosting_mode` - (Optional) Specifies the Hosting Mode, which allows for High Density partitions (that allow for up to 1000 indexes) should be supported. Possible values are `highDensity` or `default`. Defaults to `default`. Changing this forces a new Search Service to be created.

-> **Note:** `hosting_mode` can only be configured when `sku` is set to `standard3`.


* `identity` - (Optional) An `identity` block as defined below.

* `local_authentication_enabled` - (Optional) Specifies whether the Search Service allows authenticating using API Keys? Defaults to `true`.

* `network_rule_bypass_option` - (Optional) Whether to allow trusted Azure services to access a network restricted Container Registry? Possible values are `None` and `AzureServices`. Defaults to `None`.

* `partition_count` - (Optional) Specifies the number of partitions which should be created. This field cannot be set when using a `free` sku ([see the Microsoft documentation](https://learn.microsoft.com/azure/search/search-sku-tier)). Possible values include `1`, `2`, `3`, `4`, `6`, or `12`. Defaults to `1`.

-> **Note:** when `hosting_mode` is set to `highDensity` the maximum number of partitions allowed is `3`.

* `public_network_access_enabled` - (Optional) Specifies whether Public Network Access is allowed for this resource. Defaults to `true`.

* `replica_count` - (Optional) Specifies the number of Replica's which should be created for this Search Service. This field cannot be set when using a `free` sku ([see the Microsoft documentation](https://learn.microsoft.com/azure/search/search-sku-tier)).

* `semantic_search_sku` - (Optional) Specifies the Semantic Search SKU which should be used for this Search Service. Possible values include `free` and `standard`.

~> **Note:** The `semantic_search_sku` cannot be defined if your Search Services `sku` is set to `free`. The Semantic Search feature is only available in certain regions, please see the [product documentation](https://learn.microsoft.com/azure/search/semantic-search-overview#availability-and-pricing) for more information.

* `tags` - (Optional) Specifies a mapping of tags which should be assigned to this Search Service.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Search Service. Possible values are `SystemAssigned`, `UserAssigned`, and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) Specifies the list of User Assigned Managed Service Identity IDs which should be assigned to this Search Service.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Search Service.

* `customer_managed_key_encryption_compliance_status` - Describes whether the search service is compliant or not with respect to having non-customer encrypted resources. If a service has more than one non-customer encrypted resource and `Enforcement` is `enabled` then the service will be marked as `NonCompliant`. If all the resources are customer encrypted, then the service will be marked as `Compliant`.

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

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Search Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Search Service.
* `update` - (Defaults to 1 hour) Used when updating the Search Service.
* `delete` - (Defaults to 1 hour) Used when deleting the Search Service.

## Import

Search Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_search_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Search/searchServices/service1
```
