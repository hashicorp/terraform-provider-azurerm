---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_metadata_schema"
subcategory: "ApiCenter"
description: |-
  Documentation for azurerm_api_center_metadata_schema.
---

# azurerm_api_center_metadata_schema

Documentation for azurerm_api_center_metadata_schema.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_api_center_service" "example" {
  name                = "example-service"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_api_center_metadata_schema" "example" {
  name       = "example-schema"
  service_id = azurerm_api_center_service.example.id

  metadata_assignment {
    entity = "api"
  }

  schema = jsonencode({
    type = "object"
    properties = {
      name = {
        type = "string"
      }
    }
  })
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the metadata schema. Changing this value forces a new resource to be created.

* `metadata_assignment` - (Required) (see [nested schema](#schema-6e6573746564-696e707574-6d657461646174615f61737369676e6d656e74))

* `schema` - (Required) The JSON schema defining the metadata.

* `service_id` - (Required) The ID of the API Center service. Changing this value forces a new resource to be created.

---

<a id="schema-6e6573746564-696e707574-6d657461646174615f61737369676e6d656e74"></a>
The `metadata_assignment` block supports the following:

* `entity` - (Required) The type of entity to which the metadata applies. Possible values include `api`, `deployment`, and `environment`.

* `deprecated` - (Optional) Whether the metadata assignment is deprecated.

* `required` - (Optional) Whether the metadata assignment is required.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

API Center Metadata Schemas can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_center_metadata_schema.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiCenter/services/example-service/metadataSchemas/example-schema
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiCenter` - 2024-03-01
