---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_integration_runtime_self_hosted"
description: |-
  Manages a Data Factory Self-hosted Integration Runtime.
---

# azurerm_data_factory_integration_runtime_self_hosted

Manages a Data Factory Self-hosted Integration Runtime.

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

resource "azurerm_data_factory_integration_runtime_self_hosted" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Linked Service with. Changing this forces a new resource.

* `name` - (Required) The name which should be used for this Data Factory. Changing this forces a new Data Factory Self-hosted Integration Runtime to be created.

---

* `description` - (Optional) Integration runtime description.

* `rbac_authorization` - (Optional) A `rbac_authorization` block as defined below. Changing this forces a new resource to be created.

* `self_contained_interactive_authoring_enabled` - (Optional) Specifies whether enable interactive authoring function when your self-hosted integration runtime is unable to establish a connection with Azure Relay.

---

A `rbac_authorization` block supports the following:

* `resource_id` - (Required) The resource identifier of the integration runtime to be shared.

-> **Note:** RBAC Authorization creates a [linked Self-hosted Integration Runtime targeting the Shared Self-hosted Integration Runtime in resource_id](https://docs.microsoft.com/azure/data-factory/create-shared-self-hosted-integration-runtime-powershell#share-the-self-hosted-integration-runtime-with-another-data-factory). The linked Self-hosted Integration Runtime needs Contributor access granted to the Shared Self-hosted Data Factory. See example [Shared Self-hosted](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/data-factory/shared-self-hosted).

For more information on the configuration, please check out the [Azure documentation](https://docs.microsoft.com/rest/api/datafactory/integrationruntimes/createorupdate#linkedintegrationruntimerbacauthorization)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory.

* `primary_authorization_key` - The primary integration runtime authentication key.

* `secondary_authorization_key` - The secondary integration runtime authentication key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory.

## Import

Data Factories can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_integration_runtime_self_hosted.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/integrationruntimes/example
```
