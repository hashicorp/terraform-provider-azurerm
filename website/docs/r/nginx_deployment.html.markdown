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

locals {
  config_content = base64encode(<<-EOT
http {
    server {
        listen 80;
        location / {
            auth_basic "Protected Area";
            auth_basic_user_file /opt/.htpasswd;
            default_type text/html;
        }
        include site/*.conf;
    }
}
EOT
  )

  protected_content = base64encode(<<-EOT
user:$apr1$VeUA5kt.$IjjRk//8miRxDsZvD4daF1
EOT
  )

  sub_config_content = base64encode(<<-EOT
location /bbb {
	default_type text/html;
	return 200 '<!doctype html><html lang="en"><head></head><body>
		<div>this one will be updated</div>
		<div>at 10:38 am</div>
	</body></html>';
}
EOT
  )
}


resource "azurerm_nginx_deployment" "example" {
  name                      = "example-nginx"
  resource_group_name       = azurerm_resource_group.example.name
  sku                       = "publicpreview_Monthly_gmz7xq9ge3py"
  location                  = azurerm_resource_group.example.location
  managed_resource_group    = "example"
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
  configuration {
    root_file = "/etc/nginx/nginx.conf"

    config_file {
      content      = local.config_content
      virtual_path = "/etc/nginx/nginx.conf"
    }

    config_file {
      content      = local.sub_config_content
      virtual_path = "/etc/nginx/site/b.conf"
    }

    protected_file {
      content      = local.protected_content
      virtual_path = "/opt/.htpasswd"
    }
  }

  lifecycle {
    ignore_changes = [configuration.0.protected_file]
  }

}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Nginx Deployment should exist. Changing this forces a new Nginx Deployment to be created.

* `name` - (Required) The name which should be used for this Nginx Deployment. Changing this forces a new Nginx Deployment to be created.

* `location` - (Required) The Azure Region where the Nginx Deployment should exist. Changing this forces a new Nginx Deployment to be created.

* `sku` - (Required) Specifies the Nginx Deployment SKU. Possible values include `standard_Monthly`. Changing this forces a new resource to be created.

* `managed_resource_group` - (Optional) Specify the managed resource group to deploy VNet injection related network resources. Changing this forces a new Nginx Deployment to be created.

---

* `capacity` - (Optional) Specify the number of NGINX capacity units for this NGINX deployment. Defaults to `20`.

-> **Note** For more information on NGINX capacity units, please refer to the [NGINX scaling guidance documentation](https://docs.nginx.com/nginxaas/azure/quickstart/scaling/)

* `auto_scale_profile` - (Optional) An `auto_scale_profile` block as defined below.

* `diagnose_support_enabled` - (Optional) Should the diagnosis support be enabled?

* `email` - (Optional) Specify the preferred support contact email address of the user used for sending alerts and notification.

* `identity` - (Optional) An `identity` block as defined below.

* `frontend_private` - (Optional) One or more `frontend_private` blocks as defined below. Changing this forces a new Nginx Deployment to be created.

* `frontend_public` - (Optional) A `frontend_public` block as defined below. Changing this forces a new Nginx Deployment to be created.

* `logging_storage_account` - (Optional) One or more `logging_storage_account` blocks as defined below.

* `network_interface` - (Optional) One or more `network_interface` blocks as defined below. Changing this forces a new Nginx Deployment to be created.

* `automatic_upgrade_channel` - (Optional) Specify the automatic upgrade channel for the NGINX deployment. Defaults to `stable`. The possible values are `stable` and `preview`.

* `configuration` - (Optional) Specify a custom `configuration` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Nginx Deployment.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Nginx Deployment. Possible values are `UserAssigned`, `SystemAssigned`.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned.

~> **NOTE:** This is required when `type` is set to `UserAssigned`.

---

A `frontend_private` block supports the following:

* `allocation_method` - (Required) Specify the method of allocating the private IP. Possible values are `Static` and `Dynamic`.

* `ip_address` - (Required) Specify the IP Address of this private IP.

* `subnet_id` - (Required) Specify the SubNet Resource ID to this Nginx Deployment.

---

A `frontend_public` block supports the following:

* `ip_address` - (Optional) Specifies a list of Public IP Resource ID to this Nginx Deployment.

---

A `logging_storage_account` block supports the following:

* `container_name` - (Optional) Specify the container name of Storage Account for logging.

* `name` - (Optional) The account name of the StorageAccount for Nginx Logging.

---

A `network_interface` block supports the following:

* `subnet_id` - (Required) Specify The SubNet Resource ID to this Nginx Deployment.

---

A `configuration` block supports the following:

* `root_file` - (Required) Specify the root file path of this Nginx Configuration.

---

-> **NOTE:** Either `package_data` or `config_file` must be specified - but not both.

* `package_data` - (Optional) Specify the package data for this configuration.

* `config_file` - (Optional) One or more `config_file` blocks as defined below.

* `protected_file` - (Optional) One or more `protected_file` blocks with sensitive information as defined below. If specified `config_file` must also be specified.

---

A `config_file` block supports the following:

* `content` - (Required) Specifies the base-64 encoded contents of this config file.

* `virtual_path` - (Required) Specify the path of this config file.

---

A `protected_file` (Protected File) block supports the following:

* `content` - (Required) Specifies the base-64 encoded contents of this config file (Sensitive).

* `virtual_path` - (Required) Specify the path of this config file.
---

An `auto_scale_profile` block supports the following:

* `name` - (Required) Specify the name of the autoscaling profile.

* `min_capacity` - (Required) Specify the minimum number of NGINX capacity units for this NGINX Deployment.

* `max_capacity` - (Required) Specify the maximum number of NGINX capacity units for this NGINX Deployment.

-> **NOTE:** If you're using autoscaling, you should use [Terraform's `ignore_changes` functionality](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changes) to ignore changes to the `capacity` field.

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
