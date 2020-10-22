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
  value = "${data.azurerm_cognitive_account.test.primary_access_key}"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cognitive Services Account.

* `resource_group_name` - (Required) Specifies the name of the resource group where the Cognitive Services Account resides.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure location where the Cognitive Services Account exists

* `kind` - The kind of the Cognitive Services Account

* `sku_name` - The sku name of the Cognitive Services Account

* `endpoint` - The endpoint of the Cognitive Services Account

* `qna_runtime_endpoint` - If `kind` is `QnAMaker` the link to the QNA runtime.
* `primary_access_key` - The primary access key of the Cognitive Services Account

* `secondary_access_key` - The secondary access key of the Cognitive Services Account

* `tags` - A mapping of tags to assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Services Account.
