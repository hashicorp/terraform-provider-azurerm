---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_dicom_service"
description: |-
  Manages a Healthcare DICOM (Digital Imaging and Communications in Medicine) Service.
---

# azurerm_healthcare_dicom_service

Manages a Healthcare DICOM Service

## Example Usage

```hcl
resource "azurerm_healthcare_workspace" "test" {
  name                = "tfexworkspace"
  resource_group_name = "tfex-resource_group"
  location            = "east us"
}

resource "azurerm_healthcare_dicom_service" "test" {
  name         = "tfexDicom"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = "east us"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "None"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Healthcare DICOM Service. Changing this forces a new Healthcare DICOM Service to be created.

* `workspace_id` - (Required) Specifies the ID of the Healthcare Workspace where the Healthcare DICOM Service should exist. Changing this forces a new Healthcare DICOM Service to be created.

* `location` - (Required) Specifies the Azure Region where the Healthcare DICOM Service should be created. Changing this forces a new Healthcare DICOM Service to be created.

* `data_partitions_enabled` - (Optional) If data partitions are enabled or not. Defaults to `false`. Changing this forces a new Healthcare DICOM Service to be created.

* `cors` - (Optional) A `cors` block as defined below.

* `encryption_key_url` - (Optional) The URL of the key to use for encryption as part of the customer-managed key encryption settings. For more details, refer to the [Azure Customer-Managed Keys Overview](https://learn.microsoft.com/en-us/azure/storage/common/customer-managed-keys-overview).

* `identity` - (Optional) An `identity` block as defined below.

* `public_network_access_enabled` - (Optional) Whether to enabled public networks when data plane traffic coming from public networks while private endpoint is enabled. Defaults to `true`.

* `storage` - (Optional) A `storage` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Healthcare DICOM Service.

---

An `identity` block supports the following:

* `type` - (Required) The type of identity used for the Healthcare DICOM service. Possible values are `UserAssigned`, `SystemAssigned` and `SystemAssigned, UserAssigned`. If `UserAssigned` is set, an `identity_ids` must be set as well.

* `identity_ids` - (Optional) A list of User Assigned Identity IDs which should be assigned to this Healthcare DICOM service.

---

A `cors` block supports the following:

* `allowed_origins` - (Optional) A list of allowed origins for CORS.

* `allowed_headers` - (Optional) A list of allowed headers for CORS.

* `allowed_methods` - (Optional) A list of allowed methods for CORS.

* `max_age_in_seconds` - (Optional) The maximum age in seconds for the CORS configuration (must be between 0 and 99998 inclusive).

* `allow_credentials` - (Optional) Whether to allow credentials in CORS. Defaults to `false`.

---

A `storage` block supports the following:

* `file_system_name` - (Required) The filesystem name of connected storage account. Changing this forces a new Healthcare DICOM Service to be created.

* `storage_account_id` - (Required) The resource ID of connected storage account. Changing this forces a new Healthcare DICOM Service to be created.

~> **Note:** The `is_hns_enabled` needs to be set to `true` for the storage account to be used with the Healthcare DICOM Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Healthcare DICOM Service.

* `authentication` - The `authentication` block as defined below.

* `service_url` - The url of the Healthcare DICOM Services.

---
An `authentication` block supports the following:

* `authority` - The Azure Active Directory (tenant) that serves as the authentication authority to access the service. The default authority is the Directory defined in the authentication scheme in use when running Terraform.
  Authority must be registered to Azure AD and in the following format: <https://{Azure-AD-endpoint}/{tenant-id>}.

* `audience` - The intended audience to receive authentication tokens for the service. The default value is <https://dicom.azurehealthcareapis.azure.com>

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Healthcare DICOM Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare DICOM Service.
* `update` - (Defaults to 90 minutes) Used when updating the Healthcare DICOM Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Healthcare DICOM Service.

## Import

Healthcare DICOM Service can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_dicom_service.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/dicomServices/service1
```
