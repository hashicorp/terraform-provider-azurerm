---
subcategory: "Nginx"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nginx_deployment"
description: |-
  Manages a Nginx Deployment.
---

# azurerm_nginx_deployment

Manages a Nginx Deployment.

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
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Nginx Deployment should exist. Changing this forces a new Nginx Deployment to be created.

* `name` - (Required) The name which should be used for this Nginx Deployment. Changing this forces a new Nginx Deployment to be created.

* `location` - (Required) The Azure Region where the Nginx Deployment should exist. Changing this forces a new Nginx Deployment to be created.

* `sku` - (Required) Specify the Name of Nginx deployment SKU. The possible value are `publicpreview_Monthly_gmz7xq9ge3py` and `standard_Monthly`.

* `managed_resource_group` - (Optional) Specify the managed resource group to deploy VNet injection related network resources. Changing this forces a new Nginx Deployment to be created.

---

* `diagnose_support_enabled` - (Optional) Should the diagnosis support be enabled?

* `identity` - (Optional) An `identity` block as defined below.

* `frontend_private` - (Optional) One or more `frontend_private` blocks as defined below. Changing this forces a new Nginx Deployment to be created.

* `frontend_public` - (Optional) A `frontend_public` block as defined below. Changing this forces a new Nginx Deployment to be created.

* `logging_storage_account` - (Optional) One or more `logging_storage_account` blocks as defined below.

* `network_interface` - (Optional) One or more `network_interface` blocks as defined below. Changing this forces a new Nginx Deployment to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Nginx Deployment.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Nginx Deployment. Possible values is `UserAssigned` where you can specify the Service Principal IDs in the `identity_ids` field.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`.

---

A `frontend_private` block supports the following:

* `allocation_method` - (Required) Specify the methos of allocating the private IP. Possible values are `Static` and `Dynamic`.

* `ip_address` - (Required) Specify the IP Address of this private IP.

* `subnet_id` - (Required) Specify the SubNet Resource ID to this Nginx Deployment.

---

A `frontend_public` block supports the following:

* `ip_address` - (Optional) Specifies a list of Public IP Resouce ID to this Nginx Deployment.

---

A `logging_storage_account` block supports the following:

* `container_name` - (Optional) Specify the container name of Stoage Account for logging.

* `name` - (Optional) The account name of the StorageAccount for Nginx Logging.

---

A `network_interface` block supports the following:

* `subnet_id` - (Required) Specify The SubNet Resource ID to this Nginx Deployment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Nginx Deployment.

* `ip_address` - The IP address of the deployment.

* `nginx_version` - The version of deployed nginx.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Nginx Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Nginx Deployment.
* `update` - (Defaults to 30 minutes) Used when updating the Nginx Deployment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Nginx Deployment.

## Import

Nginx Deployments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nginx_deployment.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Nginx.NginxPlus/nginxDeployments/dep1
```
