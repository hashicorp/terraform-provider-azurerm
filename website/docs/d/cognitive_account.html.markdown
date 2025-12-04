---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account"
description: |-
  Gets information about an existing Cognitive Services Account.
---

# Data Source: azurerm_cognitive_account

Use this data source to access information about an existing Cognitive Services Account.

## Example Usage

```hcl
data "azurerm_cognitive_account" "test" {
  name                = "example-account"
  resource_group_name = "cognitive_account_rg"
}

output "primary_access_key" {
  value = data.azurerm_cognitive_account.test.primary_access_key
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cognitive Services Account.

* `resource_group_name` - (Required) Specifies the name of the resource group where the Cognitive Services Account resides.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cognitive Services Account.

* `custom_question_answering_search_service_id` - The ID of the search service.

* `custom_subdomain_name` - The subdomain name used for Entra ID token-based authentication.

* `customer_managed_key` - A `customer_managed_key` block as defined below.

* `dynamic_throttling_enabled` - Whether dynamic throttling is enabled for this Cognitive Services Account.

* `endpoint` - The endpoint of the Cognitive Services Account.

* `fqdns` - List of FQDNs allowed for the Cognitive Services Account.

* `identity` - A `identity` block as defined below.

* `kind` - The type of the Cognitive Services Account.

* `local_auth_enabled` - Whether local authentication methods are enabled for the Cognitive Services Account.

* `location` - The Azure location where the Cognitive Services Account exists.

* `metrics_advisor_aad_client_id` - The Microsoft Entra Application (client) ID.

* `metrics_advisor_aad_tenant_id` - The Microsoft Entra Tenant ID.

* `metrics_advisor_super_user_name` - The super user of Metrics Advisor.

* `metrics_advisor_website_name` - The website name of Metrics Advisor.

* `network_acls` -  A `network_acls` block as defined below.

* `network_injection` -  A `network_injection` block as defined below.

* `outbound_network_access_restricted` - Whether outbound network access is restricted for the Cognitive Services Account.

* `project_management_enabled` -  Whether project management is enabled.

* `public_network_access_enabled` - Whether public network access is allowed for the Cognitive Services Account.

* `qna_runtime_endpoint` - The link to the QNA runtime.

* `sku_name` - The SKU name of the Cognitive Services Account.

* `storage` - A `storage` block as defined below.

* `primary_access_key` - The primary access key of the Cognitive Services Account.

* `secondary_access_key` - The secondary access key of the Cognitive Services Account.

* `tags` - A mapping of tags to assigned to the resource.

---

A `customer_managed_key` block exports the following:

* `key_vault_key_id` - The ID of the Key Vault Key which is used to encrypt the data in this Cognitive Services Account.

* `identity_client_id` - The Client ID of the User Assigned Identity that has access to the key.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Cognitive Services Account.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Cognitive Services Account.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Cognitive Services Account.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Cognitive Services Account.

---

A `network_acls` block exports the following:

* `bypass` - Whether trusted Azure Services are allowed to access the service.

* `default_action` - The Default Action to use when no rules match from `ip_rules` / `virtual_network_rules`.

* `ip_rules` - One or more IP Addresses, or CIDR Blocks that are able to access the Cognitive Services Account.

* `virtual_network_rules` - A `virtual_network_rules` block as defined below.

---

A `network_injection` block exports the following:

* `scenario` - The feature that network injection is applied to.

* `subnet_id` - The ID of the subnet which the Agent Client is injected into.

---

A `virtual_network_rules` block exports the following:

* `subnet_id` - The ID of the subnet which is able to access this Cognitive Services Account.

* `ignore_missing_vnet_service_endpoint` - Whether missing vnet service endpoint is ignored or not.

---

A `storage` block exports the following:

* `storage_account_id` - The ID of the Storage Account resource associated with this Cognitive Services Account.

* `identity_client_id` - The client ID of the managed identity associated with the storage resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Services Account.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2025-06-01
