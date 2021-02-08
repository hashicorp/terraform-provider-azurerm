---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_image_builder_template"
description: |-
  Manages an Image template that can be used to create virtual machine images.
---

# azurerm_image_builder_template

Manages an Image template that can be used to create virtual machine images.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "identityName"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_role_definition" "example" {
  name  = "roleDefName"
  scope = azurerm_resource_group.example.id

  permissions {
    actions = [
      "Microsoft.Compute/images/write",
      "Microsoft.Compute/images/read",
      "Microsoft.Compute/images/delete"
    ]
    not_actions = []
  }

  assignable_scopes = [
    azurerm_resource_group.example.id,
  ]
}

resource "azurerm_role_assignment" "example" {
  scope              = azurerm_resource_group.example.id
  role_definition_id = azurerm_role_definition.example.role_definition_resource_id
  principal_id       = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_image_builder_template" "example" {
  name                = "acctest"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  source_platform_image {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  distribution_managed_image {
    name                = "accTestImg1"
    resource_group_name = azurerm_resource_group.example.name
    location            = azurerm_resource_group.example.location
    run_output_name     = "ouputName"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Image Builder Template. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Image Builder Template should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Image Builder Template should exist. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below.

---

* `build_timeout_minutes` - (Optional) Maximum duration to wait while building the image template. Defaults to `240`. Changing this forces a new resource to be created.

* `customizer` - (Optional) A `customizer` block as defined below. Changing this forces a new resource to be created.

-> **Note** Azure Image Builder will run through the customizers in sequential order. Any failure in any customizer will fail the build process.

* `disk_size_gb` - (Optional) Size of the OS disk in GB. Defaults to 0 to leave the same size as the source image. You can only increase the size of the OS Disk (Win and Linux) but cannot reduce the OS Disk size to smaller than the size from the source image. Changing this forces a new resource to be created.

* `distribution_managed_image` - (Optional) A `distribution_managed_image` block as defined below. Changing this forces a new resource to be created.

* `distribution_shared_image` - (Optional) A `distribution_shared_image` block as defined below. Changing this forces a new resource to be created.

* `distribution_vhd` - (Optional) A `distribution_vhd` block as defined below. Changing this forces a new resource to be created.

-> **NOTE** At least one of `distribution_managed_image`, `distribution_shared_image` and `distribution_vhd` is required to specify.

* `size` - (Optional) Size of the virtual machine used to build, customize and capture images. Defaults to `Standard_D1_v2`. Changing this forces a new resource to be created.

* `source_managed_image_id` - (Optional) The ID of an image source that is a managed image. Changing this forces a new resource to be created.

* `source_platform_image` - (Optional) A `source_platform_image` block as defined below. Changing this forces a new resource to be created.

* `source_shared_image_version_id` - (Optional) The ID of an image source that is an image version in a shared image gallery. Changing this forces a new resource to be created.

-> **NOTE** Exactly one of `source_managed_image_id`, `source_platform_image` and `source_shared_image_version_id` is required to specify.

* `subnet_id` - (Optional) The ID of the subnet that the virtual machine used to build, customize and capture images. Changing this forces a new resource to be created.

-> **Note** If you specify a VNET, Azure Image Builder does not use a public IP address. Communication from Azure Image Builder Service to the build VM uses Azure Private Link technology. Private Link service requires an IP from the given VNET and subnet. Currently, Azure doesn’t support Network Policies on these IPs. Hence, network policies need to be disabled on the subnet, to do which, you need to ensure `enforce_private_link_service_network_policies` is set to `true` in your `azurerm_subnet` resource.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `customizer` block supports the following:

* `type` - (Required) The type of customization tool you want to use on the Image. Possible values are `File`, `PowerShell`, `Shell` `WindowsRestart` and `WindowsUpdate`. If you use `File` for example, you should only specify properties starting with `file_`, otherwise it will result into error

* `file_destination_path` - (Optional) The absolute path where the file will be downloaded to the VM. This must be specified when the customizer type is `File`. Changing this forces a new resource to be created.

* `file_sha256_checksum` - (Optional) SHA256 checksum of the file provided in the `file_source_uri` field below. Changing this forces a new resource to be created.

* `file_source_uri` - (Optional) The URI of the file to be uploaded for customizing the VM. This must be specified when the customizer type is `File`. Changing this forces a new resource to be created.

* `name` - (Optional) The name of the customizer. Changing this forces a new resource to be created.

* `powershell_commands` - (Optional) List of PowerShell commands to execute. Changing this forces a new resource to be created.

* `powershell_run_as_system` - (Optional) Whether the PowerShell script will be run with elevated privileges using the local system user? Changing this forces a new resource to be created.

* `powershell_run_elevated` - (Optional) Whether the PowerShell script will be run with elevated privileges? Changing this forces a new resource to be created.

-> **Note** `powershell_run_as_system` can only be true when `powershell_run_elevated` is set to true when the customizer type is `PowerShell`.

* `powershell_script_uri` - (Optional) URI of the PowerShell script to be run for customizing. Changing this forces a new resource to be created.

-> **Note** Exactly one of `powershell_commands` and `powershell_script_uri` must be specified when the customizer type is `PowerShell`.

* `powershell_sha256_checksum` - (Optional) SHA256 checksum of the power shell script provided in the `powershell_script_uri` above. Changing this forces a new resource to be created.

* `powershell_valid_exit_codes` - (Optional) List of valid exit codes that can be returned from commands/scripts. Changing this forces a new resource to be created.

* `shell_commands` - (Optional) List of shell commands to execute. Changing this forces a new resource to be created.

* `shell_script_uri` - (Optional) URI of the shell script to be run for customizing. Changing this forces a new resource to be created.

-> **Note** Exactly one of `shell_commands` and `shell_script_uri` must be specified when the customizer type is `Shell`.

* `shell_sha256_checksum` - (Optional) SHA256 checksum of the shell script provided in the `shell_script_uri` above. Changing this forces a new resource to be created.

* `windows_restart_check_command` - (Optional) Command to check if restart succeeded. Changing this forces a new resource to be created.

* `windows_restart_command` - (Optional) Command to execute the restart. Changing this forces a new resource to be created.

* `windows_restart_timeout` - (Optional) Restart timeout specified as a string of magnitude and unit, e.g. '5m' (5 minutes) or '2h' (2 hours). Changing this forces a new resource to be created.

* `windows_update_filters` - (Optional) List of filters to select updates to apply. Changing this forces a new resource to be created.

* `windows_update_search_criteria` - (Optional) Criteria to search updates. More to refer to https://github.com/rgl/packer-provisioner-windows-update. Changing this forces a new resource to be created.

* `windows_update_limit` - (Optional) Maximum number of updates to apply at a time. Changing this forces a new resource to be created.

---

A `distribution_managed_image` block supports the following:

* `name` - (Required) The name of the to be generated image. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the to be generated Image will exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the to be generated image will exist. Changing this forces a new resource to be created.

* `run_output_name` - (Required) The name to be used for the associated RunOutput. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the to be generated image. Changing this forces a new resource to be created.

---

A `distribution_shared_image` block supports the following:

* `id` - (Required) ID of the shared image gallery image. Changing this forces a new resource to be created.

* `replica_regions` - (Required) A `replica_regions` block as defined below. Changing this forces a new resource to be created.

* `run_output_name` - (Required) The name to be used for the associated RunOutput. Changing this forces a new resource to be created.

* `exclude_from_latest` - (Optional) Should the to be created Image Version be excluded from the `latest` filter? If set to `true` this Image Version won't be returned for the `latest` version. Defaults to `false`. Changing this forces a new resource to be created.

* `storage_account_type` - (Optional) The storage account type to store the to be created Image Version. Possible values are `Standard_LRS` and `Standard_ZRS`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the to be generated shared image version. Changing this forces a new resource to be created.

---

A `distribution_vhd` block supports the following:

* `run_output_name` - (Required) The name to be used for the associated RunOutput. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the to be generated vhd. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Identity which should be assigned to the Image Builder Template. Currently only `UserAssigned` is supported. Changing this forces a new resource to be created.

* `identity_ids` - (Required) The list of user assigned identities associated with the image template. Currently only 1 id is allowed. Changing this forces a new resource to be created.

---

A `plan` block supports the following:

* `name` - (Required) The name of the Marketplace Image this Image Builder Template should be created from. Changing this forces a new resource to be created.

* `product` - (Required) The product of the Marketplace Image this Image Builder Template should be created from. Changing this forces a new resource to be created.

* `publisher` - (Required) The publisher of the Marketplace Image this Image Builder Template should be created from. Changing this forces a new resource to be created.

---

A `source_platform_image` block supports the following:

* `publisher` - (Required) The publisher of the image used to create the Image Builder Template. Changing this forces a new resource to be created.

* `offer` - (Required) The offer of the image used to create the Image Builder Template. Changing this forces a new resource to be created.

* `sku` - (Required) The sku of the image used to create the Image Builder Template. Changing this forces a new resource to be created.

* `version` - (Required) The version of the image used to create the Image Builder Template. Changing this forces a new resource to be created.

* `plan` - (Optional) A `plan` block as defined above. Changing this forces a new resource to be created.

---

A `replica_regions` block supports the following:

* `name` - (Required) The Azure Region in which the to be created Shared Image Version should exist. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Image Builder Template.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Image Builder Template.
* `read` - (Defaults to 5 minutes) Used when retrieving the Image Builder Template.
* `update` - (Defaults to 30 minutes) Used when updating the Image Builder Template.
* `delete` - (Defaults to 30 minutes) Used when deleting the Image Builder Template.

## Import

Image Builder Templates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_image_builder_template.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.VirtualMachineImages/imageTemplates/imagebuildertemplate1
```
