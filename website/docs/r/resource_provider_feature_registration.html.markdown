---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_provider_feature_registration"
description: |-
    Manages the Registration of a Resource Provider Feature.
---

# azurerm_resource_provider_feature_registration

Manages the registration of a feature within a Resource Provider.

-> **Note:** The Resource Provider must be registered before a feature can be registered.

~> **Note:** Only Preview Features which have an `ApprovalType` of `AutoApproval` can be managed in Terraform, features which require manual approval by Service Teams are unsupported. [More information on Resource Provider Preview Features can be found in this document](https://docs.microsoft.com/rest/api/resources/features)

## Example Usage

```hcl
resource "azurerm_resource_provider_feature_registration" "example" {
  name          = "EncryptionAtHost"
  provider_name = "Microsoft.Compute"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the feature to register. Changing this forces a new resource to be created.

* `provider_name` - (Required) The Resource Provider name. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Provider Feature.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when registering the Resource Provider Feature.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Provider Feature.
* `delete` - (Defaults to 30 minutes) Used when unregistering the Resource Provider Feature.

## Import

Features Registrations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_provider_feature_registration.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Features/providers/{ResourceProviderName}/features/{FeatureName}
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Features` - 2021-07-01
