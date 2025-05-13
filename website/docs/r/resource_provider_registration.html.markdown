---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_provider_registration"
description: |-
    Manages the Registration of a Resource Provider.
---

# azurerm_resource_provider_registration

Manages the registration of a Resource Provider - which allows access to the API's supported by this Resource Provider.

-> **Note:** The Azure Provider will automatically register all of the Resource Providers which it supports on launch (unless opted-out using the `skip_provider_registration` field within the provider block).

!> **Note:** The errors returned from the Azure API when a Resource Provider is unregistered are unclear (example `API version '2019-01-01' was not found for 'Microsoft.Foo'`) - please ensure that all of the necessary Resource Providers you're using are registered - if in doubt **we strongly recommend letting Terraform register these for you**.

-> **Note:** Adding or Removing a Preview Feature will re-register the Resource Provider.

## Example Usage

```hcl
resource "azurerm_resource_provider_registration" "example" {
  name = "Microsoft.PolicyInsights"
}
```

## Example Usage (Registering a Preview Feature)

```hcl
provider "azurerm" {
  features {}

  skip_provider_registration = true
}

resource "azurerm_resource_provider_registration" "example" {
  name = "Microsoft.ContainerService"

  feature {
    name       = "AKS-DataPlaneAutoApprove"
    registered = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The namespace of the Resource Provider which should be registered. Changing this forces a new resource to be created.

* `feature` - (Optional) A list of `feature` blocks as defined below.

~> **Note:** The `feature` block allows a Preview Feature to be explicitly Registered or Unregistered for this Resource Provider - once a Feature has been explicitly Registered or Unregistered, it must be specified in the Terraform Configuration (it's not possible to reset this to the default, unspecified, state).

---

A `feature` block supports the following:

* `name` - (Required) Specifies the name of the feature to register.

~> **Note:** Only Preview Features which have an `ApprovalType` of `AutoApproval` can be managed in Terraform, features which require manual approval by Service Teams are unsupported. [More information on Resource Provider Preview Features can be found in this document](https://docs.microsoft.com/rest/api/resources/features)

* `registered` - (Required) Should this feature be Registered or Unregistered?

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when registering the Resource Provider/Features.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Provider.
* `update` - (Defaults to 2 hours) Used when updating the Resource Provider/Features.
* `delete` - (Defaults to 30 minutes) Used when unregistering the Resource Provider.

## Import

Resource Provider Registrations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_provider_registration.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.PolicyInsights
```
