# Change History

## Breaking Changes

### Removed Constants

1. ActionAfterReboot.ContinueConfiguration
1. ActionAfterReboot.StopConfiguration
1. ComplianceStatus.Compliant
1. ComplianceStatus.NonCompliant
1. ComplianceStatus.Pending
1. ConfigurationMode.ApplyAndAutoCorrect
1. ConfigurationMode.ApplyAndMonitor
1. ConfigurationMode.ApplyOnly
1. Kind.DSC
1. ProvisioningState.Canceled
1. ProvisioningState.Created
1. ProvisioningState.Failed
1. ProvisioningState.Succeeded
1. Type.Consistency
1. Type.Initial

### Removed Funcs

1. Navigation.MarshalJSON() ([]byte, error)

## Additive Changes

### New Constants

1. ActionAfterReboot.ActionAfterRebootContinueConfiguration
1. ActionAfterReboot.ActionAfterRebootStopConfiguration
1. AssignmentType.AssignmentTypeApplyAndAutoCorrect
1. AssignmentType.AssignmentTypeApplyAndMonitor
1. AssignmentType.AssignmentTypeAudit
1. AssignmentType.AssignmentTypeDeployAndAutoCorrect
1. ComplianceStatus.ComplianceStatusCompliant
1. ComplianceStatus.ComplianceStatusNonCompliant
1. ComplianceStatus.ComplianceStatusPending
1. ConfigurationMode.ConfigurationModeApplyAndAutoCorrect
1. ConfigurationMode.ConfigurationModeApplyAndMonitor
1. ConfigurationMode.ConfigurationModeApplyOnly
1. Kind.KindDSC
1. ProvisioningState.ProvisioningStateCanceled
1. ProvisioningState.ProvisioningStateCreated
1. ProvisioningState.ProvisioningStateFailed
1. ProvisioningState.ProvisioningStateSucceeded
1. Type.TypeConsistency
1. Type.TypeInitial

### New Funcs

1. PossibleAssignmentTypeValues() []AssignmentType

### Struct Changes

#### New Struct Fields

1. Navigation.AssignmentType
