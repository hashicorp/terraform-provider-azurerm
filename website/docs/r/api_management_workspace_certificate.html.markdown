---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_certificate"
description: |-
  Manages a Certificate within an API Management Workspace.
---

# azurerm_api_management_workspace_certificate

Manages a Certificate within an API Management Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
}

resource "azurerm_api_management_workspace_certificate" "example" {
  name                        = "example-cert"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  certificate_data_base64     = filebase64("example.pfx")
  password                    = "terraform"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API Management Workspace Certificate. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

* `certificate_data_base64` - (Optional) Specifies the base64-encoded string containing the certificate in PKCS#12 (.pfx) format.

-> **Note:** This is required when `password` is specified. Exactly one of `certificate_data_base64` or `key_vault_secret_id` must be specified.

* `key_vault_secret_id` - (Optional) Specifies the ID of the key vault secret.

-> **Note:** This is required when `user_assigned_identity_client_id` is specified. Exactly one of `certificate_data_base64` or `key_vault_secret_id` must be specified.

* `password` - (Optional) Specifies the password used to access the `certificate_data_base64`.

* `user_assigned_identity_client_id` - (Optional) Specifies the client ID of user-assigned identity to be used for accessing the `key_vault_secret_id`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace Certificate.

* `expiration` - The expiration date of the API Management Workspace Certificate.

* `subject` - The subject name of the API Management Workspace Certificate.

* `thumbprint` - The thumbprint of the API Management Workspace Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Certificate.

## Import

API Management Workspace Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1/certificates/certificate1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
