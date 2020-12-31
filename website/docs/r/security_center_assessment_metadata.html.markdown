---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_assessment_metadata"
description: |-
  Manages the Security Center Assessment Metadata for Azure Security Center.
---

# azurerm_security_center_assessment_metadata

Manages the Security Center Assessment Metadata for Azure Security Center.

## Example Usage

```hcl
resource "azurerm_security_center_assessment_metadata" "example" {
  name            = "ca039e75-a276-4175-aebc-bcd41e4b14b7"
  display_name    = "Test Display Name"
  assessment_type = "CustomerManaged"
  severity        = "Medium"
  description     = "Test Description"
  categories      = ["Compute"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The GUID as name which should be used for this Security Center Assessment Metadata. Changing this forces a new Security Center Assessment Metadata to be created.

* `assessment_type` - (Required) The type of the Security Center Assessment. Possible values are `BuiltIn`, `CustomPolicy`, `CustomerManaged` and `VerifiedPartner`.

* `description` - (Required) The description of the Security Center Assessment.

* `display_name` - (Required) The user friendly display name of the Security Center Assessment.

* `severity` - (Required) The severity level of the Security Center Assessment. Possible values are `Low`, `Medium` and `High`.

---

* `categories` - (Optional) A list of the categories which are at risk when the Security Center Assessment is unhealthy. Possible values are `Compute`, `Data`, `IdentityAndAccess`, `IoT` and `Networking`.

* `implementation_effort` - (Optional) The implementation effort which is used to remediate this assessment. Possible values are `Low`, `Moderate` and `High`.

* `partner_data` - (Optional)  A `partner_data` block as defined below.

* `preview` - (Optional) Is this assessment in preview release status?

* `remediation_description` - (Optional) The description which is used to mitigate this security issue.

* `threats` - (Optional) A list of the threat impacts for the Security Center Assessment. Possible values are `AccountBreach`, `DataExfiltration`, `DataSpillage`, `DenialOfService`, `ElevationOfPrivilege`, `MaliciousInsider`, `MissingCoverage` and `ThreatResistance`.

* `user_impact` - (Optional) The user impact of the Security Center Assessment. Possible values are `Low`, `Moderate` and `High`.

---

An `partner_data` block exports the following:

* `partner_name` - (Required) The name of the company of the partner.

* `secret` - (Required) The secret which is used to authenticate the partner and verify if it created the Security Center Assessment with write-only.

---

* `product_name` - (Optional) The name of the product of the partner that created the Security Center Assessment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Security Center Assessment Metadata.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Security Center Assessment Metadata.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Assessment Metadata.
* `update` - (Defaults to 30 minutes) Used when updating the Security Center Assessment Metadata.
* `delete` - (Defaults to 30 minutes) Used when deleting the Security Center Assessment Metadata.

## Import

security AssessmentsMetadatums can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_assessment_metadata.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/assessmentMetadata/metadata1
```
