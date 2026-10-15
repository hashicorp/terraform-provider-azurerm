---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_discovery_bookshelf"
description: |-
  Manages a Microsoft Discovery Bookshelf.
---

# azurerm_discovery_bookshelf

Manages a Microsoft Discovery Bookshelf.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "private_endpoint" {
  name                 = "privateEndpointSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_subnet" "search" {
  name                 = "searchSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/23"]

  delegation {
    name = "Microsoft.App/environments"

    service_delegation {
      name    = "Microsoft.App/environments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_discovery_bookshelf" "example" {
  name                       = "example-bookshelf"
  resource_group_name        = azurerm_resource_group.example.name
  location                   = azurerm_resource_group.example.location
  private_endpoint_subnet_id = azurerm_subnet.private_endpoint.id
  search_subnet_id           = azurerm_subnet.search.id
  workload_identity_ids      = [azurerm_user_assigned_identity.example.id]

  tags = {
    environment = "Production"
    indexSize   = "small"
  }
}
```

-> **Note:** Microsoft Discovery Bookshelves are only available in a subset of Azure Regions, and provisioning one requires Azure OpenAI model capacity (`text-embedding-3-small`, `gpt-5-mini` and `gpt-5.2` on the `GlobalStandard` SKU) in the same region. Requests which exceed the available quota are rejected during creation.

-> **Note:** The size of the search infrastructure provisioned for a Bookshelf is controlled by the `indexSize` tag, which accepts `small`, `medium` (the default) or `large`. Each size requires a corresponding amount of `Microsoft.App` memory optimized cores quota in the region - for example `medium` requires 80 cores. See [the Microsoft Discovery documentation](https://learn.microsoft.com/azure/microsoft-discovery/how-to-index-bookshelf-knowledgebase) for more information.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Microsoft Discovery Bookshelf. Must be between 3 and 24 characters long and may only contain letters, numbers and hyphens. Changing this forces a new Microsoft Discovery Bookshelf to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Microsoft Discovery Bookshelf should exist. Changing this forces a new Microsoft Discovery Bookshelf to be created.

* `location` - (Required) The Azure Region where the Microsoft Discovery Bookshelf should exist. Changing this forces a new Microsoft Discovery Bookshelf to be created.

* `private_endpoint_subnet_id` - (Required) The ID of the Subnet which should be used for Private Endpoint Connections. Changing this forces a new Microsoft Discovery Bookshelf to be created.

* `search_subnet_id` - (Required) The ID of the Subnet which should be used for search resources. Changing this forces a new Microsoft Discovery Bookshelf to be created.

-> **Note:** The Subnet specified in `search_subnet_id` must be delegated to `Microsoft.App/environments`.

* `workload_identity_ids` - (Required) A set of one or more User Assigned Identity IDs which should be used by knowledgebase workloads. Changing this forces a new Microsoft Discovery Bookshelf to be created.

---

* `customer_managed_key` - (Optional) A `customer_managed_key` block as defined below. When specified, data at rest is encrypted with a customer managed key. Adding or removing this block forces a new Microsoft Discovery Bookshelf to be created.

-> **Note:** `log_analytics_cluster_id` is required when `customer_managed_key` is specified.

* `log_analytics_cluster_id` - (Optional) The ID of the Log Analytics Cluster which should be used for debug logs. Changing this forces a new Microsoft Discovery Bookshelf to be created.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for this Microsoft Discovery Bookshelf. Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Microsoft Discovery Bookshelf.

---

A `customer_managed_key` block supports the following:

* `identity_client_id` - (Required) The Client ID of the Identity which should be used to access the Key Vault. This must be one of the Identities specified in `workload_identity_ids`. Changing this forces a new Microsoft Discovery Bookshelf to be created.

* `key_name` - (Required) The name of the Key Vault Key.

* `key_vault_uri` - (Required) The URI of the Key Vault. Changing this forces a new Microsoft Discovery Bookshelf to be created.

* `key_version` - (Optional) The version of the Key Vault Key. When omitted the latest version of the key is used. Removing a previously specified `key_version` forces a new Microsoft Discovery Bookshelf to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Microsoft Discovery Bookshelf.

* `managed_on_behalf_of_broker_resource_ids` - A list of Managed-On-Behalf-Of broker resource IDs.

* `managed_resource_group_name` - The name of the Resource Group containing the resources managed on behalf of this Microsoft Discovery Bookshelf.

* `uri` - The Bookshelf data plane API URI.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Microsoft Discovery Bookshelf.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft Discovery Bookshelf.
* `update` - (Defaults to 1 hour) Used when updating the Microsoft Discovery Bookshelf.
* `delete` - (Defaults to 1 hour) Used when deleting the Microsoft Discovery Bookshelf.

## Import

Microsoft Discovery Bookshelves can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_discovery_bookshelf.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resources/providers/Microsoft.Discovery/bookshelves/example-bookshelf
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Discovery` - 2026-06-01
