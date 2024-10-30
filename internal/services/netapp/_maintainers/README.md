# This is a maintainer file for the Azure NetApp Files service.

Main idea of this file is to give insights to the revewiwers (and authors) on the details of the service and a few whys or do or dont's so PR review process gets faster and don't potentially introduces issues.

## Mainteiners
- paulomarquesc
- Orexii

## Acceptance Tests
- There is lack of SMB-related acct tests, the reason for that is that it requires Active Directory Domain Controller infrastucture which is not available, SMB-related tests can only be added the PR author provides such infrastructure which is not that trivial, so we should not require SMB tests until it is implemented unless it comes with Domain Controller setup automation, otherwise, these acceptance tests will be failing and causing disruptions in CI or bulk testing.
- New tests failing should not be accepted.

## Polling functions
- Some Netapp resources requires an extra type of polling mechanism, usually named wait`whatever is the resource name and operation` for example:

  ```go
    // Waiting for volume be completely provisioned
    if err := waitForVolumeCreateOrUpdate(ctx, client, id); err != nil {
        return err
    }
  ```	
  This seems redundant but it is not, some operations return from regular SDK polling as completed but in fact due to several factors it is still in progress, so we need to wait for it to be really completed, for this reason, we should not ask PR authors to remove these type of wait functions.