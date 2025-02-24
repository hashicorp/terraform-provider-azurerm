# This is a maintainer file for the Azure NetApp Files service.

This document gives insights into who is maintaining this service and includes details for why certain decisions in this service package have been made.

## Mainteiners

- paulomarquesc
- orexii

## Acceptance Tests

- There is lack of SMB-related acceptance tests because it requires Active Directory Domain Controller infrastructure which is not easy to automate properly. SMB-related tests can only be tested if the infrastructure is setup beforehand which is not that trivial. We should not require SMB tests unless it comes with Domain Controller setup automation. Without automation, the SMB acceptance tests will fail and cause disruptions in CI/bulk testing.

- New tests failing should not be accepted.

- For Azure NetApp Files, some features are highly dependent on specific regions, that's why for some acceptance tests, we will see regions defined there instead of for example `data.Locations.Primary` within the templates, this is expected and should not be changed. 

## Polling functions

- Some Netapp resources requires an extra type of polling mechanism. For example:

```go
// Waiting for volume be completely provisioned
if err := waitForVolumeCreateOrUpdate(ctx, client, id); err != nil {
  return err
}
```

  This is because some operations return from regular SDK polling as completed but due to several factors it is still in progress (e.g. ARM caching, software and hardware layer sync delays, etc.). These wait functions are necessary and should not be removed.

- Do not approve Pull Requests that relies on `<Create or Delete Operations>ThenPoll()` methods, e.g. `DeleteThenPoll()`, we should not use those for volume related operations due to some unknown [issues](https://github.com/hashicorp/pandora/issues/4571) with Pandora, those for Azure NetApp Files are not reliable, causing errors from time to time (and depending on the operation, very frequently) like this:

```text
pmarques [ ~/go/src/github.com/hashicorp/terraform-provider-azurerm ]$ make acctests SERVICE='netapp' TESTARGS=' -parallel 5 -run=TestAccNetAppVolumeGroupSAPHana_crossRegionReplication -count=1' TESTTIMEOUT='1200m'
==> Checking that code complies with gofmt requirements...
==> Checking that Custom Timeouts are used...
egrep: warning: egrep is obsolescent; using grep -E
egrep: warning: egrep is obsolescent; using grep -E
==> Checking that acceptance test packages are used...
TF_ACC=1 go test -v ./internal/services/netapp -parallel 5 -run=TestAccNetAppVolumeGroupSAPHana_crossRegionReplication -count=1 -timeout 1200m -ldflags="-X=github.com/hashicorp/terraform-provider-azurerm/version.ProviderVersion=acc"
=== RUN   TestAccNetAppVolumeGroupSAPHana_crossRegionReplication
=== PAUSE TestAccNetAppVolumeGroupSAPHana_crossRegionReplication
=== CONT  TestAccNetAppVolumeGroupSAPHana_crossRegionReplication
    testcase.go:173: Error running post-test destroy, there may be dangling resources: exit status 1

        Error: deleting `volume`: deleting replicate Volume (Subscription: "66bc9830-19b6-4987-94d2-0e487be7aa47"
        Resource Group Name: "acctestRG-netapp-241202215210177839"
        Net App Account Name: "acctest-NetAppAccount-Secondary-241202215210177839"
        Capacity Pool Name: "acctest-NetAppPool-Secondary-241202215210177839"
        Volume Name: "acctest-NetAppVolume-1-Secondary-241202215210177839"): polling after VolumesDeleteReplication: `result.Status` was nil/empty - `op.Status` was "DeleteReplication" / `op.Properties.ProvisioningState` was ""

        deleting `volume`: deleting replicate Volume (Subscription:
        "66bc9830-19b6-4987-94d2-0e487be7aa47"
        Resource Group Name: "acctestRG-netapp-241202215210177839"
        Net App Account Name: "acctest-NetAppAccount-Secondary-241202215210177839"
        Capacity Pool Name: "acctest-NetAppPool-Secondary-241202215210177839"
        Volume Name: "acctest-NetAppVolume-1-Secondary-241202215210177839"): polling
        after VolumesDeleteReplication: `result.Status` was nil/empty - `op.Status`
        was "DeleteReplication" / `op.Properties.ProvisioningState` was ""
--- FAIL: TestAccNetAppVolumeGroupSAPHana_crossRegionReplication (1375.67s)
FAIL
FAIL    github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp        1375.697s
FAIL
make: *** [GNUmakefile:103: acctests] Error 1
```

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

## Azure NetApp Files has features that requires disassociation, e.g. BackupPolicyId and SnapshotPolicyIds

- For cases where a property must have its content removed, mostly Ids (BackupPolicyId or SnapshotPolicyIds), instead of using `nil`, use `pointer.To("")`, this will trigger ANF RP to update the resource and set the value to empty string, setting as `nil` won't trigger any action within ANF RP.

E.g.

```golang
// Removing SnapshotId
update := volumes.VolumePatch{
    Properties: &volumes.VolumePatchProperties{
        DataProtection: &volumes.VolumePatchPropertiesDataProtection{
            Snapshot: &volumes.VolumeSnapshotProperties{
                SnapshotPolicyId: pointer.To(""),
            },
        },
    },
}
```

```golang
// Removing BackupPolicyId
backupPolicyIdRemoval := volumes.VolumePatch{
    Properties: &volumes.VolumePatchProperties{
        DataProtection: &volumes.VolumePatchPropertiesDataProtection{
            Backup: &volumes.VolumeBackupProperties{
                BackupPolicyId: pointer.To(""),
            },
        },
    },
}
```

## `Computed` attribute of a few configuration items

- Azure NetApp Files resources are complex enough on the backend and a few configuration items must remain in `computed = true` state and not changed under any circumstance, otherwise, if it is a `ForceNew = true` type of configuration, it may result in data loss. The best example of this is the configuration called `network_features`, this must be `computed = true` at all times, Azure NetApp Files team will soon make changes to networking that will make the property `network_features` be changed from `basic` to `standard` and to prevent data loss, this configuration should stay as is, with `computed = true` set.

```golang
"network_features": {
    Type:     pluginsdk.TypeString,
    Optional: true,
    Computed: true,
    Default:  string(volumes.NetworkFeaturesBasic),
    ValidateFunc: validation.StringInSlice([]string{
        string(volumes.NetworkFeaturesBasic),
        string(volumes.NetworkFeaturesStandard),
    }, false),
},
```

