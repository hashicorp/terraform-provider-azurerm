---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_fhir_service"
description: |-
  Manages a Healthcare FHIR (Fast Healthcare Interoperability Resources) Service.
---

# azurerm_healthcare_fhir_service

Manages a Healthcare FHIR (Fast Healthcare Interoperability Resources) Service

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_healthcare_workspace" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_healthcare_fhir_service" "example" {
  name                = "tfexfhir"
  location            = "east us"
  resource_group_name = "tfex-resource_group"
  workspace_id        = azurerm_healthcare_workspace.example.id
  kind                = "fhir-R4"

  authentication {
    authority = "https://login.microsoftonline.com/tenantId"
    audience  = "https://tfexfhir.fhir.azurehealthcareapis.com"
  }

  access_policy_object_ids = [
    data.azurerm_client_config.current.object_id
  ]

  identity {
    type = "SystemAssigned"
  }

  container_registry_login_server_url = ["tfex-container_registry_login_server"]

  cors {
    allowed_origins     = ["https://tfex.com:123", "https://tfex1.com:3389"]
    allowed_headers     = ["*"]
    allowed_methods     = ["GET", "DELETE", "PUT"]
    max_age_in_seconds  = 3600
    credentials_allowed = true
  }

  configuration_export_storage_account_name = "storage_account_name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Healthcare FHIR Service. Changing this forces a new Healthcare FHIR Service to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which to create the Healthcare FHIR Service. Changing this forces a new resource to be created.

* `workspace_id` - (Required) Specifies the id of the Healthcare Workspace where the Healthcare FHIR Service should exist. Changing this forces a new Healthcare FHIR Service to be created.

* `location` - (Required) Specifies the Azure Region where the Healthcare FHIR Service should be created. Changing this forces a new Healthcare FHIR Service to be created.

* `kind` - (Optional) Specifies the kind of the Healthcare FHIR Service. Possible values are: `fhir-Stu3` and `fhir-R4`. Defaults to `fhir-R4`. Changing this forces a new Healthcare FHIR Service to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `access_policy_object_ids` - (Optional) A list of the access policies of the service instance.

* `cors` - (Optional) A `cors` block as defined below.

* `container_registry_login_server_url` - (Optional) A list of azure container registry settings used for convert data operation of the service instance.

* `oci_artifact` - (Optional) [A list](/docs/configuration/attr-as-blocks.html) of `oci_artifact` objects as defined below to describe [OCI artifacts for export](https://learn.microsoft.com/en-gb/azure/healthcare-apis/fhir/de-identified-export).

* `authentication` - (Required) An `authentication` block as defined below.

* `configuration_export_storage_account_name` - (Optional) Specifies the name of the storage account which the operation configuration information is exported to.

* `tags` - (Optional) A mapping of tags to assign to the Healthcare FHIR Service.

---
An `identity` block supports the following:

* `type` - (Required) The type of managed identity to assign. Possible values are `UserAssigned` and `SystemAssigned`.

* `identity_ids` - (Optional) A list of one or more Resource IDs for User Assigned Managed identities to assign. Required when `type` is set to `UserAssigned`.

---
A `cors` block supports the following:

* `allowed_origins` - (Required) A set of origins to be allowed via CORS.
* `allowed_headers` - (Required) A set of headers to be allowed via CORS.
* `allowed_methods` - (Required) The methods to be allowed via CORS. Possible values are `DELETE`, `GET`, `HEAD`, `MERGE`, `POST`, `OPTIONS`, `PATCH` and `PUT`.
* `max_age_in_seconds` - (Optional) The max age to be allowed via CORS.
* `credentials_allowed` - (Optional) If credentials are allowed via CORS.

---
An `authentication` block supports the following:

* `authority` - (Required) The Azure Active Directory (tenant) that serves as the authentication authority to access the service. The default authority is the Directory defined in the authentication scheme in use when running Terraform.
  Authority must be registered to Azure AD and in the following format: <https://{Azure-AD-endpoint}/{tenant-id>}.
* `audience` - (Required) The intended audience to receive authentication tokens for the service.
* `smart_proxy_enabled` - (Optional) Whether smart proxy is enabled.

---

A `oci_artifact` block supports the following:

* `login_server` - (Required) An Azure container registry used for export operations of the service instance.

* `image_name` - (Optional) An image within Azure container registry used for export operations of the service instance.

* `digest` - (Optional) A digest of an image within Azure container registry used for export operations of the service instance to narrow the artifacts down.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Healthcare FHIR Service.

* `public_network_access_enabled` - Whether public networks access is enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Healthcare FHIR Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare FHIR Service.
* `update` - (Defaults to 90 minutes) Used when updating the Healthcare FHIR Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Healthcare FHIR Service.

## Import

Healthcare FHIR Service can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_fhir_service.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/fhirServices/service1
```
