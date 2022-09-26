---
subcategory: "Nginx"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nginx_configuration"
description: |-
  Manages a Nginx Configuration.
---

# azurerm_nginx_configuration

Manages a Nginx Configuration.

## Example Usage

```hcl
resource "azurerm_nginx_configuration" "test" {
  nginx_deployment_id = azurerm_nginx_deployment.test.id
  root_file           = "/etc/nginx/nginx.conf"

  config_file {
    content      = "aHR0cCB7DQogICAgc2VydmVyIHsNCiAgICAgICAgbGlzdGVuIDgwOw0KICAgICAgICBsb2NhdGlvbiAvIHsNCiAgICAgICAgICAgIGRlZmF1bHRfdHlwZSB0ZXh0L2h0bWw7DQogICAgICAgICAgICByZXR1cm4gMjAwICc8IWRvY3R5cGUgaHRtbD48aHRtbCBsYW5nPSJlbiI+PGhlYWQ+PC9oZWFkPjxib2R5Pg0KICAgICAgICAgICAgICAgIDxkaXY+dGhpcyBvbmUgd2lsbCBiZSB1cGRhdGVkPC9kaXY+DQogICAgICAgICAgICAgICAgPGRpdj5hdCAxMDozOCBhbTwvZGl2Pg0KICAgICAgICAgICAgPC9ib2R5PjwvaHRtbD4nOw0KICAgICAgICB9DQogICAgICAgIGluY2x1ZGUgc2l0ZS8qLmNvbmY7DQogICAgfQ0KfQ=="
    virtual_path = "/etc/nginx/nginx.conf"
  }

  config_file {
    content      = "DQogICAgICAgIGxvY2F0aW9uIC9iYmIgew0KICAgICAgICAgICAgZGVmYXVsdF90eXBlIHRleHQvaHRtbDsNCiAgICAgICAgICAgIHJldHVybiAyMDAgJzwhZG9jdHlwZSBodG1sPjxodG1sIGxhbmc9ImVuIj48aGVhZD48L2hlYWQ+PGJvZHk+DQogICAgICAgICAgICAgICAgPGRpdj50aGlzIG9uZSB3aWxsIGJlIHVwZGF0ZWQ8L2Rpdj4NCiAgICAgICAgICAgICAgICA8ZGl2PmF0IDEwOjM4IGFtPC9kaXY+DQogICAgICAgICAgICA8L2JvZHk+PC9odG1sPic7DQogICAgICAgIH0NCg=="
    virtual_path = "/etc/nginx/site/b.conf"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `config_file` - (Required) One or more `config_file` blocks as defined below.

* `nginx_deployment_id` - (Required) The ID of the Nginx Deployment. Changing this forces a new Nginx Configuration to be created.

* `root_file` - (Required) Specify the root file path of this Nginx Configuration.

---

* `package_data` - (Optional) Specify the package data for this configuration.

* `protected_file` - (Optional) One or more `config_file` blocks as defined below.

---

A `config_file` block supports the following:

* `content` - (Required) Specify the content of this config file. Content value should be encoded by base64

* `virtual_path` - (Required) Specify the path of this config file.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Nginx Configuration.

* `name` - The name of this Nginx Configuration. The value of configuration name is a fixed value as `default`. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Nginx Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Nginx Configuration.
* `update` - (Defaults to 10 minutes) Used when updating the Nginx Configuration.
* `delete` - (Defaults to 10 minutes) Used when deleting the Nginx Configuration.

## Import

Nginxs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nginx_configuration.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Nginx Configuration.NginxPlus/nginxDeployments/dep1/configurations/default
```