---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account"
sidebar_current: "docs-azurerm-datasource-cognitive-account"
description: |-
  Gets information about an existing Cognitive Services Account.

---

# Data Source: azurerm_cognitive_account

Use this data source to access information about an existing Cognitive Services Account.

~> **Note:** All arguments including the access key values will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_cognitive_account" "test" {
  name      = "my-cognitive-account"
  resource_group_name = "cognitive_account_rg"
}

output "primary_access_key" {
  value = "${data.azurerm_cognitive_account.test.primary_access_key}"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cognitive Services Account.

* `resource_group_name` - (Required) Specifies the name of the resource group where the Cognitive Services Account resides.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure location where this Cognitive Services Account exists

* `kind` - The kind of this Cognitive Services Account

* `sku` - The sku of this Cognitive Services Account as defined below.

* `endpoint` - The endpoint of this Cognitive Services Account

* `primary_access_key` - The primary access key of this Cognitive Services Account

* `secondary_access_key` - The secondary access key of this Cognitive Services Account

* `tags` - A mapping of tags to assigned to the resource.

---

* `sku` supports the following:

* `name` - The Sku Name used for this Cognitive Services Account.

* `tier` - The Sku Tier used for this Cognitive Services Account.
