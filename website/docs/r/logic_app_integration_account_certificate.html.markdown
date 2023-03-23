---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account_certificate"
description: |-
  Manages a Logic App Integration Account Certificate.
---

# azurerm_logic_app_integration_account_certificate

Manages a Logic App Integration Account Certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_logic_app_integration_account" "example" {
  name                = "example-ia"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard"
}

resource "azurerm_logic_app_integration_account_certificate" "example" {
  name                     = "example-iac"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.example.name
  public_certificate       = "MIIDbzCCAlegAwIBAgIJAIzjRD36sIbbMA0GCSqGSIb3DQEBCwUAME0xCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApTb21lLVN0YXRlMRIwEAYDVQQKDAl0ZXJyYWZvcm0xFTATBgNVBAMMDHRlcnJhZm9ybS5pbzAgFw0xNzA0MjEyMDA1MjdaGA8yMTE3MDMyODIwMDUyN1owTTELMAkGA1UEBhMCVVMxEzARBgNVBAgMClNvbWUtU3RhdGUxEjAQBgNVBAoMCXRlcnJhZm9ybTEVMBMGA1UEAwwMdGVycmFmb3JtLmlvMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3L9L5szT4+FLykTFNyyPjy/k3BQTYAfRQzP2dhnsuUKm3cdPC0NyZ+wEXIUGhoDO2YG6EYChOl8fsDqDOjloSUGKqYw++nlpHIuUgJx8IxxG2XkALCjFU7EmF+w7kn76d0ezpEIYxnLP+KG2DVornoEt1aLhv1MLmpgEZZPhDbMSLhSYWeTVRMayXLwqtfgnDumQSB+8d/1JuJqrSI4pD12JozVThzb6hsjfb6RMX4epPmrGn0PbTPEEA6awmsxBCXB0s13nNQt/O0hLM2agwvAyozilQV+s616Ckgk6DJoUkqZhDy7vPYMIRSr98fBws6zkrV6tTLjmD8xAvobePQIDAQABo1AwTjAdBgNVHQ4EFgQUXIqO421zMMmbcRRX9wctZFCQuPIwHwYDVR0jBBgwFoAUXIqO421zMMmbcRRX9wctZFCQuPIwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAr82NeT3BYJOKLlUL6Om5LjUF66ewcJjG9ltdvyQwVneMcq7t5UAPxgChzqNRVk4da8PzkXpjBJyWezHupdJNX3XqeUk2kSxqQ6/gmhqvfI3y7djrwoO6jvMEY26WqtkTNORWDP3THJJVimC3zV+KMU5UBVrEzhOVhHSU709lBP75o0BBn3xGsPqSq9k8IotIFfyAc6a+XP3+ZMpvh7wqAUml7vWa5wlcXExCx39h1balfDSLGNC4swWPCp9AMnQR0p+vMay9hNP1Eh+9QYUai14d5KS3cFV+KxE1cJR5HD/iLltnnOEbpMsB0eVOZWkFvE7Y5lW0oVSAfin5TwTJMQ=="
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account Certificate. Changing this forces a new Logic App Integration Account Certificate to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account Certificate should exist. Changing this forces a new Logic App Integration Account Certificate to be created.

* `integration_account_name` - (Required) The name of the Logic App Integration Account. Changing this forces a new Logic App Integration Account Certificate to be created.

* `key_vault_key` - (Optional) A `key_vault_key` block as documented below.

* `metadata` - (Optional) A JSON mapping of any Metadata for this Logic App Integration Account Certificate.

* `public_certificate` - (Optional) The public certificate for the Logic App Integration Account Certificate.

---

A `key_vault_key` block exports the following:

* `key_name` - (Required) The name of Key Vault Key.

* `key_vault_id` - (Required) The ID of the Key Vault.

* `key_version` - (Optional) The version of Key Vault Key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account Certificate.

## Import

Logic App Integration Account Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1/certificates/certificate1
```
