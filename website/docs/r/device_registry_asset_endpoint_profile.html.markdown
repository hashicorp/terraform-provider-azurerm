---
subcategory: "Device Registry"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_device_registry_asset_endpoint_profile"
description: |-
  Manages a Device Registry Asset Endpoint Profile.
---

# azurerm_device_registry_asset_endpoint_profile

Manages a Device Registry Asset Endpoint Profile.

## Example Usage

### Asset Endpoint Profile with Certificate authentication
```hcl
resource "azurerm_device_registry_asset_endpoint_profile" "example" {
  name                                        = "example"
  location                                    = "West US 2"
  endpoint_profile_type                       = "OpcUa"
  resource_group_id                           = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group"
  extended_location_id                        = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group/providers/Microsoft.ExtendedLocation/customLocations/my-custom-location"
  target_address                              = "opc.tcp://foo"
  discovered_asset_endpoint_profile_reference = "discoveredAssetEndpointProfile123"
  additional_configuration                    = "{\"foo\": \"bar\"}"
  authentication {
    method                                  = "Certificate"
    x509_credential_certificate_secret_name = "myCertificateRef"
  }
  tags = {
    "sensor" = "temperature,humidity"
  }
}
```

### Asset Endpoint Profile with UsernamePassword authentication
```hcl
resource "azurerm_device_registry_asset_endpoint_profile" "example" {
  name                                        = "example"
  location                                    = "West US 2"
  endpoint_profile_type                       = "OpcUa"
  resource_group_id                           = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group"
  extended_location_id                        = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group/providers/Microsoft.ExtendedLocation/customLocations/my-custom-location"
  target_address                              = "opc.tcp://foo"
  discovered_asset_endpoint_profile_reference = "discoveredAssetEndpointProfile123"
  additional_configuration                    = "{\"foo\": \"bar\"}"
  authentication {
    method                                            = "UsernamePassword"
    username_password_credential_username_secret_name = "myUsernameRef"
    username_password_credential_password_secret_name = "myPasswordRef"
  }
  tags = {
    "sensor" = "temperature,humidity"
  }
}
```

### Asset Endpoint Profile with Anonymous authentication
```hcl
resource "azurerm_device_registry_asset_endpoint_profile" "example" {
  name                                        = "example"
  location                                    = "West US 2"
  endpoint_profile_type                       = "OpcUa"
  resource_group_id                           = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group"
  extended_location_id                        = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group/providers/Microsoft.ExtendedLocation/customLocations/my-custom-location"
  target_address                              = "opc.tcp://foo"
  discovered_asset_endpoint_profile_reference = "discoveredAssetEndpointProfile123"
  additional_configuration                    = "{\"foo\": \"bar\"}"
  authentication {
    method = "Anonymous"
  }
  tags = {
    "sensor" = "temperature,humidity"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `endpoint_profile_type` - (Required) Defines the configuration for the connector type that is being used with the endpoint profile.

* `extended_location_id` - (Required) The ID of the extended location. Must provide a custom location ID.

* `location` - (Required) The Azure Region where the Device Registry Asset Endpoint Profile should exist. Changing this forces a new Device Registry Asset Endpoint Profile to be created.

* `name` - (Required) The name which should be used for this Device Registry Asset Endpoint Profile. Changing this forces a new Device Registry Asset Endpoint Profile to be created.

* `resource_group_id` - (Required) The ID of the Resource Group where the Device Registry Asset Endpoint Profile should exist.

* `target_address` - (Required) The local valid URI specifying the network address/DNS name of a southbound device. The scheme part of the targetAddress URI specifies the type of the device. The additionalConfiguration field holds further connector type specific configuration.

---

* `additional_configuration` - (Optional) Stringified JSON that contains connectivity type specific further configuration (e.g. OPC UA, Modbus, ONVIF).

* `authentication` - (Optional) A `authentication` block as defined below. Defines the client authentication mechanism to the server.

* `discovered_asset_endpoint_profile_reference` - (Optional) Reference to a discovered asset endpoint profile. Populated only if the asset endpoint profile has been created from discovery flow. Discovered asset endpoint profile name must be provided.

* `tags` - (Optional) A mapping of tags which should be assigned to the Device Registry Asset Endpoint Profile.

---

A `authentication` block supports the following:

* `method` - (Required) Defines the method to authenticate the user of the client at the server. Possible values are `Certificate`, `UsernamePassword`, and `Anonymous`.

* `username_password_credential_password_secret_name` - (Optional) The name of the secret containing the password for authentication mode UsernamePassword.

* `username_password_credential_username_secret_name` - (Optional) The name of the secret containing the username for authentication mode UsernamePassword.

* `x509_credential_certificate_secret_name` - (Optional) The name of the secret containing the certificate and private key (e.g. stored as .der/.pem or .der/.pfx) for authentication mode Certificate.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Device Registry Asset Endpoint Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Device Registry Asset Endpoint Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Device Registry Asset Endpoint Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Device Registry Asset Endpoint Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Device Registry Asset Endpoint Profile.

## Import

Device Registry Asset Endpoint Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_device_registry_asset_endpoint_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.DeviceRegistry/assetEndpointProfiles/assetEndpointProfileName
```
