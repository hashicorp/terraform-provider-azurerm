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

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service.

The following supported arguments are specific to SFTP Linked Service:

* `authentication_type` - (Required) The type of authentication used to connect to the web table source. Valid options are `MultiFactor`, `Basic` and `SshPublicKey`.

* `host` - (Required) The SFTP server hostname.

* `port` - (Required) The TCP port number that the SFTP server uses to listen for client connection. Default value is 22.

* `username` - (Required) The username used to log on to the SFTP server.

* `password` - (Required if using Basic Authentication) Password to log on to the SFTP Server for Basic Authentication.

* `private_key_content` - (Required if using `SshPublicKey` Authentication and not `private_key_path`) The private key content in OpenSSH format or the URL of an Azure Key Vault Secret of the same used to log on to the SFTP server.

* `private_key_path` - (Required if using `SshPublicKey` Authentication, a self-hosted integration runtime, and not `private_key_content`) The absolute path to the private key file that the integration runtime can access. This applies only when the self-hosted type of integration runtime is specified in `integration_runtime_name`.

* `host_key_fingerprint` - (Optional) The host key fingerprint of the SFTP server.

* `skip_host_key_validation` - (Optional) Whether to validate host key fingerprint while connecting. If set to `false`, `host_key_fingerprint` must also be set.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Linked Service.

## Import

Data Factory Linked Service's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_sftp.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
