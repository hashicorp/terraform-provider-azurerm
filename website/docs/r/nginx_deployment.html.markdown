---
subcategory: "NGINX"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nginx_deployment"
description: |-
  Manages an NGINX Deployment.
---

# azurerm_nginx_deployment

Manages an NGINX Deployment.

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
  diagnose_support_enabled  = true
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
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the NGINX Deployment should exist. Changing this forces a new NGINX Deployment to be created.

* `name` - (Required) The name which should be used for this NGINX Deployment. Changing this forces a new NGINX Deployment to be created.

* `location` - (Required) The Azure Region where the NGINX Deployment should exist. Changing this forces a new NGINX Deployment to be created.

* `sku` - (Required) Specifies the NGINX Deployment SKU. Possible values are `standardv2_Monthly`, `basic_Monthly`.

-> **Note:** If you are setting the `sku` to `basic_Monthly`, you cannot specify a `capacity` or `auto_scale_profile`; basic plans do not support scaling. Other `sku`s require either `capacity` or `auto_scale_profile`. If you're using `basic_Monthly` with deployments created before v4.0, you may need to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changes) to ignore changes to the `capacity` field.

---

* `capacity` - (Optional) Specify the number of NGINX capacity units for this NGINX deployment.

-> **Note:** For more information on NGINX capacity units, please refer to the [NGINX scaling guidance documentation](https://docs.nginx.com/nginxaas/azure/quickstart/scaling/)

* `auto_scale_profile` - (Optional) An `auto_scale_profile` block as defined below.

* `diagnose_support_enabled` - (Optional) Should the metrics be exported to Azure Monitor?

* `email` - (Optional) Specify the preferred support contact email address for receiving alerts and notifications.

* `identity` - (Optional) An `identity` block as defined below.

* `frontend_private` - (Optional) One or more `frontend_private` blocks as defined below.

* `frontend_public` - (Optional) A `frontend_public` block as defined below.

* `network_interface` - (Optional) One or more `network_interface` blocks as defined below.

* `automatic_upgrade_channel` - (Optional) Specify the automatic upgrade channel for the NGINX deployment. Defaults to `stable`. The possible values are `stable` and `preview`.

* `web_application_firewall` - (Optional) A `web_application_firewall` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the NGINX Deployment.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the NGINX Deployment. Possible values are `SystemAssigned`, `UserAssigned` or `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned.

~> **Note:** This is required when `type` is set to `UserAssigned`.

---

A `frontend_private` block supports the following:

* `allocation_method` - (Required) Specify the method for allocating the private IP. Possible values are `Static` and `Dynamic`.

* `ip_address` - (Required) Specify the private IP Address.

* `subnet_id` - (Required) Specify the Subnet Resource ID for this NGINX Deployment.

---

A `frontend_public` block supports the following:

* `ip_address` - (Optional) Specifies a list of Public IP Resource ID to this NGINX Deployment.

---

A `network_interface` block supports the following:

* `subnet_id` - (Required) Specify The Subnet Resource ID for this NGINX Deployment.

---

An `auto_scale_profile` block supports the following:

* `name` - (Required) Specify the name of the autoscaling profile.

* `min_capacity` - (Required) Specify the minimum number of NGINX capacity units for this NGINX Deployment.

* `max_capacity` - (Required) Specify the maximum number of NGINX capacity units for this NGINX Deployment.

-> **Note:** If you're using autoscaling with deployments created before v4.0, you may need to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changes) to ignore changes to the `capacity` field.

---

A `web_application_firewall` - block supports the following:

* `activation_state_enabled` - (Required) Whether WAF is enabled/disabled for this NGINX Deployment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NGINX Deployment.

* `ip_address` - The IP address of the NGINX Deployment.

* `nginx_version` - The version of the NGINX Deployment.

* `dataplane_api_endpoint` - The dataplane API endpoint of the NGINX Deployment.

* `web_application_firewall.status` - A `web_application_firewall.status` block as defined below:

---

A `web_application_firewall.status` - block supports the following:

* `attack_signatures_package` - Indicates the version of the attack signatures package used by NGINX App Protect.

* `bot_signatures_package` - Indicates the version of the bot signatures package used by NGINX App Protect.

* `threat_campaigns_package` - Indicates the version of the threat campaigns package used by NGINX App Protect.

* `component_versions` - Indicates the version of the WAF Engine and Nginx WAF Module used by NGINX App Protect.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the NGINX Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the NGINX Deployment.
* `update` - (Defaults to 30 minutes) Used when updating the NGINX Deployment.
* `delete` - (Defaults to 30 minutes) Used when deleting the NGINX Deployment.

## Import

NGINX Deployments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nginx_deployment.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Nginx.NginxPlus/nginxDeployments/dep1
```
