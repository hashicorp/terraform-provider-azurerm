---
subcategory: "IoT Operations"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotoperations_instance"
description: |-
  Manages an Azure IoT Operations Instance.
---

# azurerm_iotoperations_instance

Manages an Azure IoT Operations Instance.

An IoT Operations Instance is the core management resource that provides orchestration and runtime capabilities for IoT workloads. It serves as the foundational platform for deploying brokers, dataflows, and other IoT components, managing schema registries, and coordinating edge computing operations.

## Example Usage

### Basic IoT Operations Instance

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iotoperations_instance" "example" {
  name                   = "example-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/example-registry"
  extended_location_name = "microsoftiotoperations"
  extended_location_type = "CustomLocation"
}
```

### IoT Operations Instance with Description and Version

```hcl
resource "azurerm_iotoperations_instance" "production" {
  name                   = "production-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/prod-registry"
  extended_location_name = "production-edge-cluster"
  extended_location_type = "CustomLocation"
  
  description = "Production IoT Operations instance for manufacturing edge computing"
  version     = "1.0.0"

  tags = {
    Environment = "Production"
    Department  = "Manufacturing"
    CostCenter  = "CC-1234"
  }
}
```

### Development Environment Instance

```hcl
resource "azurerm_iotoperations_instance" "development" {
  name                   = "dev-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/dev-registry"
  extended_location_name = "development-edge-cluster"
  extended_location_type = "CustomLocation"
  
  description = "Development environment for IoT Operations testing and validation"
  
  tags = {
    Environment = "Development"
    Team        = "IoT-Development"
    Project     = "EdgeComputing"
  }
}
```

### Multi-Region Manufacturing Setup

```hcl
# Primary manufacturing site
resource "azurerm_iotoperations_instance" "manufacturing_primary" {
  name                   = "mfg-primary-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/manufacturing-registry"
  extended_location_name = "primary-manufacturing-site"
  extended_location_type = "CustomLocation"
  
  description = "Primary manufacturing site IoT Operations instance"
  version     = "2.0.0"

  tags = {
    Environment = "Production"
    Site        = "Primary"
    Region      = "WestEurope"
    Application = "Manufacturing"
  }
}

# Secondary manufacturing site
resource "azurerm_iotoperations_instance" "manufacturing_secondary" {
  name                   = "mfg-secondary-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = "East US"
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/manufacturing-registry"
  extended_location_name = "secondary-manufacturing-site"
  extended_location_type = "CustomLocation"
  
  description = "Secondary manufacturing site IoT Operations instance"
  version     = "2.0.0"

  tags = {
    Environment = "Production"
    Site        = "Secondary"
    Region      = "EastUS"
    Application = "Manufacturing"
  }
}
```

### Edge Computing Instance for Retail

```hcl
resource "azurerm_iotoperations_instance" "retail_edge" {
  name                   = "retail-edge-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/retail-registry"
  extended_location_name = "retail-store-cluster"
  extended_location_type = "CustomLocation"
  
  description = "Retail edge computing instance for in-store IoT operations"
  
  tags = {
    Environment = "Production"
    Industry    = "Retail"
    StoreType   = "Flagship"
    Location    = "Downtown"
  }
}
```

### Smart City Infrastructure Instance

```hcl
resource "azurerm_iotoperations_instance" "smart_city" {
  name                   = "smartcity-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/smartcity-registry"
  extended_location_name = "city-infrastructure-cluster"
  extended_location_type = "CustomLocation"
  
  description = "Smart city infrastructure management and monitoring instance"
  version     = "3.0.0"

  tags = {
    Environment   = "Production"
    Application   = "SmartCity"
    Municipality  = "CityName"
    Department    = "PublicWorks"
    Compliance    = "SOC2"
  }
}
```

### Healthcare IoT Instance

```hcl
resource "azurerm_iotoperations_instance" "healthcare" {
  name                   = "healthcare-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/healthcare-registry"
  extended_location_name = "hospital-edge-cluster"
  extended_location_type = "CustomLocation"
  
  description = "Healthcare IoT Operations for medical device management and patient monitoring"
  
  tags = {
    Environment = "Production"
    Industry    = "Healthcare"
    Compliance  = "HIPAA"
    Facility    = "HospitalMain"
    Department  = "IT"
  }
}
```

### Energy Management Instance

```hcl
resource "azurerm_iotoperations_instance" "energy_management" {
  name                   = "energy-mgmt-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/energy-registry"
  extended_location_name = "power-plant-cluster"
  extended_location_type = "CustomLocation"
  
  description = "Energy management and grid monitoring IoT Operations instance"
  version     = "1.5.0"

  tags = {
    Environment = "Production"
    Industry    = "Energy"
    FacilityType = "PowerPlant"
    GridZone    = "Zone-A"
    Criticality = "High"
  }
}
```

### Multi-Tenant SaaS Instance

```hcl
resource "azurerm_iotoperations_instance" "saas_platform" {
  name                   = "saas-platform-instance"
  resource_group_name    = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  schema_registry_ref   = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/saas-registry"
  extended_location_name = "saas-edge-cluster"
  extended_location_type = "CustomLocation"
  
  description = "Multi-tenant SaaS platform IoT Operations instance"
  
  tags = {
    Environment = "Production"
    ServiceType = "SaaS"
    Tier        = "Premium"
    Scaling     = "Auto"
    Monitoring  = "24x7"
  }
}
```

### Minimal Configuration for Testing

```hcl
resource "azurerm_iotoperations_instance" "minimal" {
  name                = "test-instance"
  resource_group_name = azurerm_resource_group.example.name
  location           = azurerm_resource_group.example.location
  schema_registry_ref = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DeviceRegistry/schemaRegistries/test-registry"
  
  # Minimal configuration without extended location for basic testing
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the IoT Operations Instance. Must be between 3-63 characters, lowercase alphanumeric with dashes, starting and ending with alphanumeric characters. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the IoT Operations Instance should exist. Must be between 1-90 characters. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the IoT Operations Instance should exist. Changing this forces a new resource to be created.

* `schema_registry_ref` - (Required) The resource ID reference to the Device Registry Schema Registry that this IoT Operations Instance will use for schema management. This registry stores and manages schemas for data validation and transformation. Changing this forces a new resource to be created.

* `description` - (Optional) A description of the IoT Operations Instance. Use this to document the purpose, environment, or specific use case of the instance.

* `version` - (Optional) The version of the IoT Operations Instance. If not specified, Azure will assign a default version and may update it during the resource lifecycle.

* `extended_location_name` - (Optional) The name of the extended location where the IoT Operations Instance should be deployed. This is typically a Custom Location representing an edge cluster or on-premises infrastructure. Changing this forces a new resource to be created.

* `extended_location_type` - (Optional) The type of the extended location. Must be `CustomLocation` when specified. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the IoT Operations Instance.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ID of the IoT Operations Instance.

* `provisioning_state` - The provisioning state of the IoT Operations Instance. Possible values include `Succeeded`, `Failed`, `Canceled`, `Creating`, `Updating`, and `Deleting`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Operations Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Operations Instance.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Operations Instance.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Operations Instance.

## Import

An IoT Operations Instance can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotoperations_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.IoTOperations/instances/instance1
```
