---
subcategory: "Nginx"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nginx_certificate"
description: |-
  Manages a Certificate for an NGinx Deployment.
---

# azurerm_nginx_certificate

Manages a Certificate for an NGinx Deployment.

## Example Usage

```hcl
resource "azurerm_nginx_certificate" "test" {
  name                     = "examplecert"
  nginx_deployment_id      = azurerm_nginx_deployment.test.id
  key_virtual_path         = "/src/cert/soservermekey.key"
  certificate_virtual_path = "/src/cert/server.cert"
  key_vault_secret_id      = azurerm_key_vault_certificate.test.secret_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Nginx Certificate. Changing this forces a new Nginx Certificate to be created.

* `nginx_deployment_id` - (Required) The ID of the Nginx Deployment that this Certificate should be associated with. Changing this forces a new Nginx Certificate to be created.

* `certificate_virtual_path` - (Required) Specify the path to the cert file of this certificate. Changing this forces a new Nginx Certificate to be created.

* `key_virtual_path` - (Required) Specify the path to the key file of this certificate. Changing this forces a new Nginx Certificate to be created.

* `key_vault_secret_id` - (Required) Specify the ID of the Key Vault Secret for this certificate. Changing this forces a new Nginx Certificate to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this Nginx Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Nginx Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Nginx Certificate.
* `delete` - (Defaults to 10 minutes) Used when deleting the Nginx Certificate.

## Import

An Nginx Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nginx_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Nginx.NginxPlus/nginxDeployments/deploy1/certificates/cer1
```
