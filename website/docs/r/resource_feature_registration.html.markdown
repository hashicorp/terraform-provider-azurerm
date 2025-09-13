---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_provider_registration"
description: |-
    Manages the Registration of a Resource Provider.
---

# azurerm_resource_feature_registration

Manages the registration of a feature within a Resource Provider - which allows to active features of an Resource Provider.

-> **Note:** The Resource Provider must be registered before a feature can be registered.
~> **Note:** Only Preview Features which have an `ApprovalType` of `AutoApproval` can be managed in Terraform, features which require manual approval by Service Teams are unsupported. [More information on Resource Provider Preview Features can be found in this document](https://docs.microsoft.com/rest/api/resources/features)

## Example Usage

```hcl
resource "azurerm_resource_feature_registration" "example" {
  name          = "EncryptionAtHost"
  provider_name = "Microsoft.Compute"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the feature to register
* `provider_name` - (Required) The namespace of the Resource Provider which should be registered. Changing this forces a new resource to be created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Resource Provider/Features.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Provider.

## Import

Features Registrations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_provider_registration.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Features/featureProviders/Microsoft.Compute/subscriptionFeatureRegistrations/EncryptionAtHost
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Features` - 2021-07-01
