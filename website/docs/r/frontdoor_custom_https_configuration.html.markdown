---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_custom_https_configuration"
description: |-
  Manages the Custom Https Configuration for an Azure Front Door Frontend Endpoint.
---

# azurerm_frontdoor_custom_https_configuration

Manages the Custom Https Configuration for an Azure Front Door Frontend Endpoint..

-> **NOTE:** Defining custom https configurations using a separate `azurerm_frontdoor_custom_https_configuration` resource allows for parallel creation/update.

!> **BREAKING CHANGE:** In order to address the ordering issue we have changed the design on how to retrieve existing sub resources such as frontend endpoints. Existing design will be deprecated and will result in an incorrect configuration. Please refer to the updated documentation below for more information.

!> **BREAKING CHANGE:** The `resource_group_name` field has been removed as of the `v2.58.0` provider release. If the `resource_group_name` field has been defined in your current `azurerm_frontdoor_custom_https_configuration` resource configuration file please remove it else you will receive a `An argument named "resource_group_name" is not expected here.` error. If your pre-existing Front Door instance contained inline `custom_https_configuration` blocks there are additional steps that will need to be completed to succefully migrate your Front Door onto the `v2.58.0` provider which [can be found in this guide](../guides/2.58.0-frontdoor-upgrade-guide.html).

!> **Be Aware:** Azure is rolling out a breaking change on Friday 9th April which may cause issues with the CDN/FrontDoor resources. [More information is available in this Github issue](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11231) - however unfortunately this may necessitate a breaking change to the CDN and FrontDoor resources, more information will be posted [in the Github issue](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11231) as the necessary changes are identified.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "FrontDoorExampleResourceGroup"
  location = "West Europe"
}

data "azurerm_key_vault" "vault" {
  name                = "example-vault"
  resource_group_name = "example-vault-rg"
}

resource "azurerm_frontdoor" "example" {
  name                                         = "example-FrontDoor"
  resource_group_name                          = azurerm_resource_group.example.name
  enforce_backend_pools_certificate_name_check = false

  routing_rule {
    name               = "exampleRoutingRule1"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = ["exampleFrontendEndpoint1"]
    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = "exampleBackendBing"
    }
  }

  backend_pool_load_balancing {
    name = "exampleLoadBalancingSettings1"
  }

  backend_pool_health_probe {
    name = "exampleHealthProbeSetting1"
  }

  backend_pool {
    name = "exampleBackendBing"
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = "exampleLoadBalancingSettings1"
    health_probe_name   = "exampleHealthProbeSetting1"
  }

  frontend_endpoint {
    name      = "exampleFrontendEndpoint1"
    host_name = "example-FrontDoor.azurefd.net"
  }

  frontend_endpoint {
    name      = "exampleFrontendEndpoint2"
    host_name = "examplefd1.examplefd.net"
  }
}

resource "azurerm_frontdoor_custom_https_configuration" "example_custom_https_0" {
  frontend_endpoint_id              = azurerm_frontdoor.example.frontend_endpoints["exampleFrontendEndpoint1"]
  custom_https_provisioning_enabled = false
}

resource "azurerm_frontdoor_custom_https_configuration" "example_custom_https_1" {
  frontend_endpoint_id              = azurerm_frontdoor.example.frontend_endpoints["exampleFrontendEndpoint2"]
  custom_https_provisioning_enabled = true

  custom_https_configuration {
    certificate_source                      = "AzureKeyVault"
    azure_key_vault_certificate_secret_name = "examplefd1"
    azure_key_vault_certificate_vault_id    = data.azurerm_key_vault.vault.id
  }
}
```

## Argument Reference

The `custom_https_configuration` block is also valid inside an `azurerm_frontdoor_custom_https_configuration`, which supports the following arguments: 

* `frontend_endpoint_id` - (Required) The ID of the FrontDoor Frontend Endpoint which this configuration refers to.

* `custom_https_provisioning_enabled` - (Required) Should the HTTPS protocol be enabled for this custom domain associated with the Front Door?

* `custom_https_configuration` - (Optional) A `custom_https_configuration` block as defined above.

---

The `custom_https_configuration` block supports the following:

* `certificate_source` - (Optional) Certificate source to encrypted `HTTPS` traffic with. Allowed values are `FrontDoor` or `AzureKeyVault`. Defaults to `FrontDoor`.

The following attributes are only valid if `certificate_source` is set to `AzureKeyVault`:

* `azure_key_vault_certificate_vault_id` - (Required) The ID of the Key Vault containing the SSL certificate.

* `azure_key_vault_certificate_secret_name` - (Required) The name of the Key Vault secret representing the full certificate PFX.

* `azure_key_vault_certificate_secret_version` - (Optional) The version of the Key Vault secret representing the full certificate PFX. Defaults to `Latest`.

~> **Note:** In order to enable the use of your own custom `HTTPS certificate` you must grant `Azure Front Door Service` access to your key vault. For instuctions on how to configure your `Key Vault` correctly please refer to the [product documentation](https://docs.microsoft.com/en-us/azure/frontdoor/front-door-custom-domain-https#option-2-use-your-own-certificate).


## Attributes Reference

* `id` - The ID of the Azure Front Door Custom Https Configuration.

* `custom_https_configuration` - A `custom_https_configuration` block as defined below.

The `custom_https_configuration` block exports the following:

* `minimum_tls_version` - Minimum client TLS version supported.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating a Custom Https Configuration.
* `update` - (Defaults to 6 hours) Used when updating a Custom Https Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving a Custom Https Configuration.
* `delete` - (Defaults to 6 hours) Used when deleting a Custom Https Configuration.

## Import

Front Door Custom Https Configurations can be imported using the `resource id` of the Front Door Custom Https Configuration, e.g.

```shell
terraform import azurerm_frontdoor_custom_https_configuration.example_custom_https_1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/frontDoors/frontdoor1/customHttpsConfiguration/endpoint1
```
