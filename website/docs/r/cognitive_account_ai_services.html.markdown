---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_ai_services"
description: |-
  Manages a Cognitive AI Services Account.
---

# azurerm_cognitive_account_ai_services

Manages a Cognitive AI Services Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account_ai_services" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "S0"

  tags = {
    Acceptance = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cognitive AI Service Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cognitive AI Service Account is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this Cognitive AI Service Account. Possible values are `F0`, `F1`, `S0`, `S`, `S1`, `S2`, `S3`, `S4`, `S5`, `S6`, `P0`, `P1`, `P2`, `E0` and `DC0`.

-> **NOTE:** SKU `DC0` is the commitment tier for Cognitive AI Services containers running in disconnected environments. You must obtain approval from Microsoft by submitting the [request form](https://aka.ms/csdisconnectedcontainers) first, before you can use this SKU. More information on [Purchase a commitment plan to use containers in disconnected environments](https://learn.microsoft.com/en-us/azure/cognitive-services/containers/disconnected-containers?tabs=stt#purchase-a-commitment-plan-to-use-containers-in-disconnected-environments).

* `custom_subdomain_name` - (Optional) The subdomain name used for token-based authentication. This property is required when `network_acls` is specified. Changing this forces a new resource to be created.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as documented below.

* `fqdns` - (Optional) List of FQDNs allowed for the Cognitive AI Service Account.

* `identity` - (Optional) An `identity` block as defined below.

* `local_auth_enabled` - (Optional) Whether local authentication methods is enabled for the Cognitive AI Service Account. Defaults to `true`.

* `network_acls` - (Optional) A `network_acls` block as defined below. When this property is specified, `custom_subdomain_name` is also required to be set.

* `outbound_network_access_restricted` - (Optional) Whether outbound network access is restricted for the Cognitive AI Service Account. Defaults to `false`.

* `public_network_access` - (Optional) Whether public network access is allowed for the Cognitive AI Service Account. Possible values are `Enabled` and `Disabled`. Defaults to `Enabled`.

* `storage` - (Optional) A `storage` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `network_acls` block supports the following:

* `default_action` - (Required) The Default Action to use when no rules match from `ip_rules` / `virtual_network_rules`. Possible values are `Allow` and `Deny`.

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the Cognitive AI Service Account.

* `virtual_network_rules` - (Optional) A `virtual_network_rules` block as defined below.

---

A `virtual_network_rules` block supports the following:

* `subnet_id` - (Required) The ID of the subnet which should be able to access this Cognitive AI Service Account.

* `ignore_missing_vnet_service_endpoint` - (Optional) Whether ignore missing vnet service endpoint or not. Default to `false`.

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Encrypt the data in this Cognitive AI Service Account.

* `identity_client_id` - (Optional) The Client ID of the User Assigned Identity that has access to the key. This property only needs to be specified when there're multiple identities attached to the Cognitive AI Service Account.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Cognitive AI Service Account. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Cognitive AI Service Account.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `storage` block supports the following:

* `storage_account_id` - (Required) Full resource id of a Microsoft.Storage resource.

* `identity_client_id` - (Optional) The client ID of the managed identity associated with the storage resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive AI Service Account.

* `endpoint` - The endpoint used to connect to the Cognitive AI Service Account.

* `identity` - An `identity` block as defined below.

* `primary_access_key` - A primary access key which can be used to connect to the Cognitive AI Service Account.

* `secondary_access_key` - The secondary access key which can be used to connect to the Cognitive AI Service Account.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive AI Service Account.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive AI Service Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive AI Service Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive AI Service Account.

## Import

Cognitive AI Service Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_ai_services.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1
```
