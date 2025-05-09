---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_managed_certificate"
description: |-
  Manages a Container App Environment Managed Certificate.
---

# azurerm_container_app_environment_managed_certificate

Manages a Container App Environment Managed Certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_app_environment" "example" {
  name                = "myEnvironment"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_dns_cname_record" "example" {
  resource_group_name = var.dns_zone_resource_group_name
  name                = "friendlydomain"
  zone_name           = var.dns_zone_name
  ttl                 = 300

  record = azurerm_container_app.example.ingress.0.fqdn
}

resource "azurerm_dns_txt_record" "example" {
  resource_group_name = var.dns_zone_resource_group_name
  name                = "asuid.friendlydomain"
  zone_name           = var.dns_zone_name
  ttl                 = 300

  record {
    value = azurerm_container_app.example.custom_domain_verification_id
  }
}

resource "azurerm_container_app_custom_domain" "example" {
  name             = trimsuffix(trimprefix(azurerm_dns_txt_record.example.fqdn, "asuid."), ".")
  container_app_id = azurerm_container_app.example.id

  lifecycle {
    // When using an Azure created Managed Certificate these values must be added to ignore_changes to prevent resource recreation.
    ignore_changes = [certificate_binding_type, container_app_environment_certificate_id]
  }
}

resource "azurerm_container_app_environment_managed_certificate" "example" {
  name                           = "mymanagedcertificate"
  container_app_environment_id   = azurerm_container_app_environment.example.id
  subject_name                   = "${azurerm_dns_cname_record.example.name}.${azurerm_dns_cname_record.example.zone_name}"
  domain_control_validation_type = "CNAME"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Apps Environment Managed Certificate. Changing this forces a new resource to be created.

* `container_app_environment_id` - (Required) The Container App Managed Environment ID to configure this Certificate on. Changing this forces a new resource to be created.

* `domain_control_validation_type` - (Required) Type of domain control validation. Possible values include `CNAME`, `HTTP` and `TXT`. Changing this forces a new resource to be created.

* `subject_name` - (Required) The subject name of the certificate. Changing this forces a new resource to be created.

---

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Managed Certificate

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment Managed Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment Managed Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Managed Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment Managed Certificate.

## Import

A Container App Environment Managed Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment_managed_certificate.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myenv/managedCertificates/mycertificate"
```
