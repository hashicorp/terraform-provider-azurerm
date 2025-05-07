---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_synapse_workspace"
description: |-
  Gets information about an existing Synapse Workspace.
---

# Data Source: azurerm_synapse_workspace

Use this data source to access information about an existing Synapse Workspace.

## Example Usage

```hcl
data "azurerm_synapse_workspace" "example" {
  name                = "existing"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_synapse_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Synapse Workspace.

* `resource_group_name` - (Required) The name of the Resource Group where the Synapse Workspace exists.

## Attributes Reference

the following Attributes are exported:

* `id` - The ID of the synapse Workspace.

* `location` - The Azure location where the Synapse Workspace exists.

* `connectivity_endpoints` - A map of Connectivity endpoints for this Synapse Workspace.

* `tags` - A mapping of tags assigned to the resource.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Synapse Workspace.

---

The `identity` block exports the following:

* `type` - The Identity Type for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Workspace.
