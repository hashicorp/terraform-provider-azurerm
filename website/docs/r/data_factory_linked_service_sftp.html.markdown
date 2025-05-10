---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_sftp"
description: |-
  Manages a Linked Service (connection) between an SFTP Server and Azure Data Factory.
---

# azurerm_data_factory_linked_service_sftp

Manages a Linked Service (connection) between a SFTP Server and Azure Data Factory.

~> **Note:** All arguments including the client secret will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_linked_service_sftp" "example" {
  name                = "example"
  data_factory_id     = azurerm_data_factory.example.id
  authentication_type = "Basic"
  host                = "http://www.bing.com"
  port                = 22
  username            = "foo"
  password            = "bar"
}
```

## Argument Reference

The following supported arguments are common across all Azure Data Factory Linked Services:

* `name` - (Required) Specifies the name of the Data Factory Linked Service. Changing this forces a new resource to be created. Must be unique within a data factory. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `description` - (Optional) The description for the Data Factory Linked Service.

* `integration_runtime_name` - (Optional) The name of the integration runtime to associate with the Data Factory Linked Service.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service.

The following supported arguments are specific to SFTP Linked Service:

* `authentication_type` - (Required) The type of authentication used to connect to the SFTP server. Valid options are `MultiFactor`, `Basic` and `SshPublicKey`.

* `host` - (Required) The SFTP server hostname.

* `port` - (Required) The TCP port number that the SFTP server uses to listen for client connection. Default value is 22.

* `username` - (Required) The username used to log on to the SFTP server.

* `password` - (Optional) Password to log on to the SFTP Server for Basic Authentication.

* `key_vault_password` - (Optional) A `key_vault_password` block as defined below.

~> **Note:** Either `password` or `key_vault_password` is required when `authentication_type` is set to `Basic`.

* `private_key_content_base64` - (Optional) The Base64 encoded private key content in OpenSSH format used to log on to the SFTP server.

* `key_vault_private_key_content_base64` - (Optional) A `key_vault_private_key_content_base64` block as defined below.

* `private_key_path` - (Optional) The absolute path to the private key file that the self-hosted integration runtime can access.

~> **Note:** `private_key_path` only applies when using a self-hosted integration runtime (instead of the default Azure provided runtime), as indicated by supplying a value for `integration_runtime_name`.

* `private_key_passphrase` - (Optional) The passphrase for the private key if the key is encrypted.

* `key_vault_private_key_passphrase` - (Optional) A `key_vault_private_key_passphrase` block as defined below.

~> **Note:** One of `private_key_content_base64` or `private_key_path` (or their Key Vault equivalent) is required when `authentication_type` is set to `SshPublicKey`.

* `host_key_fingerprint` - (Optional) The host key fingerprint of the SFTP server.

* `skip_host_key_validation` - (Optional) Whether to validate host key fingerprint while connecting. If set to `false`, `host_key_fingerprint` must also be set.

---

A `key_vault_password` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the name of the secret containing the password.

---

A `key_vault_private_key_content_base64` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the name of the secret containing the Base64 encoded SSH private key.

---

A `key_vault_private_key_passphrase` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the name of the secret containing the SSH private key passphrase.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Linked Service.

## Import

Data Factory Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_sftp.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
