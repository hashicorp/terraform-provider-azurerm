---
subcategory: "Extended Location"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_extended_location_custom_location"
description: |-
  Manages a Custom Location within an Extended Location.
---

# azurerm_extended_location_custom_location

Manages a Custom Location within an Extended Location.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_arc_kubernetes_cluster" "example" {
  name                         = "example-akcc"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = "West Europe"
  agent_public_key_certificate = filebase64("testdata/public.cer")

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}

resource "azurerm_arc_kubernetes_cluster_extension" "example" {
  name           = "example-ext"
  cluster_id     = azurerm_arc_kubernetes_cluster.example.id
  extension_type = "microsoft.flux"
}

resource "azurerm_extended_location_custom_location" "example" {
  name                = "example-custom-location"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  cluster_extension_ids = [
    azurerm_arc_kubernetes_cluster_extension.example.id
  ]
  display_name     = "example-custom-location"
  namespace        = "example-namespace"
  host_resource_id = azurerm_arc_kubernetes_cluster.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Custom Location. Changing this forces a new Custom Location to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Custom Location should exist. Changing this forces a new Custom Location to be created.

* `location` - (Required) Specifies the Azure location where the Custom Location should exist. Changing this forces a new Custom Location to be created.

* `namespace` - (Required) Specifies the namespace of the Custom Location.Changing this forces a new Custom Location to be created.

* `cluster_extension_ids` - (Required) Specifies the list of Cluster Extension IDs.

* `host_resource_id` - (Required) Specifies the host resource ID.

* `authentication` - (Optional) An `authentication` block as defined below.

* `display_name` - (Optional) Specifies the display name of the Custom Location.

* `host_type` - (Optional) Specifies the host type of the Custom Location. The only possible values is `KubernetesCluster`.

---

An `authentication` block supports the following:

* `type` - (Required) Specifies the type of authentication.

* `value` - (Required) Specifies the value of authentication.

## Attributes Reference

* `id` - The ID of the Custom Location.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Custom Location.
* `read` - (Defaults to 5 minutes) Used when retrieving the Custom Location.
* `update` - (Defaults to 30 minutes) Used when updating the Custom Location.
* `delete` - (Defaults to 30 minutes) Used when deleting the Custom Location.

## Import

Custom Locations can be imported using the resource id, e.g.

```shell
terraform import azurerm_extended_location_custom_location.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.ExtendedLocation/customLocations/example-custom-location
```

