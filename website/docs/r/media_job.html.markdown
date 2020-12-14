---
subcategory: "Media"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_media_job"
description: |-
  Manages a Media Job.
---

# azurerm_media_job

Manages a Media Job.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "media-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "example" {
  name                = "examplemediaacc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  storage_account {
    id         = azurerm_storage_account.example.id
    is_primary = true
  }
}

resource "azurerm_media_transform" "example" {
  name                        = "transform1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  description                 = "My transform description"
  output {
    relative_priority = "Normal"
    on_error_action   = "ContinueJob"
    builtin_preset {
      preset_name = "AACGoodQualityAudio"
    }
  }
}

resource "azurerm_media_asset" "input" {
  name                        = "input"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  description                 = "Input Asset description"
}

resource "azurerm_media_asset" "output" {
  name                        = "output"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  description                 = "Output Asset description"
}

resource "azurerm_media_job" "example" {
  name                        = "job1"
  resource_group_name         = azurerm_resource_group.example.name
  media_services_account_name = azurerm_media_services_account.example.name
  transform_name              = azurerm_media_transform.example.name
  description                 = "My Job description"
  priority                    = "Normal"
  input_asset {
    asset_name = azurerm_media_asset.input.name
  }
  output_asset {
    asset_name = azurerm_media_asset.output.name
  }
}
```

## Arguments Reference

The following arguments are supported:

* `input_asset` - (Required) A `input_asset` block as defined below. Changing this forces a new Media Job to be created.

* `media_services_account_name` - (Required) The Media Services account name. Changing this forces a new Transform to be created.

* `name` - (Required) The name which should be used for this Media Job. Changing this forces a new Media Job to be created.

* `output_asset` - (Required) One or more `output_asset` blocks as defined below. Changing this forces a new Media Job to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Media Job should exist. Changing this forces a new Media Job to be created.

* `transform_name` - (Required) The Transform name. Changing this forces a new Media Job to be created.

---

* `description` - (Optional) 	Optional customer supplied description of the Job.

* `priority` - (Optional) Priority with which the job should be processed. Higher priority jobs are processed before lower priority jobs. Changing this forces a new Media Job to be created.

---

A `input_asset` block supports the following:

* `asset_name` - (Required) The name of the input Asset. Changing this forces a new Media Job to be created.

* `label` - (Optional) A label that is assigned to a JobInputClip, that is used to satisfy a reference used in the Transform. For example, a Transform can be authored so as to take an image file with the label 'xyz' and apply it as an overlay onto the input video before it is encoded. When submitting a Job, exactly one of the JobInputs should be the image file, and it should have the label 'xyz'.

---

A `output_asset` block supports the following:

* `asset_name` - (Required) The name of the output Asset. Changing this forces a new Media Job to be created.

* `label` - (Optional) A label that is assigned to a JobOutput in order to help uniquely identify it. This is useful when your Transform has more than one TransformOutput, whereby your Job has more than one JobOutput. In such cases, when you submit the Job, you will add two or more JobOutputs, in the same order as TransformOutputs in the Transform. Subsequently, when you retrieve the Job, either through events or on a GET request, you can use the label to easily identify the JobOutput. If a label is not provided, a default value of '{presetName}_{outputIndex}' will be used, where the preset name is the name of the preset in the corresponding TransformOutput and the output index is the relative index of the this JobOutput within the Job. Note that this index is the same as the relative index of the corresponding TransformOutput within its Transform.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Media Job.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Media Job.
* `read` - (Defaults to 5 minutes) Used when retrieving the Media Job.
* `update` - (Defaults to 30 minutes) Used when updating the Media Job.
* `delete` - (Defaults to 30 minutes) Used when deleting the Media Job.

## Import

Media Jobs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_media_job.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Media/mediaservices/account1/transforms/transform1/jobs/job1
```
