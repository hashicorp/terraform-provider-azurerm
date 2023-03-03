---
subcategory: "Nginx"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nginx_configuration"
description: |-
  Manages the configuration for a Nginx Deployment.
---

# azurerm_nginx_configuration

Manages the configuration for a Nginx Deployment.

## Example Usage

```hcl
resource "azurerm_nginx_configuration" "test" {
  nginx_deployment_id = azurerm_nginx_deployment.test.id
  root_file           = "/etc/nginx/nginx.conf"

  config_file {
    content = base64encode(<<-EOT
http {
    server {
        listen 80;
        location / {
            default_type text/html;
            return 200 '<!doctype html><html lang="en"><head></head><body>
                <div>this one will be updated</div>
                <div>at 10:38 am</div>
            </body></html>';
        }
        include site/*.conf;
    }
}
EOT
    )
    virtual_path = "/etc/nginx/nginx.conf"
  }

  config_file {
    content = base64encode(<<-EOT
location /bbb {
 default_type text/html;
 return 200 '<!doctype html><html lang="en"><head></head><body>
  <div>this one will be updated</div>
  <div>at 10:38 am</div>
 </body></html>';
}
EOT
    )
    virtual_path = "/etc/nginx/site/b.conf"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `nginx_deployment_id` - (Required) The ID of the Nginx Deployment. Changing this forces a new Nginx Configuration to be created.

* `root_file` - (Required) Specify the root file path of this Nginx Configuration.

---

-> **NOTE:** Either `package_data` or `config_file` must be specified - but not both.

* `package_data` - (Optional) Specify the package data for this configuration.

* `config_file` - (Optional) One or more `config_file` blocks as defined below.

* `protected_file` - (Optional) One or more `protected_file` (Protected File) blocks with sensitive information as defined below. If specified `config_file` must also be specified.

---

A `config_file` block supports the following:

* `content` - (Required) Specifies the base-64 encoded contents of this config file.

* `virtual_path` - (Required) Specify the path of this config file.

---

A `protected_file` (Protected File) block supports the following:

* `content` - (Required) Specifies the base-64 encoded contents of this config file (Sensitive).

* `virtual_path` - (Required) Specify the path of this config file.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this Nginx Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Nginx Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Nginx Configuration.
* `update` - (Defaults to 10 minutes) Used when updating the Nginx Configuration.
* `delete` - (Defaults to 10 minutes) Used when deleting the Nginx Configuration.

## Import

Nginxs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nginx_configuration.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Nginx.NginxPlus/nginxDeployments/dep1/configurations/default
```
