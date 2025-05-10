---
subcategory: "NGINX"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nginx_api_key"
description: |-
  Manages a Dataplane API Key for an NGINX Deployment.
---

# azurerm_nginx_api_key

Manages the Dataplane API Key for an Nginx Deployment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "delegation"

    service_delegation {
      name = "NGINX.NGINXPLUS/nginxDeployments"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_nginx_deployment" "example" {
  name                      = "example-nginx"
  resource_group_name       = azurerm_resource_group.example.name
  sku                       = "standardv2_Monthly"
  location                  = azurerm_resource_group.example.location
  automatic_upgrade_channel = "stable"

  frontend_public {
    ip_address = [azurerm_public_ip.example.id]
  }
  network_interface {
    subnet_id = azurerm_subnet.example.id
  }

  capacity = 20

  email = "user@test.com"
}

resource "azurerm_nginx_api_key" "example" {
  name                = "example-api-key"
  nginx_deployment_id = azurerm_nginx_deployment.example.id
  # We don't recommend hard coding the secret_text value. Please refer to the secret_text reference below for sources on how to manage sensitive input variables.
  secret_text   = "727c8642-6807-4254-9d02-ae93bfad21de"
  end_date_time = "2027-01-01T00:00:00Z"
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the NGINX Dataplane API Key. Changing this forces a new resource to be created.

- `nginx_deployment_id` - (Required) The ID of the NGINX Deployment that the API key is associated with. Changing this forces a new resource to be created.

- `end_date_time` - (Required) The RFC3339 formatted date-time after which this Dataplane API Key is no longer valid. The maximum value is now+2y.

- `secret_text` - (Required) The value used as the Dataplane API Key. The API key requirements can be found in the [NGINXaaS Documentation](https://docs.nginx.com/nginxaas/azure/quickstart/loadbalancer-kubernetes/#create-an-nginxaas-data-plane-api-key).

-> **Note:** The `secret_text` contains a Dataplane API Key that can be used to modify NGINX upstream servers. The following sources are useful in learning to manage sensitive data.

  - [Sensitive Data in State](https://developer.hashicorp.com/terraform/language/state/sensitive-data)
  - [Protect sensitive input variables](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `hint` - The first three characters of the secret text to help identify it in use.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the NGINX Dataplane API Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the NGINX Dataplane API Key.
* `update` - (Defaults to 5 minutes) Used when updating the NGINX Dataplane API Key.
* `delete` - (Defaults to 5 minutes) Used when deleting the NGINX Dataplane API Key.

## Import

An NGINX Dataplane API Key can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nginx_api_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Nginx.NginxPlus/nginxDeployments/deploy1/apiKeys/key1
```
