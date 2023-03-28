---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_custom_domain"
description: |-
  Manages an Azure Web Pubsub Custom Domain.
---

# azurerm_web_pubsub_custom_domain

Manages an Azure Web Pubsub Custom Domain.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_web_pubsub" "example" {
  name                = "example-webpubsub"
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
    object_id = azurerm_web_pubsub_service.test.identity[0].principal_id
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
resource "azurerm_web_pubsub_custom_certificate_binding" "test" {
  name                  = "example-certbinding"
  web_pubsub_service_id = azurerm_web_pubsub_service.example.id
  custom_certificate_id = azurerm_key_vault_certificate.example.id
  certificate_version   = "ec6faxx"
  depends_on            = [azurerm_key_vault_access_policy.example]
}


resource "azurerm_web_pubsub_custom_domain" "test" {
  name                             = "example-domain"
  web_pubsub_service_id            = azurerm_web_pubsub.test.id
  domain_name                      = "tftest.com"
  web_pubsub_custom_certificate_id = azurerm_web_pubsub_custom_certificate_binding.test.id
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Web Pubsub Custom Domain. Changing this forces a new resource to be created.

* `web_pubsub_id` - (Required) The Web Pubsub ID of the Web Pubsub Custom Domain. Changing this forces a new resource to be created.

* `web_pubsub_custom_certificate_id` - (Required) The Web Pubsub custom certificate id of the Web Pubsub Custom Domain service. Changing this forces a new resource to be created.

* `domain_name` - (Required) The custom domain name of the Web Pubsub Custom Domain service.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Web Pubsub Custom Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Custom Domain of the Web Pubsub service
* `read` - (Defaults to 5 minutes) Used when retrieving the Custom Domain of the Web Pubsub service
* `update` - (Defaults to 30 minutes) Used when updating the Custom Domain of the Web Pubsub service
* `delete` - (Defaults to 30 minutes) Used when deleting the Custom Domain of the Web Pubsub service

## Import

Custom Domain for a Web Pubsub service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_web_pubsub_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/WebPubSub/webpubsub1/customDomains/customDomain1
```
