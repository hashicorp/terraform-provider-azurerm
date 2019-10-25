## Example: Virtual Machine with Managed Disks from a Custom Image

This example provisions a Virtual Machine with Managed Disks from a Custom Image that already exists.

Notes:

- The files involved in this example are split out to make it easier to read, however all of the resources could be combined into a single file if needed.
- This example assumes the Custom Image specified exists - if it doesn't this example will fail.

### Variables

* `prefix` - (Required) The Prefix used for all resources in this example.
* `location` - (Required) The Azure Region in which the resources in this example should exist.
* `tags` - (Optional) Any tags which should be assigned to the resources in this example.

* `custom_image_resource_group_name` - (Required) The name of the Resource Group in which the Custom Image exists.
* `custom_image_name` - (Required) The name of the Custom Image to provision this Virtual Machine from.
