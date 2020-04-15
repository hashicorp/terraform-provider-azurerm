---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_share_invitation"
description: |-
  Gets information about an existing Data Share Invitation.
---

# Data Source: azurerm_data_share_invitation

Use this data source to access information about an existing Data Share Invitation.

## Example Usage

```hcl
data "azurerm_data_share_invitation" "example" {
  name = "existing"
  share_id = "TODO"
}

output "id" {
  value = data.azurerm_data_share_invitation.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Data Share Invitation. Changing this forces a new Data Share Invitation to be created.

* `share_id` - (Required) The ID of the TODO. Changing this forces a new Data Share Invitation to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Share Invitation.

* `invitation_id` - The ID of the TODO.

* `invitation_status` - TODO.

* `target_email` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Invitation.