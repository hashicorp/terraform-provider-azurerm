subcategory	layout	page_title	description
IoT Operations	azurerm	Azure Resource Manager: azurerm_iotoperations_instance	Manages an IoT Operations Instance.

azurerm_iotoperations_instance
Manages an IoT Operations Instance.

Example Usage
----------------
resource "azurerm_resource_group" "example" {
  name     = "example-iotoperations"
  location = "West Europe"
}

resource "azurerm_iotoperations_instance" "example" {
  name                = "example-iotinstance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  # Optional extended location
  extended_location_name = "your-extended-location"
  extended_location_type = "CustomLocation"

  description = "Example IoT Operations instance"
  version     = "1.0.0"

  tags = {
    Environment = "Dev"
    Owner       = "team"
  }
}

Arguments Reference
-------------------
The following arguments are supported:

- name - (Required) The name of the IoT Operations Instance. Changing this forces a new Instance.

- resource_group_name - (Required) The name of the Resource Group where the IoT Operations Instance should exist. Changing this forces a new Instance.

- location - (Required) The Azure Region where the IoT Operations Instance should exist. Changing this forces a new Instance.

- extended_location_name - (Optional) Name of an extended location to place the resource in.

- extended_location_type - (Optional) Type of the extended location (string).

- description - (Optional) Description for the IoT Operations Instance.

- version - (Optional) Version information for the instance.

- tags - (Optional) A mapping of tags which should be assigned to the Instance.

Attributes Reference
--------------------
In addition to the Arguments above, the following attributes are exported:

- id - The ID of the IoT Operations Instance.

- provisioning_state - (Computed) The provisioning state of the instance.

- location - (Computed) The Azure region of the resource.

- extended_location_name - (Computed) When present, the name of the extended location.

- extended_location_type - (Computed) When present, the type of the extended location.

- description - (Computed) The description stored on the instance (if present).

- version - (Computed) The version stored on the instance (if present).

- tags - (Computed) Tags assigned to the instance.

Timeouts
--------
This resource does not expose a timeouts block in the provider schema; operations use provider defaults.

Import
------
IoT Operations Instances can be imported using the resource id, e.g.:

terraform import azurerm_iotoperations_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.IoTOperations/instances/example-instance

API Providers
-------------
This resource uses the following Azure API Providers:

- Microsoft.IoTOperations - 2024-11-01

