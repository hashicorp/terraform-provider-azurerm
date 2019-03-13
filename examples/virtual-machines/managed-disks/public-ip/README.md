## Example: Basic Virtual Machine with a Public IP Address using a Managed Disk

This example provisions a Virtual Machine with no Data Disks with a Managed Disk as the main OS Disk and a Public IP.

Notes:

- The files involved in this example are split out to make it easier to read, however all of the resources could be combined into a single file if needed.
- Your Public SSH Key will be uploaded to this instance so that you can SSH into it. This example assumes this is located at `~/.ssh/id_rsa.pub`.

### Variables

* `prefix` - (Required) The Prefix used for all resources in this example.
* `location` - (Required) The Azure Region in which the resources in this example should exist.
* `tags` - (Optional) Any tags which should be assigned to the resources in this example.
