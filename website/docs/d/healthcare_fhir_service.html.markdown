---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_fhir_service"
description: |-
  Get information about an existing Healthcare FHIR (Fast Healthcare Interoperability Resources) Service.
---

# Data Source: azurerm_healthcare_fhir_service

Use this data source to access information about an existing Healthcare FHIR Service(Fast Healthcare Interoperability Resources).

## Example Usage

```hcl
data "azurerm_healthcare_fhir_service" "example" {
  name         = "example-healthcare"
  workspace_id = data.azurerm_healthcare_fhir_service.example.workspace_id
}

output "healthcare_fhir_service_id" {
  value = data.azurerm_healthcare_fhir_service.example.id
}
```

## Argument Reference

* `name` - The name of the Healthcare FHIR Service.

* `workspace_id` - The name of the Healthcare Workspace in which the Healthcare FHIR Service exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Healthcare FHIR Service.

* `location` - The Azure Region where the Healthcare FHIR Service is located.

* `kind` - The kind of the Healthcare FHIR Service. 

* `identity` - The `identity` block as defined below.

* `access_policy_object_ids` - The list of the access policies of the service instance.

* `cors` - The `cors` block as defined below.

* `container_registry_login_server_url` - The list of azure container registry settings used for convert data operation of the service instance.

* `authentication` - The `authentication` block as defined below.

* `configuration_export_storage_account_name` - The name of the storage account which the operation configuration information is exported to.

* `public_network_access_enabled` - Is public networks access enabled when data plane traffic coming from public networks while private endpoint is enabled?

* `tags` - The map of tags assigned to the Healthcare FHIR Service.

---
An `identity` block exports the following:

* `type` - The type of identity used for the Healthcare FHIR service.

* `principal_id` - The Principal ID associated with this System Assigned Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this System Assigned Managed Service Identity.

---
A `cors` block exports the following:

* `allowed_origins` - The set of origins to be allowed via CORS.
* `allowed_headers` - The set of headers to be allowed via CORS.
* `allowed_methods` - The methods to be allowed via CORS.
* `max_age_in_seconds` - The max age to be allowed via CORS.
* `credentials_allowed` - Are credentials allowed via CORS?

---
An `authentication` block exports the following:

* `authority` - The Azure Active Directory (tenant) that serves as the authentication authority to access the service. The default authority is the Directory defined in the authentication scheme in use when running Terraform.
  Authority must be registered to Azure AD and in the following format: https://{Azure-AD-endpoint}/{tenant-id}.
* `audience` - The intended audience to receive authentication tokens for the service. The default value is https://<name>.fhir.azurehealthcareapis.com

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare FHIR Service.

