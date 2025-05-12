---
subcategory: "NGINX"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nginx_certificate"
description: |-
  Manages a Certificate for an NGINX Deployment.
---

# azurerm_nginx_certificate

Manages a Certificate for an NGINX Deployment.

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
  name                     = "example-nginx"
  resource_group_name      = azurerm_resource_group.example.name
  sku                      = "publicpreview_Monthly_gmz7xq9ge3py"
  location                 = azurerm_resource_group.example.location
  managed_resource_group   = "example"
  diagnose_support_enabled = true

  frontend_public {
    ip_address = [azurerm_public_ip.example.id]
  }
  network_interface {
    subnet_id = azurerm_subnet.example.id
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                = "examplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "GetIssuers",
      "Import",
      "List",
      "ListIssuers",
      "ManageContacts",
      "ManageIssuers",
      "SetIssuers",
      "Update",
    ]
  }
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "imported-cert"
  key_vault_id = azurerm_key_vault.example.id

  certificate {
    contents = filebase64("certificate-to-import.pfx")
    password = ""
  }
}

resource "azurerm_nginx_certificate" "example" {
  name                     = "examplecert"
  nginx_deployment_id      = azurerm_nginx_deployment.example.id
  key_virtual_path         = "/src/cert/soservermekey.key"
  certificate_virtual_path = "/src/cert/server.cert"
  key_vault_secret_id      = azurerm_key_vault_certificate.example.secret_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this NGINX Certificate. Changing this forces a new NGINX Certificate to be created.

* `nginx_deployment_id` - (Required) The ID of the NGINX Deployment that this Certificate should be associated with. Changing this forces a new NGINX Certificate to be created.

* `certificate_virtual_path` - (Required) Specify the path to the certificate file of this certificate.

* `key_virtual_path` - (Required) Specify the path to the key file of this certificate.

* `key_vault_secret_id` - (Required) Specify the ID of the Key Vault Secret for this certificate.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this NGINX Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the NGINX Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the NGINX Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the NGINX Certificate.
* `delete` - (Defaults to 10 minutes) Used when deleting the NGINX Certificate.

## Import

An NGINX Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nginx_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Nginx.NginxPlus/nginxDeployments/deploy1/certificates/cer1
```
