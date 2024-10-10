---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_dicom_service"
description: |-
  Get information about an existing Healthcare DICOM (Digital Imaging and Communications in Medicine) Service
---

# Data Source: azurerm_healthcare_dicom_service

Use this data source to access information about an existing Healthcare DICOM Service

## Example Usage

```hcl
data "azurerm_healthcare_dicom_service" "example" {
  name         = "example-healthcare_dicom_service"
  workspace_id = data.azurerm_healthcare_workspace.example.id
}

output "azurerm_healthcare_dicom_service" {
  value = data.azurerm_healthcare_dicom_service.example.id
}
```

## Argument Reference

* `name` - The name of the Healthcare DICOM Service

* `workspace_id` - The ID of the Healthcare Workspace in which the Healthcare DICOM Service exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Healthcare DICOM Service.

* `location` - The Azure Region where the Healthcare DICOM Service is located.

* `authentication` - The `authentication` block as defined below.

* `data_partitions_enabled` - If data partitions are enabled or not.

* `cors` - The `cors` block as defined below.

* `encryption_key_url` - The URL of the key to use for encryption as part of the customer-managed key encryption settings.

* `service_url` - The url of the Healthcare DICOM Services.

* `storage` - The `storage` block as defined below.

* `tags` - A map of tags assigned to the Healthcare DICOM Service.

---
An `authentication` exports the following:

* `authority` - The Azure Active Directory (tenant) that serves as the authentication authority to access the service. The default authority is the Directory defined in the authentication scheme in use when running Terraform.
  Authority must be registered to Azure AD and in the following format: <https://{Azure-AD-endpoint}/{tenant-id>}.

* `audience` - The intended audience to receive authentication tokens for the service. The default value is <https://dicom.azurehealthcareapis.azure.com>

---

A `cors` exports the following:

* `allowed_origins` - A list of allowed origins for CORS.

* `allowed_headers` - A list of allowed headers for CORS.

* `allowed_methods` - A list of allowed methods for CORS.

* `max_age_in_seconds` - The maximum age in seconds for the CORS configuration.

* `allow_credentials` - Whether to allow credentials in CORS.

---

A `storage` block exports the following:

* `file_system_name` - The filesystem name of connected storage account.

* `storage_account_id` - The resource ID of connected storage account.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare DICOM Service.
