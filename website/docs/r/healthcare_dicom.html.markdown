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

* `workspace_id` - (Required) Specifies the id of the Healthcare Workspace where the Healthcare DICOM Service should exist. Changing this forces a new Healthcare DICOM Service to be created.

* `location` - (Required) Specifies the Azure Region where the Healthcare DICOM Service should be created. Changing this forces a new Healthcare DICOM Service to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `public_network_access_enabled` - (Optional) Whether to enabled public networks when data plane traffic coming from public networks while private endpoint is enabled. Defaults to `true`.

* `tags` - (Optional) A mapping of tags to assign to the Healthcare DICOM Service.

---

An `identity` block supports the following:

* `type` - (Required) The type of identity used for the Healthcare DICOM service. Possible values are `UserAssigned`, `SystemAssigned` and `SystemAssigned, UserAssigned`. If `UserAssigned` is set, an `identity_ids` must be set as well.

* `identity_ids` - (Optional) A list of User Assigned Identity IDs which should be assigned to this Healthcare DICOM service.

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
* `update` - (Defaults to 90 minutes) Used when updating the Healthcare DICOM Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare DICOM Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Healthcare DICOM Service.

## Import

Healthcare DICOM Service can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_dicom_service.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/dicomServices/service1
```
