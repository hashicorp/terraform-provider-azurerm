## Example: Encrypted Disks

This example provisions a Managed Disk which uses Customer Managed Encryption via a Disk Encryption Set.

## Notes

- Disk Encryption Sets are in Preview and only available in a limited set of regions.
- At this time Terraform's Azure Provider does not support enabling Soft Delete or Purge Protection; instead we use a `null_resource` to invoke the Azure CLI to add this behaviour. 
