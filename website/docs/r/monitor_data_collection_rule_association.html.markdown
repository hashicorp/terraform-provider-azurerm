---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_data_collection_rule_association"
description: |-
  Manages a Data Collection Rule Association.
---

# azurerm_monitor_data_collection_rule_association

Manages a Data Collection Rule Association.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "virtualnetwork"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                = "machine"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  size                = "Standard_B1ls"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  admin_password = "example-Password@7890"

  disable_password_authentication = false

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_monitor_data_collection_rule" "example" {
  name                = "example-dcr"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  destinations {
    azure_monitor_metrics {
      name = "example-destination-metrics"
    }
  }
  data_flow {
    streams      = ["Microsoft-InsightsMetrics"]
    destinations = ["example-destination-metrics"]
  }
}

resource "azurerm_monitor_data_collection_endpoint" "example" {
  name                = "example-dce"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

# associate to a Data Collection Rule
resource "azurerm_monitor_data_collection_rule_association" "example1" {
  name                    = "example1-dcra"
  target_resource_id      = azurerm_linux_virtual_machine.example.id
  data_collection_rule_id = azurerm_monitor_data_collection_rule.example.id
  description             = "example"
}

# associate to a Data Collection Endpoint
resource "azurerm_monitor_data_collection_rule_association" "example2" {
  target_resource_id          = azurerm_linux_virtual_machine.example.id
  data_collection_endpoint_id = azurerm_monitor_data_collection_endpoint.example.id
  description                 = "example"
}

```

## Arguments Reference

The following arguments are supported:

* `target_resource_id` - (Required) The ID of the Azure Resource which to associate to a Data Collection Rule or a Data Collection Endpoint. Changing this forces a new resource to be created.

---

* `name` - (Optional) The name which should be used for this Data Collection Rule Association. Changing this forces a new Data Collection Rule Association to be created. Defaults to `configurationAccessEndpoint`.

-> **Note:** `name` is required when `data_collection_rule_id` is specified. And when `data_collection_endpoint_id` is specified, the `name` is populated with `configurationAccessEndpoint`.

* `data_collection_endpoint_id` - (Optional) The ID of the Data Collection Endpoint which will be associated to the target resource.

* `data_collection_rule_id` - (Optional) The ID of the Data Collection Rule which will be associated to the target resource.

-> **Note:** Exactly one of `data_collection_endpoint_id` and `data_collection_rule_id` blocks must be specified.

* `description` - (Optional) The description of the Data Collection Rule Association.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Collection Rule Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Collection Rule Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Collection Rule Association.
* `update` - (Defaults to 30 minutes) Used when updating the Data Collection Rule Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Collection Rule Association.

## Import

Data Collection Rules Association can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_data_collection_rule_association.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Insights/dataCollectionRuleAssociations/dca1
```
