---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_assessment_metadata"
description: |-
  Manages the Security Center Assessment Metadata for Azure Security Center.
---

# azurerm_security_center_assessment_metadata

Manages the Security Center Assessment Metadata for Azure Security Center.

~> **NOTE:** This resource has been deprecated in favour of the `azurerm_security_center_assessment_policy` resource and will be removed in the next major version of the AzureRM Provider. The new resource shares the same fields as this one, and information on migrating across [can be found in this guide](../guides/migrating-between-renamed-resources.html).

## Example Usage

```hcl
resource "azurerm_security_center_assessment_metadata" "example" {
  display_name = "Test Display Name"
  severity     = "Medium"
  description  = "Test Description"
}
```

## Arguments Reference

The following arguments are supported:

* `description` - (Required) The description of the Security Center Assessment.

* `display_name` - (Required) The user-friendly display name of the Security Center Assessment.

* `severity` - (Required) The severity level of the Security Center Assessment. Possible values are `Low`, `Medium` and `High`. Defaults to `Medium`.

---

* `implementation_effort` - (Optional) The implementation effort which is used to remediate the Security Center Assessment. Possible values are `Low`, `Moderate` and `High`.

* `remediation_description` - (Optional) The description which is used to mitigate the security issue.

* `threats` - (Optional) A list of the threat impacts for the Security Center Assessment. Possible values are `AccountBreach`, `DataExfiltration`, `DataSpillage`, `DenialOfService`, `ElevationOfPrivilege`, `MaliciousInsider`, `MissingCoverage` and `ThreatResistance`.

* `user_impact` - (Optional) The user impact of the Security Center Assessment. Possible values are `Low`, `Moderate` and `High`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Security Center Assessment Metadata.

* `name` - The GUID as the name of the Security Center Assessment Metadata.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Security Center Assessment Metadata.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Assessment Metadata.
* `update` - (Defaults to 30 minutes) Used when updating the Security Center Assessment Metadata.
* `delete` - (Defaults to 30 minutes) Used when deleting the Security Center Assessment Metadata.

## Import

Security Assessments Metadata can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_assessment_metadata.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/assessmentMetadata/metadata1
```
