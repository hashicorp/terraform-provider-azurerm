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

* `workspace_id` - The id of the Healthcare Workspace in which the Healthcare DICOM Service exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Healthcare DICOM Service.

* `location` - The Azure Region where the Healthcare DICOM Service is located.

* `authentication` - The `authentication` block as defined below.

* `service_url` - The url of the Healthcare DICOM Services.

* `tags` - A map of tags assigned to the Healthcare DICOM Service.

---
An `authentication` supports the following:

* `authority` - The Azure Active Directory (tenant) that serves as the authentication authority to access the service. The default authority is the Directory defined in the authentication scheme in use when running Terraform.
  Authority must be registered to Azure AD and in the following format: <https://{Azure-AD-endpoint}/{tenant-id>}.

* `audience` - The intended audience to receive authentication tokens for the service. The default value is <https://dicom.azurehealthcareapis.azure.com>

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare DICOM Service.
