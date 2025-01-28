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

```hcl
resource "azurerm_device_registry_asset_endpoint_profile" "example" {
  name = "example"
  resource_group_name = "example"
  location = "West Europe"
  extended_location_name = "example"
  extended_location_type = "TODO"
  endpoint_profile_type = "TODO"
  target_address = "TODO"
}
```

## Arguments Reference

The following arguments are supported:

* `endpoint_profile_type` - (Required) Defines the configuration for the connector type that is being used with the endpoint profile.

* `extended_location_name` - (Required) The extended location name.

* `extended_location_type` - (Required) The extended location type.

* `location` - (Required) The Azure Region where the Asset Endpoint Profile should exist. Changing this forces a new Asset Endpoint Profile to be created.

* `name` - (Required) The name which should be used for this Asset Endpoint Profile.

* `resource_group_name` - (Required) The name of the Resource Group where the Asset Endpoint Profile should exist. Changing this forces a new Asset Endpoint Profile to be created.

* `target_address` - (Required) The local valid URI specifying the network address/DNS name of a southbound device. The scheme part of the targetAddress URI specifies the type of the device. The additionalConfiguration field holds further connector type specific configuration.

---

* `additional_configuration` - (Optional) Stringified JSON that contains connectivity type specific further configuration (e.g. OPC UA, Modbus, ONVIF).

* `authentication_method` - (Optional) Defines the method to authenticate the user of the client at the server. Defaults to `Certificate`.

* `discovered_asset_endpoint_profile_ref` - (Optional) Reference to a discovered asset endpoint profile. Populated only if the asset endpoint profile has been created from discovery flow. Discovered asset endpoint profile name must be provided.

* `tags` - (Optional) A mapping of tags which should be assigned to the Asset Endpoint Profile.

* `username_password_credentials_password_secret_name` - (Optional) The name of the secret containing the password for authentication mode UsernamePassword.

* `username_password_credentials_username_secret_name` - (Optional) The name of the secret containing the username for authentication mode UsernamePassword.

* `x509_credentials_certificate_secret_name` - (Optional) The name of the secret containing the certificate and private key (e.g. stored as .der/.pem or .der/.pfx) for authentication mode Certificate.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Asset Endpoint Profile.

* `provisioning_state` - Provisioning state of the resource.

* `status` - Read only object to reflect changes that have occurred on the Edge. Similar to Kubernetes status property for custom resources.

* `type` - Azure resource type. Defaults to `Microsoft.DeviceRegistry/AssetEndpointProfiles`.

* `uuid` - Globally unique, immutable, non-reusable id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Asset Endpoint Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Asset Endpoint Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Asset Endpoint Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Asset Endpoint Profile.

## Import

Asset Endpoint Profile can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_device_registry_asset_endpoint_profile.example C:/Program Files/Git/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/adr-terraform-rg/providers/Microsoft.DeviceRegistry/assetendpointprofiles/test-asset-endpoint-profile
```