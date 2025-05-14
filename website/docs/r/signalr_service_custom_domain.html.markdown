---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_signalr_service_custom_domain"
description: |-
  Manages an Azure SignalR Custom Domain.
---

# azurerm_signalr_service_custom_domain

Manages an Azure SignalR Custom Domain.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_signalr_service" "example" {
  name                = "example-signalr"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Premium_P1"
    capacity = 1
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault" "example" {
  name                = "example-keyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Get",
      "List",
    ]

    secret_permissions = [
      "Get",
      "List",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_signalr_service.test.identity[0].principal_id

    certificate_permissions = [
      "Create",
      "Get",
      "List",
    ]

    secret_permissions = [
      "Get",
      "List",
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

resource "azurerm_signalr_service_custom_certificate" "test" {
  name                  = "example-cert"
  signalr_service_id    = azurerm_signalr_service.example.id
  custom_certificate_id = azurerm_key_vault_certificate.example.id

  depends_on = [azurerm_key_vault_access_policy.example]
}

resource "azurerm_signalr_service_custom_domain" "test" {
  name                          = "example-domain"
  signalr_service_id            = azurerm_signalr_service.test.id
  domain_name                   = "tftest.com"
  signalr_custom_certificate_id = azurerm_signalr_service_custom_certificate.test.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the SignalR Custom Domain. Changing this forces a new resource to be created.

* `domain_name` - (Required) Specifies the custom domain name of the SignalR Custom Domain. Changing this forces a new resource to be created.

-> **Note:** Please ensure the custom domain name is included in the Subject Alternative Names of the selected SignalR Custom Certificate.

* `signalr_service_id` - (Required) Specifies the SignalR ID of the SignalR Custom Domain. Changing this forces a new resource to be created.

* `signalr_custom_certificate_id` - (Required) Specifies the SignalR Custom Certificate ID of the SignalR Custom Domain. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SignalR Custom Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the custom domain of the SignalR service
* `read` - (Defaults to 5 minutes) Used when retrieving the custom domain of the SignalR service
* `delete` - (Defaults to 30 minutes) Used when deleting the custom domain of the SignalR service

## Import

Custom Domain for a SignalR service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_signalr_service_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/signalR/signalr1/customDomains/customDomain1
```
