---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_dicom_service"
description: |-
  Get information about an existing Healthcare Dicom Service
---

# Data Source: azurerm_healthcare_dicom_service

Use this data source to access information about an existing Healthcare Dicom Service

## Example Usage

```hcl
data "azurerm_healthcare_dicom_service" "example" {
  name                = "example-healthcare_dicom_service"
  workspace_id = "example_healthcare_workspace"
}

output "azurerm_healthcare_dicom_service" {
  value = data.azurerm_healthcare_dicom_service.example.id
}
```

## Argument Reference

* `name` - The name of the Healthcare Dicom Service

* `workspace_id` - The name of the Healthcare Workspace in which the Healthcare Dicom Service exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Healthcare Workspace.

* `location` - The Azure Region where the Healthcare Dicom Service is located.

* `authentication_configuration` - The `authentication_configuration` block as defined below.

* `service_url` - The url of the Healthcare Dicom Services.

* `tags` - A map of tags assigned to the Healthcare Dicom Service.

---
An `authentication_configuration` supports the following:

* `authority` - The Azure Active Directory (tenant) that serves as the authentication authority to access the service. The default authority is the Directory defined in the authentication scheme in use when running Terraform.
  Authority must be registered to Azure AD and in the following format: https://{Azure-AD-endpoint}/{tenant-id}.

* `audience` - The intended audience to receive authentication tokens for the service. The default value is https://dicom.azurehealthcareapis.azure.com

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Dicom Service.
