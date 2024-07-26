---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azure_ai_services"
description: |-
  Manages Azure AI Services.
---

# azurerm_azure_ai_services

Manages Azure AI Services.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_azure_ai_services" "example" {
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

* `name` - (Required) Specifies the name of the Azure AI Services. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Azure AI Services is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this Azure AI Services. Possible values are `F0`, `F1`, `S0`, `S`, `S1`, `S2`, `S3`, `S4`, `S5`, `S6`, `P0`, `P1`, `P2`, `E0` and `DC0`.

-> **NOTE:** SKU `DC0` is the commitment tier for Azure AI Services containers running in disconnected environments. You must obtain approval from Microsoft by submitting the [request form](https://aka.ms/csdisconnectedcontainers) first, before you can use this SKU. More information on [Purchase a commitment plan to use containers in disconnected environments](https://learn.microsoft.com/en-us/azure/cognitive-services/containers/disconnected-containers?tabs=stt#purchase-a-commitment-plan-to-use-containers-in-disconnected-environments).

* `custom_subdomain_name` - (Optional) The subdomain name used for token-based authentication. This property is required when `network_acls` is specified. Changing this forces a new resource to be created.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as documented below.

* `fqdns` - (Optional) List of FQDNs allowed for the Azure AI Services.

* `identity` - (Optional) An `identity` block as defined below.

* `local_authentication_enabled` - (Optional) Whether local authentication methods is enabled for the Azure AI Services. Defaults to `true`.

* `network_acls` - (Optional) A `network_acls` block as defined below. When this property is specified, `custom_subdomain_name` is also required to be set.

* `outbound_network_access_restricted` - (Optional) Whether outbound network access is restricted for the Azure AI Services. Defaults to `false`.

* `public_network_access` - (Optional) Whether public network access is allowed for the Azure AI Services. Possible values are `Enabled` and `Disabled`. Defaults to `Enabled`.

* `storage` - (Optional) A `storage` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `network_acls` block supports the following:

* `default_action` - (Required) The Default Action to use when no rules match from `ip_rules` / `virtual_network_rules`. Possible values are `Allow` and `Deny`.

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the Azure AI Services.

* `virtual_network_rules` - (Optional) A `virtual_network_rules` block as defined below.

---

A `virtual_network_rules` block supports the following:

* `subnet_id` - (Required) The ID of the subnet which should be able to access this Azure AI Services.

* `ignore_missing_vnet_service_endpoint` - (Optional) Whether ignore missing vnet service endpoint or not. Default to `false`.

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Optional) The ID of the Key Vault Key which should be used to Encrypt the data in this Azure AI Services. Exactly one of `key_vault_key_id`, `managed_hsm_key_id` must be specified.

* `managed_hsm_key_id` - (Optional) The ID of the managed HSM Key which should be used to Encrypt the data in this Azure AI Services. Exactly one of `key_vault_key_id`, `managed_hsm_key_id` must be specified.

* `identity_client_id` - (Optional) The Client ID of the User Assigned Identity that has access to the key. This property only needs to be specified when there're multiple identities attached to the Azure AI Service.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Azure AI Services. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Azure AI Services.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `storage` block supports the following:

* `storage_account_id` - (Required) Full resource id of a Microsoft.Storage resource.

* `identity_client_id` - (Optional) The client ID of the managed identity associated with the storage resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure AI Services.

* `endpoint` - The endpoint used to connect to the Azure AI Services.

* `identity` - An `identity` block as defined below.

* `primary_access_key` - A primary access key which can be used to connect to the Azure AI Services.

* `secondary_access_key` - The secondary access key which can be used to connect to the Azure AI Services.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure AI Services.
* `update` - (Defaults to 30 minutes) Used when updating the Azure AI Services.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure AI Services.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure AI Services.

## Import

Azure AI Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_azure_ai_services.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1
```
