---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_custom_https_configuration"
description: |-
  Manages the Custom Https Configuration for an Azure Front Door Frontend Endpoint.
---

# azurerm_frontdoor_custom_https_configuration

Manages the Custom Https Configuration for an Azure Front Door Frontend Endpoint..

~> **NOTE:** Custom https configurations for a Front Door Frontend Endpoint can be defined both within [the `azurerm_frontdoor` resource](frontdoor.html) via the `custom_https_configuration` block and by using a separate resource, as described in the following sections.

-> **NOTE:** Defining custom https configurations using a separate `azurerm_frontdoor_custom_https_configuration` resource allows for parallel creation/update.
 

```hcl
resource "azurerm_resource_group" "example" {
  name     = "FrontDoorExampleResourceGroup"
  location = "EastUS2"
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
  frontend_endpoint_id              = azurerm_frontdoor.example.frontend_endpoint[0].id
  custom_https_provisioning_enabled = false
}

resource "azurerm_frontdoor_custom_https_configuration" "example_custom_https_1" {
  frontend_endpoint_id              = azurerm_frontdoor.example.frontend_endpoint[1].id
  custom_https_provisioning_enabled = true

  custom_https_configuration {
    certificate_source                         = "AzureKeyVault"
    azure_key_vault_certificate_secret_name    = "examplefd1"
    azure_key_vault_certificate_secret_version = "ec8d0737e0df4f4gb52ecea858e97a73"
    azure_key_vault_certificate_vault_id       = data.azurerm_key_vault.vault.id
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

* `azure_key_vault_certificate_secret_version` - (Required) The version of the Key Vault secret representing the full certificate PFX.

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

Front Door Custom Https Configurations can be imported using the `resource id` of the Frontend Endpoint, e.g.

```shell
terraform import azurerm_frontdoor_custom_https_configuration.example_custom_https_1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/frontdoors/frontdoor1/frontendEndpoints/endpoint1
```
