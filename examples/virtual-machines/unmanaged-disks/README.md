## Example: Virtual Machine with Unmanaged Disks

This example provisions a Virtual Machine with Unmanaged Disks, with one OS Disk and one Data Disk.

Notes:

- The files involved in this example are split out to make it easier to read, however all of the resources could be combined into a single file if needed.

### Variables

* `prefix` - (Required) The Prefix used for all resources in this example.
* `location` - (Required) The Azure Region in which the resources in this example should exist.
* `tags` - (Optional) Any tags which should be assigned to the resources in this example.
