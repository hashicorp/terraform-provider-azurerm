---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_rai_blocklist"
description: |-
  Manages a Cognitive Account Rai Blocklist.
---

# azurerm_cognitive_account_rai_blocklist

Manages a Cognitive Account Rai Blocklist.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "Brazil South"
}

resource "azurerm_cognitive_account" "example" {
  name                = "example-ca"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "OpenAI"
  sku_name            = "S0"
}

resource "azurerm_cognitive_account_rai_blocklist" "example" {
  name                 = "example-crb"
  cognitive_account_id = azurerm_cognitive_account.example.id
  description          = "Azure OpenAI Rai Blocklist"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Account Rai Blocklist. Changing this forces a new Cognitive Account Rai Blocklist to be created.

* `cognitive_account_id` - (Required) The ID of the Cognitive Services Account. Changing this forces a new Cognitive Account Rai Blocklist to be created.

---

* `description` - (Optional) A short description for the Cognitive Account Rai Blocklist.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cognitive Account Rai Blocklist.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Account Rai Blocklist.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Account Rai Blocklist.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Account Rai Blocklist.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Account Rai Blocklist.

## Import

Cognitive Account Rai Blocklist can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_rai_blocklist.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CognitiveServices/accounts/account1/raiBlocklists/raiblocklist1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices`: 2024-10-01
