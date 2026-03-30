---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_next_generation_firewall_metrics"
description: |-
  Manages a Palo Alto Next Generation Firewall Metrics configuration.
---

# azurerm_palo_alto_next_generation_firewall_metrics

Manages a Palo Alto Next Generation Firewall Metrics configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_public_ip" "example" {
  name                = "example-publicip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_network_security_group" "example" {
  name                = "example-nsg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "trust" {
  name                 = "example-trust-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "trusted"

    service_delegation {
      name = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "trust" {
  subnet_id                 = azurerm_subnet.trust.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_subnet" "untrust" {
  name                 = "example-untrust-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "untrusted"

    service_delegation {
      name = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "untrust" {
  subnet_id                 = azurerm_subnet.untrust.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_palo_alto_local_rulestack" "example" {
  name                = "example-rulestack"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  depends_on = [
    azurerm_subnet_network_security_group_association.trust,
    azurerm_subnet_network_security_group_association.untrust,
  ]
}

resource "azurerm_palo_alto_local_rulestack_rule" "example" {
  name         = "example-rule"
  rulestack_id = azurerm_palo_alto_local_rulestack.example.id
  priority     = 1001
  action       = "Allow"
  protocol     = "application-default"
  applications = ["any"]

  destination {
    cidrs = ["any"]
  }

  source {
    cidrs = ["any"]
  }
}

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack" "example" {
  name                = "example-ngfw"
  resource_group_name = azurerm_resource_group.example.name
  rulestack_id        = azurerm_palo_alto_local_rulestack.example.id

  network_profile {
    public_ip_address_ids = [azurerm_public_ip.example.id]

    vnet_configuration {
      virtual_network_id  = azurerm_virtual_network.example.id
      trusted_subnet_id   = azurerm_subnet.trust.id
      untrusted_subnet_id = azurerm_subnet.untrust.id
    }
  }

  depends_on = [azurerm_palo_alto_local_rulestack_rule.example]
}

resource "azurerm_palo_alto_next_generation_firewall_metrics" "example" {
  firewall_id                            = azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack.example.id
  application_insights_connection_string = azurerm_application_insights.example.connection_string
  application_insights_resource_id       = azurerm_application_insights.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `firewall_id` - (Required) The ID of the Palo Alto Next Generation Firewall. Changing this forces a new resource to be created.

* `application_insights_connection_string` - (Required) The connection string of the Application Insights resource used for metrics collection.

* `application_insights_resource_id` - (Required) The resource ID of the Application Insights resource used for metrics collection.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Palo Alto Next Generation Firewall Metrics configuration. This is the same as the `firewall_id`.

* `pan_etag` - A read-only string representing the last create or update operation on this resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Next Generation Firewall Metrics.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Next Generation Firewall Metrics.
* `update` - (Defaults to 30 minutes) Used when updating the Palo Alto Next Generation Firewall Metrics.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Next Generation Firewall Metrics.

## Import

Palo Alto Next Generation Firewall Metrics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_next_generation_firewall_metrics.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/firewalls/myFirewall
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `PaloAltoNetworks.Cloudngfw` - 2025-10-08
