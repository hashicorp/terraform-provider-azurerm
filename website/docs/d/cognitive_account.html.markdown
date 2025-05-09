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

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cognitive Services Account.

* `resource_group_name` - (Required) Specifies the name of the resource group where the Cognitive Services Account resides.

## Attributes Reference

The following attributes are exported:

* `identity` - A `identity` block as defined below.

* `location` - The Azure location where the Cognitive Services Account exists

* `local_auth_enabled` - Whether local authentication methods is enabled for the Cognitive Account.

* `kind` - The kind of the Cognitive Services Account

* `sku_name` - The SKU name of the Cognitive Services Account

* `endpoint` - The endpoint of the Cognitive Services Account

* `qna_runtime_endpoint` - If `kind` is `QnAMaker` the link to the QNA runtime.

* `primary_access_key` - The primary access key of the Cognitive Services Account

* `secondary_access_key` - The secondary access key of the Cognitive Services Account

-> **Note:** The `primary_access_key` and `secondary_access_key` properties are only available when `local_auth_enabled` is `true`.

* `tags` - A mapping of tags to assigned to the resource.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Cognitive Account.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Cognitive Account.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Cognitive Account.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Cognitive Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Services Account.
