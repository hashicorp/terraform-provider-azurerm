---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_endpoint_application_security_group_association"
description: |-
  Manages an association between Private Endpoint and Application Security Group.

---

# azurerm_private_endpoint_application_security_group_association

Manages an association between Private Endpoint and Application Security Group.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-PEASGAsso"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "examplevnet"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "service" {
  name                 = "examplenetservice"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.5.1.0/24"]

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "examplenetendpoint"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_public_ip" "example" {
  name                = "examplepip"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "examplelb"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  frontend_ip_configuration {
    name                 = azurerm_public_ip.example.name
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_private_link_service" "example" {
  name                           = "examplePLS"
  location                       = azurerm_resource_group.example.location
  resource_group_name            = azurerm_resource_group.example.name
  auto_approval_subscription_ids = [data.azurerm_subscription.current.subscription_id]
  visibility_subscription_ids    = [data.azurerm_subscription.current.subscription_id]

  nat_ip_configuration {
    name      = "primaryIpConfiguration"
    primary   = true
    subnet_id = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.example.frontend_ip_configuration[0].id
  ]
}

resource "azurerm_private_endpoint" "example" {
  name                = "example-privatelink"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = azurerm_private_link_service.example.name
    is_manual_connection           = false
    private_connection_resource_id = azurerm_private_link_service.example.id
  }
}

resource "azurerm_application_security_group" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_endpoint_application_security_group_association" "example" {
  private_endpoint_id           = azurerm_private_endpoint.example.id
  application_security_group_id = azurerm_application_security_group.example.id
}
```

## Argument Reference

The following arguments are supported:

* `application_security_group_id` - (Required) The id of application security group to associate. Changing this forces a new resource to be created.

* `private_endpoint_id` - (Required) The id of private endpoint to associate. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The (Terraform specific) ID of the association between Private Endpoint and Application Security Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the association between Private Endpoint and Application Security Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the association between Private Endpoint and Application Security Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the association between Private Endpoint and Application Security Group.

## Import

Associations between Private Endpoint and Application Security Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_endpoint_application_security_group_association.association1 "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/privateEndpoints/endpoints1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/applicationSecurityGroups/securityGroup1",
```

-> **Note:** This ID is specific to Terraform - and is of the format `{privateEndpointId}|{applicationSecurityGroupId}`.
