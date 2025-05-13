---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_custom_domain"
description: |-
  Manages a Container App Custom Domain.
---

# azurerm_container_app_custom_domain

Manages a Container App Custom Domain.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "contoso.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_dns_txt_record" "example" {
  name                = "asuid.example"
  resource_group_name = azurerm_dns_zone.example.resource_group_name
  zone_name           = azurerm_dns_zone.example.name
  ttl                 = 300

  record {
    value = azurerm_container_app.example.custom_domain_verification_id
  }
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "example" {
  name                       = "Example-Environment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_container_app_environment_certificate" "example" {
  name                         = "myfriendlyname"
  container_app_environment_id = azurerm_container_app_environment.example.id
  certificate_blob             = filebase64("path/to/certificate_file.pfx")
  certificate_password         = "$3cretSqu1rreL"
}

resource "azurerm_container_app" "example" {
  name                         = "example-app"
  container_app_environment_id = azurerm_container_app_environment.example.id
  resource_group_name          = azurerm_resource_group.example.name
  revision_mode                = "Single"

  template {
    container {
      name   = "examplecontainerapp"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
  ingress {
    allow_insecure_connections = false
    external_enabled           = true
    target_port                = 5000
    transport                  = "http"
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }
}

resource "azurerm_container_app_custom_domain" "example" {
  name                                     = trimsuffix(trimprefix(azurerm_dns_txt_record.api.fqdn, "asuid."), ".")
  container_app_id                         = azurerm_container_app.example.id
  container_app_environment_certificate_id = azurerm_container_app_environment_certificate.example.id
  certificate_binding_type                 = "SniEnabled"
}

```

## Example Usage - Managed Certificate

```hcl
resource "azurerm_container_app_custom_domain" "example" {
  name             = trimsuffix(trimprefix(azurerm_dns_txt_record.api.fqdn, "asuid."), ".")
  container_app_id = azurerm_container_app.example.id

  lifecycle {
    // When using an Azure created Managed Certificate these values must be added to ignore_changes to prevent resource recreation.
    ignore_changes = [certificate_binding_type, container_app_environment_certificate_id]
  }
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The fully qualified name of the Custom Domain. Must be the CN or a named SAN in the certificate specified by the `container_app_environment_certificate_id`. Changing this forces a new resource to be created.

~> **Note:** The Custom Domain verification TXT record requires a prefix of `asuid.`, however, this must be trimmed from the `name` property here. See the [official docs](https://learn.microsoft.com/en-us/azure/container-apps/custom-domains-certificates) for more information.

* `container_app_id` - (Required) The ID of the Container App to which this Custom Domain should be bound. Changing this forces a new resource to be created.

* `container_app_environment_certificate_id` - (Optional) The ID of the Container App Environment Certificate to use. Changing this forces a new resource to be created.

-> **Note:** Omit this value if you wish to use an Azure Managed certificate. You must create the relevant DNS verification steps before this process will be successful.

* `certificate_binding_type` - (Optional) The Certificate Binding type. Possible values include `Disabled` and `SniEnabled`.  Required with `container_app_environment_certificate_id`. Changing this forces a new resource to be created.

!> **Note:** If using an Azure Managed Certificate `container_app_environment_certificate_id` and `certificate_binding_type` should be added to `ignore_changes` to prevent resource recreation due to these values being modified asynchronously outside of Terraform.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `container_app_environment_managed_certificate_id` - The ID of the Container App Environment Managed Certificate to use.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App.

## Import

A Container App Custom Domain can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_custom_domain.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/containerApps/myContainerApp/customDomainName/mycustomdomain.example.com"
```
