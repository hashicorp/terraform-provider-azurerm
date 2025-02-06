# This is a maintainer file for the Azure NetApp Files service.

This document gives insights into who is maintaining this service and includes details for why certain decisions in this service package have been made.

## Mainteiners

- paulomarquesc
- orexii

## Acceptance Tests

- There is lack of SMB-related acceptance tests because it requires Active Directory Domain Controller infrastructure which is not easily automatable. SMB-related tests can only be tested if the infrastructure is setup beforehand which is not that trivial. We should not require SMB tests unless it comes with Domain Controller setup automation. Without automation, the SMB acceptance tests will fail and cause disruptions in CI/bulk testing.

- New tests failing should not be accepted.

## Polling functions

- Some Netapp resources requires an extra type of polling mechanism. For example:

```go
// Waiting for volume be completely provisioned
if err := waitForVolumeCreateOrUpdate(ctx, client, id); err != nil {
  return err
}
```

  This is because some operations return from regular SDK polling as completed but due to several factors it is still in progress (e.g. ARM caching, software and hardware layer sync delays, etc.). These wait functions are necessary and should not be removed.

## Data loss prevention protection

- Due to possibility of a volume to be deleted due to configuration changes on config file or changes made outside of Terraform, we have decided to not allow deletion of volumes by default. This is to prevent data loss. If you want to delete a volume, you need to set the feature block configuration `prevent_deletion_if_contains_resources` argument to `true`.

  In our tests we have added this block so the tests can delete the resources, but in production customers should not allow deletion of volumes by default.

```hcl
features {
  resource_group {
    prevent_deletion_if_contains_resources = false
  }
}
```

## Resource Group Deletion on Tests

- A feature block configuration needs to be set to ensure tests conclude when Microsoft maintainers execute their testing. Due to some Microsoft subscription management policies, certain resources are created by default during the tests. This block allows the deletion of the resource group at the end of the tests.

```hcl
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
```
