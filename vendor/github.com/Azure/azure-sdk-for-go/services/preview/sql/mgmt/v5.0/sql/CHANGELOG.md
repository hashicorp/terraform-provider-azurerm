# Change History

## Breaking Changes

### Removed Constants

1. BackupStorageRedundancy1.BackupStorageRedundancy1Geo
1. BackupStorageRedundancy1.BackupStorageRedundancy1Local
1. BackupStorageRedundancy1.BackupStorageRedundancy1Zone
1. CurrentBackupStorageRedundancy.CurrentBackupStorageRedundancyGeo
1. CurrentBackupStorageRedundancy.CurrentBackupStorageRedundancyLocal
1. CurrentBackupStorageRedundancy.CurrentBackupStorageRedundancyZone
1. RequestedBackupStorageRedundancy.RequestedBackupStorageRedundancyGeo
1. RequestedBackupStorageRedundancy.RequestedBackupStorageRedundancyLocal
1. RequestedBackupStorageRedundancy.RequestedBackupStorageRedundancyZone
1. Status1.Status1Failed
1. Status1.Status1Succeeded
1. StorageAccountType1.StorageAccountType1GRS
1. StorageAccountType1.StorageAccountType1LRS
1. StorageAccountType1.StorageAccountType1ZRS
1. TargetBackupStorageRedundancy.TargetBackupStorageRedundancyGeo
1. TargetBackupStorageRedundancy.TargetBackupStorageRedundancyLocal
1. TargetBackupStorageRedundancy.TargetBackupStorageRedundancyZone
1. TransparentDataEncryptionActivityStatus.TransparentDataEncryptionActivityStatusDecrypting
1. TransparentDataEncryptionActivityStatus.TransparentDataEncryptionActivityStatusEncrypting
1. TransparentDataEncryptionStatus.TransparentDataEncryptionStatusDisabled
1. TransparentDataEncryptionStatus.TransparentDataEncryptionStatusEnabled

### Removed Type Aliases

1. string.BackupStorageRedundancy1
1. string.CurrentBackupStorageRedundancy
1. string.RequestedBackupStorageRedundancy
1. string.Status1
1. string.StorageAccountType1
1. string.TargetBackupStorageRedundancy
1. string.TransparentDataEncryptionActivityStatus
1. string.TransparentDataEncryptionStatus

### Removed Funcs

1. *OperationsHealth.UnmarshalJSON([]byte) error
1. *OperationsHealthListResultIterator.Next() error
1. *OperationsHealthListResultIterator.NextWithContext(context.Context) error
1. *OperationsHealthListResultPage.Next() error
1. *OperationsHealthListResultPage.NextWithContext(context.Context) error
1. *ReplicationLinksUnlinkFuture.UnmarshalJSON([]byte) error
1. *TransparentDataEncryption.UnmarshalJSON([]byte) error
1. *TransparentDataEncryptionActivity.UnmarshalJSON([]byte) error
1. *UpdateManagedInstanceDNSServersOperation.UnmarshalJSON([]byte) error
1. DNSRefreshConfigurationProperties.MarshalJSON() ([]byte, error)
1. NewOperationsHealthClient(string) OperationsHealthClient
1. NewOperationsHealthClientWithBaseURI(string, string) OperationsHealthClient
1. NewOperationsHealthListResultIterator(OperationsHealthListResultPage) OperationsHealthListResultIterator
1. NewOperationsHealthListResultPage(OperationsHealthListResult, func(context.Context, OperationsHealthListResult) (OperationsHealthListResult, error)) OperationsHealthListResultPage
1. NewTransparentDataEncryptionActivitiesClient(string) TransparentDataEncryptionActivitiesClient
1. NewTransparentDataEncryptionActivitiesClientWithBaseURI(string, string) TransparentDataEncryptionActivitiesClient
1. OperationsHealth.MarshalJSON() ([]byte, error)
1. OperationsHealthClient.ListByLocation(context.Context, string) (OperationsHealthListResultPage, error)
1. OperationsHealthClient.ListByLocationComplete(context.Context, string) (OperationsHealthListResultIterator, error)
1. OperationsHealthClient.ListByLocationPreparer(context.Context, string) (*http.Request, error)
1. OperationsHealthClient.ListByLocationResponder(*http.Response) (OperationsHealthListResult, error)
1. OperationsHealthClient.ListByLocationSender(*http.Request) (*http.Response, error)
1. OperationsHealthListResult.IsEmpty() bool
1. OperationsHealthListResult.MarshalJSON() ([]byte, error)
1. OperationsHealthListResultIterator.NotDone() bool
1. OperationsHealthListResultIterator.Response() OperationsHealthListResult
1. OperationsHealthListResultIterator.Value() OperationsHealth
1. OperationsHealthListResultPage.NotDone() bool
1. OperationsHealthListResultPage.Response() OperationsHealthListResult
1. OperationsHealthListResultPage.Values() []OperationsHealth
1. OperationsHealthProperties.MarshalJSON() ([]byte, error)
1. PossibleBackupStorageRedundancy1Values() []BackupStorageRedundancy1
1. PossibleCurrentBackupStorageRedundancyValues() []CurrentBackupStorageRedundancy
1. PossibleRequestedBackupStorageRedundancyValues() []RequestedBackupStorageRedundancy
1. PossibleStatus1Values() []Status1
1. PossibleStorageAccountType1Values() []StorageAccountType1
1. PossibleTargetBackupStorageRedundancyValues() []TargetBackupStorageRedundancy
1. PossibleTransparentDataEncryptionActivityStatusValues() []TransparentDataEncryptionActivityStatus
1. PossibleTransparentDataEncryptionStatusValues() []TransparentDataEncryptionStatus
1. ReplicationLinksClient.Unlink(context.Context, string, string, string, string, UnlinkParameters) (ReplicationLinksUnlinkFuture, error)
1. ReplicationLinksClient.UnlinkPreparer(context.Context, string, string, string, string, UnlinkParameters) (*http.Request, error)
1. ReplicationLinksClient.UnlinkResponder(*http.Response) (autorest.Response, error)
1. ReplicationLinksClient.UnlinkSender(*http.Request) (ReplicationLinksUnlinkFuture, error)
1. ResourceIdentityWithUserAssignedIdentities.MarshalJSON() ([]byte, error)
1. TransparentDataEncryption.MarshalJSON() ([]byte, error)
1. TransparentDataEncryptionActivitiesClient.ListByConfiguration(context.Context, string, string, string) (TransparentDataEncryptionActivityListResult, error)
1. TransparentDataEncryptionActivitiesClient.ListByConfigurationPreparer(context.Context, string, string, string) (*http.Request, error)
1. TransparentDataEncryptionActivitiesClient.ListByConfigurationResponder(*http.Response) (TransparentDataEncryptionActivityListResult, error)
1. TransparentDataEncryptionActivitiesClient.ListByConfigurationSender(*http.Request) (*http.Response, error)
1. TransparentDataEncryptionActivity.MarshalJSON() ([]byte, error)
1. TransparentDataEncryptionActivityProperties.MarshalJSON() ([]byte, error)
1. UpdateManagedInstanceDNSServersOperation.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. DNSRefreshConfigurationProperties
1. OperationsHealth
1. OperationsHealthClient
1. OperationsHealthListResult
1. OperationsHealthListResultIterator
1. OperationsHealthListResultPage
1. OperationsHealthProperties
1. ReplicationLinksUnlinkFuture
1. ResourceIdentityWithUserAssignedIdentities
1. TransparentDataEncryption
1. TransparentDataEncryptionActivitiesClient
1. TransparentDataEncryptionActivity
1. TransparentDataEncryptionActivityListResult
1. TransparentDataEncryptionActivityProperties
1. UnlinkParameters
1. UpdateManagedInstanceDNSServersOperation

#### Removed Struct Fields

1. DatabaseUpdate.*DatabaseProperties
1. ManagedInstanceProperties.StorageAccountType
1. RestorableDroppedDatabaseProperties.ElasticPoolID
1. TransparentDataEncryptionProperties.Status
1. VirtualClusterProperties.Family
1. VirtualClusterProperties.MaintenanceConfigurationID

### Signature Changes

#### Funcs

1. ElasticPoolsClient.ListByServer
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, string, *int64
1. ElasticPoolsClient.ListByServerComplete
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, string, *int64
1. ElasticPoolsClient.ListByServerPreparer
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, string, *int64
1. LedgerDigestUploadsClient.CreateOrUpdate
	- Returns
		- From: LedgerDigestUploads, error
		- To: LedgerDigestUploadsCreateOrUpdateFuture, error
1. LedgerDigestUploadsClient.CreateOrUpdateSender
	- Returns
		- From: *http.Response, error
		- To: LedgerDigestUploadsCreateOrUpdateFuture, error
1. LedgerDigestUploadsClient.Disable
	- Returns
		- From: LedgerDigestUploads, error
		- To: LedgerDigestUploadsDisableFuture, error
1. LedgerDigestUploadsClient.DisableSender
	- Returns
		- From: *http.Response, error
		- To: LedgerDigestUploadsDisableFuture, error
1. ReplicationLinksClient.Delete
	- Returns
		- From: autorest.Response, error
		- To: ReplicationLinksDeleteFuture, error
1. ReplicationLinksClient.DeleteSender
	- Returns
		- From: *http.Response, error
		- To: ReplicationLinksDeleteFuture, error
1. ReplicationLinksClient.FailoverAllowDataLossResponder
	- Returns
		- From: autorest.Response, error
		- To: ReplicationLink, error
1. ReplicationLinksClient.FailoverResponder
	- Returns
		- From: autorest.Response, error
		- To: ReplicationLink, error
1. ServerConnectionPoliciesClient.CreateOrUpdate
	- Returns
		- From: ServerConnectionPolicy, error
		- To: ServerConnectionPoliciesCreateOrUpdateFuture, error
1. ServerConnectionPoliciesClient.CreateOrUpdateSender
	- Returns
		- From: *http.Response, error
		- To: ServerConnectionPoliciesCreateOrUpdateFuture, error
1. SyncGroupsClient.ListLogs
	- Params
		- From: context.Context, string, string, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string, SyncGroupsType, string
1. SyncGroupsClient.ListLogsComplete
	- Params
		- From: context.Context, string, string, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string, SyncGroupsType, string
1. SyncGroupsClient.ListLogsPreparer
	- Params
		- From: context.Context, string, string, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string, SyncGroupsType, string
1. TransparentDataEncryptionsClient.CreateOrUpdate
	- Params
		- From: context.Context, string, string, string, TransparentDataEncryption
		- To: context.Context, string, string, string, LogicalDatabaseTransparentDataEncryption
	- Returns
		- From: TransparentDataEncryption, error
		- To: LogicalDatabaseTransparentDataEncryption, error
1. TransparentDataEncryptionsClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, string, string, TransparentDataEncryption
		- To: context.Context, string, string, string, LogicalDatabaseTransparentDataEncryption
1. TransparentDataEncryptionsClient.CreateOrUpdateResponder
	- Returns
		- From: TransparentDataEncryption, error
		- To: LogicalDatabaseTransparentDataEncryption, error
1. TransparentDataEncryptionsClient.Get
	- Returns
		- From: TransparentDataEncryption, error
		- To: LogicalDatabaseTransparentDataEncryption, error
1. TransparentDataEncryptionsClient.GetResponder
	- Returns
		- From: TransparentDataEncryption, error
		- To: LogicalDatabaseTransparentDataEncryption, error
1. VirtualClustersClient.UpdateDNSServers
	- Returns
		- From: UpdateManagedInstanceDNSServersOperation, error
		- To: VirtualClustersUpdateDNSServersFuture, error
1. VirtualClustersClient.UpdateDNSServersResponder
	- Returns
		- From: UpdateManagedInstanceDNSServersOperation, error
		- To: UpdateVirtualClusterDNSServersOperation, error
1. VirtualClustersClient.UpdateDNSServersSender
	- Returns
		- From: *http.Response, error
		- To: VirtualClustersUpdateDNSServersFuture, error

#### Struct Fields

1. CopyLongTermRetentionBackupParametersProperties.TargetBackupStorageRedundancy changed type from TargetBackupStorageRedundancy to BackupStorageRedundancy
1. DatabaseProperties.CurrentBackupStorageRedundancy changed type from CurrentBackupStorageRedundancy to BackupStorageRedundancy
1. DatabaseProperties.RequestedBackupStorageRedundancy changed type from RequestedBackupStorageRedundancy to BackupStorageRedundancy
1. ManagedDatabaseRestoreDetailsProperties.NumberOfFilesDetected changed type from *int64 to *int32
1. ManagedDatabaseRestoreDetailsProperties.PercentCompleted changed type from *float64 to *int32
1. ManagedDatabaseRestoreDetailsProperties.UnrestorableFiles changed type from *[]string to *[]ManagedDatabaseRestoreDetailsUnrestorableFileProperties
1. ReplicationLinksFailoverAllowDataLossFuture.Result changed type from func(ReplicationLinksClient) (autorest.Response, error) to func(ReplicationLinksClient) (ReplicationLink, error)
1. ReplicationLinksFailoverFuture.Result changed type from func(ReplicationLinksClient) (autorest.Response, error) to func(ReplicationLinksClient) (ReplicationLink, error)
1. RestorableDroppedDatabaseProperties.BackupStorageRedundancy changed type from BackupStorageRedundancy1 to BackupStorageRedundancy
1. StorageCapability.StorageAccountType changed type from StorageAccountType1 to StorageAccountType
1. UpdateLongTermRetentionBackupParametersProperties.RequestedBackupStorageRedundancy changed type from RequestedBackupStorageRedundancy to BackupStorageRedundancy

## Additive Changes

### New Constants

1. AdvancedThreatProtectionState.AdvancedThreatProtectionStateDisabled
1. AdvancedThreatProtectionState.AdvancedThreatProtectionStateEnabled
1. AdvancedThreatProtectionState.AdvancedThreatProtectionStateNew
1. BackupStorageRedundancy.BackupStorageRedundancyGeoZone
1. DNSRefreshOperationStatus.DNSRefreshOperationStatusFailed
1. DNSRefreshOperationStatus.DNSRefreshOperationStatusInProgress
1. DNSRefreshOperationStatus.DNSRefreshOperationStatusSucceeded
1. DatabaseIdentityType.DatabaseIdentityTypeNone
1. DatabaseIdentityType.DatabaseIdentityTypeUserAssigned
1. DatabaseStatus.DatabaseStatusStarting
1. DatabaseStatus.DatabaseStatusStopped
1. DatabaseStatus.DatabaseStatusStopping
1. IdentityType.IdentityTypeSystemAssignedUserAssigned
1. MoveOperationMode.MoveOperationModeCopy
1. MoveOperationMode.MoveOperationModeMove
1. ProvisioningState1.ProvisioningState1Accepted
1. ProvisioningState1.ProvisioningState1Canceled
1. ProvisioningState1.ProvisioningState1Created
1. ProvisioningState1.ProvisioningState1Deleted
1. ProvisioningState1.ProvisioningState1NotSpecified
1. ProvisioningState1.ProvisioningState1Registering
1. ProvisioningState1.ProvisioningState1Running
1. ProvisioningState1.ProvisioningState1TimedOut
1. ProvisioningState1.ProvisioningState1Unrecognized
1. ReplicationLinkType.ReplicationLinkTypeSTANDBY
1. ReplicationMode.ReplicationModeAsync
1. ReplicationMode.ReplicationModeSync
1. RuleSeverity.RuleSeverityHigh
1. RuleSeverity.RuleSeverityInformational
1. RuleSeverity.RuleSeverityLow
1. RuleSeverity.RuleSeverityMedium
1. RuleSeverity.RuleSeverityObsolete
1. RuleStatus.RuleStatusFinding
1. RuleStatus.RuleStatusInternalError
1. RuleStatus.RuleStatusNonFinding
1. RuleType.RuleTypeBaselineExpected
1. RuleType.RuleTypeBinary
1. RuleType.RuleTypeNegativeList
1. RuleType.RuleTypePositiveList
1. SecondaryType.SecondaryTypeStandby
1. ServicePrincipalType.ServicePrincipalTypeNone
1. ServicePrincipalType.ServicePrincipalTypeSystemAssigned
1. SyncGroupsType.SyncGroupsTypeAll
1. SyncGroupsType.SyncGroupsTypeError
1. SyncGroupsType.SyncGroupsTypeSuccess
1. SyncGroupsType.SyncGroupsTypeWarning
1. VulnerabilityAssessmentState.VulnerabilityAssessmentStateDisabled
1. VulnerabilityAssessmentState.VulnerabilityAssessmentStateEnabled

### New Type Aliases

1. string.AdvancedThreatProtectionState
1. string.DNSRefreshOperationStatus
1. string.DatabaseIdentityType
1. string.MoveOperationMode
1. string.ReplicationMode
1. string.RuleSeverity
1. string.RuleStatus
1. string.RuleType
1. string.ServicePrincipalType
1. string.SyncGroupsType
1. string.VulnerabilityAssessmentState

### New Funcs

1. *DatabaseAdvancedThreatProtection.UnmarshalJSON([]byte) error
1. *DatabaseAdvancedThreatProtectionListResultIterator.Next() error
1. *DatabaseAdvancedThreatProtectionListResultIterator.NextWithContext(context.Context) error
1. *DatabaseAdvancedThreatProtectionListResultPage.Next() error
1. *DatabaseAdvancedThreatProtectionListResultPage.NextWithContext(context.Context) error
1. *DatabaseSQLVulnerabilityAssessmentBaselineSet.UnmarshalJSON([]byte) error
1. *DatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator.Next() error
1. *DatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator.NextWithContext(context.Context) error
1. *DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage.Next() error
1. *DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage.NextWithContext(context.Context) error
1. *DatabaseSQLVulnerabilityAssessmentExecuteScanExecuteFuture.UnmarshalJSON([]byte) error
1. *DatabaseSQLVulnerabilityAssessmentRuleBaseline.UnmarshalJSON([]byte) error
1. *DatabaseSQLVulnerabilityAssessmentRuleBaselineInput.UnmarshalJSON([]byte) error
1. *DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput.UnmarshalJSON([]byte) error
1. *DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator.Next() error
1. *DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator.NextWithContext(context.Context) error
1. *DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage.Next() error
1. *DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage.NextWithContext(context.Context) error
1. *DistributedAvailabilityGroup.UnmarshalJSON([]byte) error
1. *DistributedAvailabilityGroupsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *DistributedAvailabilityGroupsDeleteFuture.UnmarshalJSON([]byte) error
1. *DistributedAvailabilityGroupsListResultIterator.Next() error
1. *DistributedAvailabilityGroupsListResultIterator.NextWithContext(context.Context) error
1. *DistributedAvailabilityGroupsListResultPage.Next() error
1. *DistributedAvailabilityGroupsListResultPage.NextWithContext(context.Context) error
1. *DistributedAvailabilityGroupsUpdateFuture.UnmarshalJSON([]byte) error
1. *EndpointCertificate.UnmarshalJSON([]byte) error
1. *EndpointCertificateListResultIterator.Next() error
1. *EndpointCertificateListResultIterator.NextWithContext(context.Context) error
1. *EndpointCertificateListResultPage.Next() error
1. *EndpointCertificateListResultPage.NextWithContext(context.Context) error
1. *IPv6FirewallRule.UnmarshalJSON([]byte) error
1. *IPv6FirewallRuleListResultIterator.Next() error
1. *IPv6FirewallRuleListResultIterator.NextWithContext(context.Context) error
1. *IPv6FirewallRuleListResultPage.Next() error
1. *IPv6FirewallRuleListResultPage.NextWithContext(context.Context) error
1. *LedgerDigestUploadsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *LedgerDigestUploadsDisableFuture.UnmarshalJSON([]byte) error
1. *LogicalDatabaseTransparentDataEncryption.UnmarshalJSON([]byte) error
1. *LogicalDatabaseTransparentDataEncryptionListResultIterator.Next() error
1. *LogicalDatabaseTransparentDataEncryptionListResultIterator.NextWithContext(context.Context) error
1. *LogicalDatabaseTransparentDataEncryptionListResultPage.Next() error
1. *LogicalDatabaseTransparentDataEncryptionListResultPage.NextWithContext(context.Context) error
1. *LogicalServerAdvancedThreatProtectionListResultIterator.Next() error
1. *LogicalServerAdvancedThreatProtectionListResultIterator.NextWithContext(context.Context) error
1. *LogicalServerAdvancedThreatProtectionListResultPage.Next() error
1. *LogicalServerAdvancedThreatProtectionListResultPage.NextWithContext(context.Context) error
1. *ManagedDatabaseAdvancedThreatProtection.UnmarshalJSON([]byte) error
1. *ManagedDatabaseAdvancedThreatProtectionListResultIterator.Next() error
1. *ManagedDatabaseAdvancedThreatProtectionListResultIterator.NextWithContext(context.Context) error
1. *ManagedDatabaseAdvancedThreatProtectionListResultPage.Next() error
1. *ManagedDatabaseAdvancedThreatProtectionListResultPage.NextWithContext(context.Context) error
1. *ManagedDatabaseMoveOperationListResultIterator.Next() error
1. *ManagedDatabaseMoveOperationListResultIterator.NextWithContext(context.Context) error
1. *ManagedDatabaseMoveOperationListResultPage.Next() error
1. *ManagedDatabaseMoveOperationListResultPage.NextWithContext(context.Context) error
1. *ManagedDatabaseMoveOperationResult.UnmarshalJSON([]byte) error
1. *ManagedDatabasesCancelMoveFuture.UnmarshalJSON([]byte) error
1. *ManagedDatabasesCompleteMoveFuture.UnmarshalJSON([]byte) error
1. *ManagedDatabasesStartMoveFuture.UnmarshalJSON([]byte) error
1. *ManagedInstanceAdvancedThreatProtection.UnmarshalJSON([]byte) error
1. *ManagedInstanceAdvancedThreatProtectionListResultIterator.Next() error
1. *ManagedInstanceAdvancedThreatProtectionListResultIterator.NextWithContext(context.Context) error
1. *ManagedInstanceAdvancedThreatProtectionListResultPage.Next() error
1. *ManagedInstanceAdvancedThreatProtectionListResultPage.NextWithContext(context.Context) error
1. *ManagedInstanceAdvancedThreatProtectionSettingsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ManagedInstanceDtc.UnmarshalJSON([]byte) error
1. *ManagedInstanceDtcListResultIterator.Next() error
1. *ManagedInstanceDtcListResultIterator.NextWithContext(context.Context) error
1. *ManagedInstanceDtcListResultPage.Next() error
1. *ManagedInstanceDtcListResultPage.NextWithContext(context.Context) error
1. *ManagedInstanceDtcsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ManagedServerDNSAlias.UnmarshalJSON([]byte) error
1. *ManagedServerDNSAliasListResultIterator.Next() error
1. *ManagedServerDNSAliasListResultIterator.NextWithContext(context.Context) error
1. *ManagedServerDNSAliasListResultPage.Next() error
1. *ManagedServerDNSAliasListResultPage.NextWithContext(context.Context) error
1. *ManagedServerDNSAliasesAcquireFuture.UnmarshalJSON([]byte) error
1. *ManagedServerDNSAliasesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ManagedServerDNSAliasesDeleteFuture.UnmarshalJSON([]byte) error
1. *ReplicationLinksDeleteFuture.UnmarshalJSON([]byte) error
1. *ServerAdvancedThreatProtection.UnmarshalJSON([]byte) error
1. *ServerAdvancedThreatProtectionSettingsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ServerConnectionPoliciesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ServerConnectionPolicyListResultIterator.Next() error
1. *ServerConnectionPolicyListResultIterator.NextWithContext(context.Context) error
1. *ServerConnectionPolicyListResultPage.Next() error
1. *ServerConnectionPolicyListResultPage.NextWithContext(context.Context) error
1. *ServerTrustCertificate.UnmarshalJSON([]byte) error
1. *ServerTrustCertificatesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ServerTrustCertificatesDeleteFuture.UnmarshalJSON([]byte) error
1. *ServerTrustCertificatesListResultIterator.Next() error
1. *ServerTrustCertificatesListResultIterator.NextWithContext(context.Context) error
1. *ServerTrustCertificatesListResultPage.Next() error
1. *ServerTrustCertificatesListResultPage.NextWithContext(context.Context) error
1. *SynapseLinkWorkspace.UnmarshalJSON([]byte) error
1. *SynapseLinkWorkspaceListResultIterator.Next() error
1. *SynapseLinkWorkspaceListResultIterator.NextWithContext(context.Context) error
1. *SynapseLinkWorkspaceListResultPage.Next() error
1. *SynapseLinkWorkspaceListResultPage.NextWithContext(context.Context) error
1. *UpdateVirtualClusterDNSServersOperation.UnmarshalJSON([]byte) error
1. *VirtualClustersUpdateDNSServersFuture.UnmarshalJSON([]byte) error
1. *VulnerabilityAssessment.UnmarshalJSON([]byte) error
1. *VulnerabilityAssessmentExecuteScanExecuteFuture.UnmarshalJSON([]byte) error
1. *VulnerabilityAssessmentListResultIterator.Next() error
1. *VulnerabilityAssessmentListResultIterator.NextWithContext(context.Context) error
1. *VulnerabilityAssessmentListResultPage.Next() error
1. *VulnerabilityAssessmentListResultPage.NextWithContext(context.Context) error
1. *VulnerabilityAssessmentScanListResultIterator.Next() error
1. *VulnerabilityAssessmentScanListResultIterator.NextWithContext(context.Context) error
1. *VulnerabilityAssessmentScanListResultPage.Next() error
1. *VulnerabilityAssessmentScanListResultPage.NextWithContext(context.Context) error
1. *VulnerabilityAssessmentScanRecordListResultTypeIterator.Next() error
1. *VulnerabilityAssessmentScanRecordListResultTypeIterator.NextWithContext(context.Context) error
1. *VulnerabilityAssessmentScanRecordListResultTypePage.Next() error
1. *VulnerabilityAssessmentScanRecordListResultTypePage.NextWithContext(context.Context) error
1. *VulnerabilityAssessmentScanRecordType.UnmarshalJSON([]byte) error
1. *VulnerabilityAssessmentScanResults.UnmarshalJSON([]byte) error
1. AdvancedThreatProtectionProperties.MarshalJSON() ([]byte, error)
1. Baseline.MarshalJSON() ([]byte, error)
1. BaselineAdjustedResult.MarshalJSON() ([]byte, error)
1. BenchmarkReference.MarshalJSON() ([]byte, error)
1. DatabaseAdvancedThreatProtection.MarshalJSON() ([]byte, error)
1. DatabaseAdvancedThreatProtectionListResult.IsEmpty() bool
1. DatabaseAdvancedThreatProtectionListResult.MarshalJSON() ([]byte, error)
1. DatabaseAdvancedThreatProtectionListResultIterator.NotDone() bool
1. DatabaseAdvancedThreatProtectionListResultIterator.Response() DatabaseAdvancedThreatProtectionListResult
1. DatabaseAdvancedThreatProtectionListResultIterator.Value() DatabaseAdvancedThreatProtection
1. DatabaseAdvancedThreatProtectionListResultPage.NotDone() bool
1. DatabaseAdvancedThreatProtectionListResultPage.Response() DatabaseAdvancedThreatProtectionListResult
1. DatabaseAdvancedThreatProtectionListResultPage.Values() []DatabaseAdvancedThreatProtection
1. DatabaseAdvancedThreatProtectionSettingsClient.CreateOrUpdate(context.Context, string, string, string, DatabaseAdvancedThreatProtection) (DatabaseAdvancedThreatProtection, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.CreateOrUpdatePreparer(context.Context, string, string, string, DatabaseAdvancedThreatProtection) (*http.Request, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.CreateOrUpdateResponder(*http.Response) (DatabaseAdvancedThreatProtection, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.Get(context.Context, string, string, string) (DatabaseAdvancedThreatProtection, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.GetResponder(*http.Response) (DatabaseAdvancedThreatProtection, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.ListByDatabase(context.Context, string, string, string) (DatabaseAdvancedThreatProtectionListResultPage, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.ListByDatabaseComplete(context.Context, string, string, string) (DatabaseAdvancedThreatProtectionListResultIterator, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.ListByDatabasePreparer(context.Context, string, string, string) (*http.Request, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.ListByDatabaseResponder(*http.Response) (DatabaseAdvancedThreatProtectionListResult, error)
1. DatabaseAdvancedThreatProtectionSettingsClient.ListByDatabaseSender(*http.Request) (*http.Response, error)
1. DatabaseIdentity.MarshalJSON() ([]byte, error)
1. DatabaseSQLVulnerabilityAssessmentBaselineSet.MarshalJSON() ([]byte, error)
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResult.IsEmpty() bool
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResult.MarshalJSON() ([]byte, error)
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator.NotDone() bool
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator.Response() DatabaseSQLVulnerabilityAssessmentBaselineSetListResult
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator.Value() DatabaseSQLVulnerabilityAssessmentBaselineSet
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage.NotDone() bool
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage.Response() DatabaseSQLVulnerabilityAssessmentBaselineSetListResult
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage.Values() []DatabaseSQLVulnerabilityAssessmentBaselineSet
1. DatabaseSQLVulnerabilityAssessmentBaselineSetProperties.MarshalJSON() ([]byte, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.CreateOrUpdate(context.Context, string, string, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput) (DatabaseSQLVulnerabilityAssessmentBaselineSet, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.CreateOrUpdatePreparer(context.Context, string, string, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.CreateOrUpdateResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentBaselineSet, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.Get(context.Context, string, string, string) (DatabaseSQLVulnerabilityAssessmentBaselineSet, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.GetResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentBaselineSet, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.GetSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.ListBySQLVulnerabilityAssessment(context.Context, string, string, string) (DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.ListBySQLVulnerabilityAssessmentComplete(context.Context, string, string, string) (DatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.ListBySQLVulnerabilityAssessmentPreparer(context.Context, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.ListBySQLVulnerabilityAssessmentResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentBaselineSetListResult, error)
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient.ListBySQLVulnerabilityAssessmentSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentExecuteScanClient.Execute(context.Context, string, string, string) (DatabaseSQLVulnerabilityAssessmentExecuteScanExecuteFuture, error)
1. DatabaseSQLVulnerabilityAssessmentExecuteScanClient.ExecutePreparer(context.Context, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentExecuteScanClient.ExecuteResponder(*http.Response) (autorest.Response, error)
1. DatabaseSQLVulnerabilityAssessmentExecuteScanClient.ExecuteSender(*http.Request) (DatabaseSQLVulnerabilityAssessmentExecuteScanExecuteFuture, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaseline.MarshalJSON() ([]byte, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineInput.MarshalJSON() ([]byte, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput.MarshalJSON() ([]byte, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListInputProperties.MarshalJSON() ([]byte, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult.IsEmpty() bool
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult.MarshalJSON() ([]byte, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator.NotDone() bool
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator.Response() DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator.Value() DatabaseSQLVulnerabilityAssessmentRuleBaseline
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage.NotDone() bool
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage.Response() DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage.Values() []DatabaseSQLVulnerabilityAssessmentRuleBaseline
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.CreateOrUpdate(context.Context, string, string, string, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineInput) (DatabaseSQLVulnerabilityAssessmentRuleBaseline, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.CreateOrUpdatePreparer(context.Context, string, string, string, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineInput) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.CreateOrUpdateResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentRuleBaseline, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.Delete(context.Context, string, string, string, string) (autorest.Response, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.DeletePreparer(context.Context, string, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.DeleteSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.Get(context.Context, string, string, string, string) (DatabaseSQLVulnerabilityAssessmentRuleBaseline, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.GetResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentRuleBaseline, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.GetSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.ListByBaseline(context.Context, string, string, string) (DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.ListByBaselineComplete(context.Context, string, string, string) (DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.ListByBaselinePreparer(context.Context, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.ListByBaselineResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult, error)
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient.ListByBaselineSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentScanResultClient.Get(context.Context, string, string, string, string, string) (VulnerabilityAssessmentScanResults, error)
1. DatabaseSQLVulnerabilityAssessmentScanResultClient.GetPreparer(context.Context, string, string, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentScanResultClient.GetResponder(*http.Response) (VulnerabilityAssessmentScanResults, error)
1. DatabaseSQLVulnerabilityAssessmentScanResultClient.GetSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentScanResultClient.ListByScan(context.Context, string, string, string, string) (VulnerabilityAssessmentScanListResultPage, error)
1. DatabaseSQLVulnerabilityAssessmentScanResultClient.ListByScanComplete(context.Context, string, string, string, string) (VulnerabilityAssessmentScanListResultIterator, error)
1. DatabaseSQLVulnerabilityAssessmentScanResultClient.ListByScanPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentScanResultClient.ListByScanResponder(*http.Response) (VulnerabilityAssessmentScanListResult, error)
1. DatabaseSQLVulnerabilityAssessmentScanResultClient.ListByScanSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentScansClient.Get(context.Context, string, string, string, string) (VulnerabilityAssessmentScanRecordType, error)
1. DatabaseSQLVulnerabilityAssessmentScansClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentScansClient.GetResponder(*http.Response) (VulnerabilityAssessmentScanRecordType, error)
1. DatabaseSQLVulnerabilityAssessmentScansClient.GetSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessments(context.Context, string, string, string) (VulnerabilityAssessmentScanRecordListResultTypePage, error)
1. DatabaseSQLVulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessmentsComplete(context.Context, string, string, string) (VulnerabilityAssessmentScanRecordListResultTypeIterator, error)
1. DatabaseSQLVulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessmentsPreparer(context.Context, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessmentsResponder(*http.Response) (VulnerabilityAssessmentScanRecordListResultType, error)
1. DatabaseSQLVulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessmentsSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient.Get(context.Context, string, string, string) (VulnerabilityAssessment, error)
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient.GetResponder(*http.Response) (VulnerabilityAssessment, error)
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient.ListByDatabase(context.Context, string, string, string) (VulnerabilityAssessmentListResultPage, error)
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient.ListByDatabaseComplete(context.Context, string, string, string) (VulnerabilityAssessmentListResultIterator, error)
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient.ListByDatabasePreparer(context.Context, string, string, string) (*http.Request, error)
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient.ListByDatabaseResponder(*http.Response) (VulnerabilityAssessmentListResult, error)
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient.ListByDatabaseSender(*http.Request) (*http.Response, error)
1. DatabaseUpdateProperties.MarshalJSON() ([]byte, error)
1. DatabaseUserIdentity.MarshalJSON() ([]byte, error)
1. DistributedAvailabilityGroup.MarshalJSON() ([]byte, error)
1. DistributedAvailabilityGroupProperties.MarshalJSON() ([]byte, error)
1. DistributedAvailabilityGroupsClient.CreateOrUpdate(context.Context, string, string, string, DistributedAvailabilityGroup) (DistributedAvailabilityGroupsCreateOrUpdateFuture, error)
1. DistributedAvailabilityGroupsClient.CreateOrUpdatePreparer(context.Context, string, string, string, DistributedAvailabilityGroup) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.CreateOrUpdateResponder(*http.Response) (DistributedAvailabilityGroup, error)
1. DistributedAvailabilityGroupsClient.CreateOrUpdateSender(*http.Request) (DistributedAvailabilityGroupsCreateOrUpdateFuture, error)
1. DistributedAvailabilityGroupsClient.Delete(context.Context, string, string, string) (DistributedAvailabilityGroupsDeleteFuture, error)
1. DistributedAvailabilityGroupsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. DistributedAvailabilityGroupsClient.DeleteSender(*http.Request) (DistributedAvailabilityGroupsDeleteFuture, error)
1. DistributedAvailabilityGroupsClient.Get(context.Context, string, string, string) (DistributedAvailabilityGroup, error)
1. DistributedAvailabilityGroupsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.GetResponder(*http.Response) (DistributedAvailabilityGroup, error)
1. DistributedAvailabilityGroupsClient.GetSender(*http.Request) (*http.Response, error)
1. DistributedAvailabilityGroupsClient.ListByInstance(context.Context, string, string) (DistributedAvailabilityGroupsListResultPage, error)
1. DistributedAvailabilityGroupsClient.ListByInstanceComplete(context.Context, string, string) (DistributedAvailabilityGroupsListResultIterator, error)
1. DistributedAvailabilityGroupsClient.ListByInstancePreparer(context.Context, string, string) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.ListByInstanceResponder(*http.Response) (DistributedAvailabilityGroupsListResult, error)
1. DistributedAvailabilityGroupsClient.ListByInstanceSender(*http.Request) (*http.Response, error)
1. DistributedAvailabilityGroupsClient.Update(context.Context, string, string, string, DistributedAvailabilityGroup) (DistributedAvailabilityGroupsUpdateFuture, error)
1. DistributedAvailabilityGroupsClient.UpdatePreparer(context.Context, string, string, string, DistributedAvailabilityGroup) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.UpdateResponder(*http.Response) (DistributedAvailabilityGroup, error)
1. DistributedAvailabilityGroupsClient.UpdateSender(*http.Request) (DistributedAvailabilityGroupsUpdateFuture, error)
1. DistributedAvailabilityGroupsListResult.IsEmpty() bool
1. DistributedAvailabilityGroupsListResult.MarshalJSON() ([]byte, error)
1. DistributedAvailabilityGroupsListResultIterator.NotDone() bool
1. DistributedAvailabilityGroupsListResultIterator.Response() DistributedAvailabilityGroupsListResult
1. DistributedAvailabilityGroupsListResultIterator.Value() DistributedAvailabilityGroup
1. DistributedAvailabilityGroupsListResultPage.NotDone() bool
1. DistributedAvailabilityGroupsListResultPage.Response() DistributedAvailabilityGroupsListResult
1. DistributedAvailabilityGroupsListResultPage.Values() []DistributedAvailabilityGroup
1. EndpointCertificate.MarshalJSON() ([]byte, error)
1. EndpointCertificateListResult.IsEmpty() bool
1. EndpointCertificateListResult.MarshalJSON() ([]byte, error)
1. EndpointCertificateListResultIterator.NotDone() bool
1. EndpointCertificateListResultIterator.Response() EndpointCertificateListResult
1. EndpointCertificateListResultIterator.Value() EndpointCertificate
1. EndpointCertificateListResultPage.NotDone() bool
1. EndpointCertificateListResultPage.Response() EndpointCertificateListResult
1. EndpointCertificateListResultPage.Values() []EndpointCertificate
1. EndpointCertificatesClient.Get(context.Context, string, string, string) (EndpointCertificate, error)
1. EndpointCertificatesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. EndpointCertificatesClient.GetResponder(*http.Response) (EndpointCertificate, error)
1. EndpointCertificatesClient.GetSender(*http.Request) (*http.Response, error)
1. EndpointCertificatesClient.ListByInstance(context.Context, string, string) (EndpointCertificateListResultPage, error)
1. EndpointCertificatesClient.ListByInstanceComplete(context.Context, string, string) (EndpointCertificateListResultIterator, error)
1. EndpointCertificatesClient.ListByInstancePreparer(context.Context, string, string) (*http.Request, error)
1. EndpointCertificatesClient.ListByInstanceResponder(*http.Response) (EndpointCertificateListResult, error)
1. EndpointCertificatesClient.ListByInstanceSender(*http.Request) (*http.Response, error)
1. IPv6FirewallRule.MarshalJSON() ([]byte, error)
1. IPv6FirewallRuleListResult.IsEmpty() bool
1. IPv6FirewallRuleListResult.MarshalJSON() ([]byte, error)
1. IPv6FirewallRuleListResultIterator.NotDone() bool
1. IPv6FirewallRuleListResultIterator.Response() IPv6FirewallRuleListResult
1. IPv6FirewallRuleListResultIterator.Value() IPv6FirewallRule
1. IPv6FirewallRuleListResultPage.NotDone() bool
1. IPv6FirewallRuleListResultPage.Response() IPv6FirewallRuleListResult
1. IPv6FirewallRuleListResultPage.Values() []IPv6FirewallRule
1. IPv6FirewallRulesClient.CreateOrUpdate(context.Context, string, string, string, IPv6FirewallRule) (IPv6FirewallRule, error)
1. IPv6FirewallRulesClient.CreateOrUpdatePreparer(context.Context, string, string, string, IPv6FirewallRule) (*http.Request, error)
1. IPv6FirewallRulesClient.CreateOrUpdateResponder(*http.Response) (IPv6FirewallRule, error)
1. IPv6FirewallRulesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. IPv6FirewallRulesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. IPv6FirewallRulesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. IPv6FirewallRulesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. IPv6FirewallRulesClient.DeleteSender(*http.Request) (*http.Response, error)
1. IPv6FirewallRulesClient.Get(context.Context, string, string, string) (IPv6FirewallRule, error)
1. IPv6FirewallRulesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. IPv6FirewallRulesClient.GetResponder(*http.Response) (IPv6FirewallRule, error)
1. IPv6FirewallRulesClient.GetSender(*http.Request) (*http.Response, error)
1. IPv6FirewallRulesClient.ListByServer(context.Context, string, string) (IPv6FirewallRuleListResultPage, error)
1. IPv6FirewallRulesClient.ListByServerComplete(context.Context, string, string) (IPv6FirewallRuleListResultIterator, error)
1. IPv6FirewallRulesClient.ListByServerPreparer(context.Context, string, string) (*http.Request, error)
1. IPv6FirewallRulesClient.ListByServerResponder(*http.Response) (IPv6FirewallRuleListResult, error)
1. IPv6FirewallRulesClient.ListByServerSender(*http.Request) (*http.Response, error)
1. LogicalDatabaseTransparentDataEncryption.MarshalJSON() ([]byte, error)
1. LogicalDatabaseTransparentDataEncryptionListResult.IsEmpty() bool
1. LogicalDatabaseTransparentDataEncryptionListResult.MarshalJSON() ([]byte, error)
1. LogicalDatabaseTransparentDataEncryptionListResultIterator.NotDone() bool
1. LogicalDatabaseTransparentDataEncryptionListResultIterator.Response() LogicalDatabaseTransparentDataEncryptionListResult
1. LogicalDatabaseTransparentDataEncryptionListResultIterator.Value() LogicalDatabaseTransparentDataEncryption
1. LogicalDatabaseTransparentDataEncryptionListResultPage.NotDone() bool
1. LogicalDatabaseTransparentDataEncryptionListResultPage.Response() LogicalDatabaseTransparentDataEncryptionListResult
1. LogicalDatabaseTransparentDataEncryptionListResultPage.Values() []LogicalDatabaseTransparentDataEncryption
1. LogicalServerAdvancedThreatProtectionListResult.IsEmpty() bool
1. LogicalServerAdvancedThreatProtectionListResult.MarshalJSON() ([]byte, error)
1. LogicalServerAdvancedThreatProtectionListResultIterator.NotDone() bool
1. LogicalServerAdvancedThreatProtectionListResultIterator.Response() LogicalServerAdvancedThreatProtectionListResult
1. LogicalServerAdvancedThreatProtectionListResultIterator.Value() ServerAdvancedThreatProtection
1. LogicalServerAdvancedThreatProtectionListResultPage.NotDone() bool
1. LogicalServerAdvancedThreatProtectionListResultPage.Response() LogicalServerAdvancedThreatProtectionListResult
1. LogicalServerAdvancedThreatProtectionListResultPage.Values() []ServerAdvancedThreatProtection
1. ManagedDatabaseAdvancedThreatProtection.MarshalJSON() ([]byte, error)
1. ManagedDatabaseAdvancedThreatProtectionListResult.IsEmpty() bool
1. ManagedDatabaseAdvancedThreatProtectionListResult.MarshalJSON() ([]byte, error)
1. ManagedDatabaseAdvancedThreatProtectionListResultIterator.NotDone() bool
1. ManagedDatabaseAdvancedThreatProtectionListResultIterator.Response() ManagedDatabaseAdvancedThreatProtectionListResult
1. ManagedDatabaseAdvancedThreatProtectionListResultIterator.Value() ManagedDatabaseAdvancedThreatProtection
1. ManagedDatabaseAdvancedThreatProtectionListResultPage.NotDone() bool
1. ManagedDatabaseAdvancedThreatProtectionListResultPage.Response() ManagedDatabaseAdvancedThreatProtectionListResult
1. ManagedDatabaseAdvancedThreatProtectionListResultPage.Values() []ManagedDatabaseAdvancedThreatProtection
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.CreateOrUpdate(context.Context, string, string, string, ManagedDatabaseAdvancedThreatProtection) (ManagedDatabaseAdvancedThreatProtection, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.CreateOrUpdatePreparer(context.Context, string, string, string, ManagedDatabaseAdvancedThreatProtection) (*http.Request, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.CreateOrUpdateResponder(*http.Response) (ManagedDatabaseAdvancedThreatProtection, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.Get(context.Context, string, string, string) (ManagedDatabaseAdvancedThreatProtection, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.GetResponder(*http.Response) (ManagedDatabaseAdvancedThreatProtection, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.ListByDatabase(context.Context, string, string, string) (ManagedDatabaseAdvancedThreatProtectionListResultPage, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.ListByDatabaseComplete(context.Context, string, string, string) (ManagedDatabaseAdvancedThreatProtectionListResultIterator, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.ListByDatabasePreparer(context.Context, string, string, string) (*http.Request, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.ListByDatabaseResponder(*http.Response) (ManagedDatabaseAdvancedThreatProtectionListResult, error)
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient.ListByDatabaseSender(*http.Request) (*http.Response, error)
1. ManagedDatabaseMoveOperationListResult.IsEmpty() bool
1. ManagedDatabaseMoveOperationListResult.MarshalJSON() ([]byte, error)
1. ManagedDatabaseMoveOperationListResultIterator.NotDone() bool
1. ManagedDatabaseMoveOperationListResultIterator.Response() ManagedDatabaseMoveOperationListResult
1. ManagedDatabaseMoveOperationListResultIterator.Value() ManagedDatabaseMoveOperationResult
1. ManagedDatabaseMoveOperationListResultPage.NotDone() bool
1. ManagedDatabaseMoveOperationListResultPage.Response() ManagedDatabaseMoveOperationListResult
1. ManagedDatabaseMoveOperationListResultPage.Values() []ManagedDatabaseMoveOperationResult
1. ManagedDatabaseMoveOperationResult.MarshalJSON() ([]byte, error)
1. ManagedDatabaseMoveOperationResultProperties.MarshalJSON() ([]byte, error)
1. ManagedDatabaseMoveOperationsClient.Get(context.Context, string, string, uuid.UUID) (ManagedDatabaseMoveOperationResult, error)
1. ManagedDatabaseMoveOperationsClient.GetPreparer(context.Context, string, string, uuid.UUID) (*http.Request, error)
1. ManagedDatabaseMoveOperationsClient.GetResponder(*http.Response) (ManagedDatabaseMoveOperationResult, error)
1. ManagedDatabaseMoveOperationsClient.GetSender(*http.Request) (*http.Response, error)
1. ManagedDatabaseMoveOperationsClient.ListByLocation(context.Context, string, string, *bool, string) (ManagedDatabaseMoveOperationListResultPage, error)
1. ManagedDatabaseMoveOperationsClient.ListByLocationComplete(context.Context, string, string, *bool, string) (ManagedDatabaseMoveOperationListResultIterator, error)
1. ManagedDatabaseMoveOperationsClient.ListByLocationPreparer(context.Context, string, string, *bool, string) (*http.Request, error)
1. ManagedDatabaseMoveOperationsClient.ListByLocationResponder(*http.Response) (ManagedDatabaseMoveOperationListResult, error)
1. ManagedDatabaseMoveOperationsClient.ListByLocationSender(*http.Request) (*http.Response, error)
1. ManagedDatabaseRestoreDetailsBackupSetProperties.MarshalJSON() ([]byte, error)
1. ManagedDatabaseRestoreDetailsUnrestorableFileProperties.MarshalJSON() ([]byte, error)
1. ManagedDatabasesClient.CancelMove(context.Context, string, string, string, ManagedDatabaseMoveDefinition) (ManagedDatabasesCancelMoveFuture, error)
1. ManagedDatabasesClient.CancelMovePreparer(context.Context, string, string, string, ManagedDatabaseMoveDefinition) (*http.Request, error)
1. ManagedDatabasesClient.CancelMoveResponder(*http.Response) (autorest.Response, error)
1. ManagedDatabasesClient.CancelMoveSender(*http.Request) (ManagedDatabasesCancelMoveFuture, error)
1. ManagedDatabasesClient.CompleteMove(context.Context, string, string, string, ManagedDatabaseMoveDefinition) (ManagedDatabasesCompleteMoveFuture, error)
1. ManagedDatabasesClient.CompleteMovePreparer(context.Context, string, string, string, ManagedDatabaseMoveDefinition) (*http.Request, error)
1. ManagedDatabasesClient.CompleteMoveResponder(*http.Response) (autorest.Response, error)
1. ManagedDatabasesClient.CompleteMoveSender(*http.Request) (ManagedDatabasesCompleteMoveFuture, error)
1. ManagedDatabasesClient.StartMove(context.Context, string, string, string, ManagedDatabaseStartMoveDefinition) (ManagedDatabasesStartMoveFuture, error)
1. ManagedDatabasesClient.StartMovePreparer(context.Context, string, string, string, ManagedDatabaseStartMoveDefinition) (*http.Request, error)
1. ManagedDatabasesClient.StartMoveResponder(*http.Response) (autorest.Response, error)
1. ManagedDatabasesClient.StartMoveSender(*http.Request) (ManagedDatabasesStartMoveFuture, error)
1. ManagedInstanceAdvancedThreatProtection.MarshalJSON() ([]byte, error)
1. ManagedInstanceAdvancedThreatProtectionListResult.IsEmpty() bool
1. ManagedInstanceAdvancedThreatProtectionListResult.MarshalJSON() ([]byte, error)
1. ManagedInstanceAdvancedThreatProtectionListResultIterator.NotDone() bool
1. ManagedInstanceAdvancedThreatProtectionListResultIterator.Response() ManagedInstanceAdvancedThreatProtectionListResult
1. ManagedInstanceAdvancedThreatProtectionListResultIterator.Value() ManagedInstanceAdvancedThreatProtection
1. ManagedInstanceAdvancedThreatProtectionListResultPage.NotDone() bool
1. ManagedInstanceAdvancedThreatProtectionListResultPage.Response() ManagedInstanceAdvancedThreatProtectionListResult
1. ManagedInstanceAdvancedThreatProtectionListResultPage.Values() []ManagedInstanceAdvancedThreatProtection
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.CreateOrUpdate(context.Context, string, string, ManagedInstanceAdvancedThreatProtection) (ManagedInstanceAdvancedThreatProtectionSettingsCreateOrUpdateFuture, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.CreateOrUpdatePreparer(context.Context, string, string, ManagedInstanceAdvancedThreatProtection) (*http.Request, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.CreateOrUpdateResponder(*http.Response) (ManagedInstanceAdvancedThreatProtection, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.CreateOrUpdateSender(*http.Request) (ManagedInstanceAdvancedThreatProtectionSettingsCreateOrUpdateFuture, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.Get(context.Context, string, string) (ManagedInstanceAdvancedThreatProtection, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.GetResponder(*http.Response) (ManagedInstanceAdvancedThreatProtection, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.ListByInstance(context.Context, string, string) (ManagedInstanceAdvancedThreatProtectionListResultPage, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.ListByInstanceComplete(context.Context, string, string) (ManagedInstanceAdvancedThreatProtectionListResultIterator, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.ListByInstancePreparer(context.Context, string, string) (*http.Request, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.ListByInstanceResponder(*http.Response) (ManagedInstanceAdvancedThreatProtectionListResult, error)
1. ManagedInstanceAdvancedThreatProtectionSettingsClient.ListByInstanceSender(*http.Request) (*http.Response, error)
1. ManagedInstanceDtc.MarshalJSON() ([]byte, error)
1. ManagedInstanceDtcListResult.IsEmpty() bool
1. ManagedInstanceDtcListResult.MarshalJSON() ([]byte, error)
1. ManagedInstanceDtcListResultIterator.NotDone() bool
1. ManagedInstanceDtcListResultIterator.Response() ManagedInstanceDtcListResult
1. ManagedInstanceDtcListResultIterator.Value() ManagedInstanceDtc
1. ManagedInstanceDtcListResultPage.NotDone() bool
1. ManagedInstanceDtcListResultPage.Response() ManagedInstanceDtcListResult
1. ManagedInstanceDtcListResultPage.Values() []ManagedInstanceDtc
1. ManagedInstanceDtcProperties.MarshalJSON() ([]byte, error)
1. ManagedInstanceDtcsClient.CreateOrUpdate(context.Context, string, string, ManagedInstanceDtc) (ManagedInstanceDtcsCreateOrUpdateFuture, error)
1. ManagedInstanceDtcsClient.CreateOrUpdatePreparer(context.Context, string, string, ManagedInstanceDtc) (*http.Request, error)
1. ManagedInstanceDtcsClient.CreateOrUpdateResponder(*http.Response) (ManagedInstanceDtc, error)
1. ManagedInstanceDtcsClient.CreateOrUpdateSender(*http.Request) (ManagedInstanceDtcsCreateOrUpdateFuture, error)
1. ManagedInstanceDtcsClient.Get(context.Context, string, string) (ManagedInstanceDtc, error)
1. ManagedInstanceDtcsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ManagedInstanceDtcsClient.GetResponder(*http.Response) (ManagedInstanceDtc, error)
1. ManagedInstanceDtcsClient.GetSender(*http.Request) (*http.Response, error)
1. ManagedInstanceDtcsClient.ListByManagedInstance(context.Context, string, string) (ManagedInstanceDtcListResultPage, error)
1. ManagedInstanceDtcsClient.ListByManagedInstanceComplete(context.Context, string, string) (ManagedInstanceDtcListResultIterator, error)
1. ManagedInstanceDtcsClient.ListByManagedInstancePreparer(context.Context, string, string) (*http.Request, error)
1. ManagedInstanceDtcsClient.ListByManagedInstanceResponder(*http.Response) (ManagedInstanceDtcListResult, error)
1. ManagedInstanceDtcsClient.ListByManagedInstanceSender(*http.Request) (*http.Response, error)
1. ManagedServerDNSAlias.MarshalJSON() ([]byte, error)
1. ManagedServerDNSAliasListResult.IsEmpty() bool
1. ManagedServerDNSAliasListResult.MarshalJSON() ([]byte, error)
1. ManagedServerDNSAliasListResultIterator.NotDone() bool
1. ManagedServerDNSAliasListResultIterator.Response() ManagedServerDNSAliasListResult
1. ManagedServerDNSAliasListResultIterator.Value() ManagedServerDNSAlias
1. ManagedServerDNSAliasListResultPage.NotDone() bool
1. ManagedServerDNSAliasListResultPage.Response() ManagedServerDNSAliasListResult
1. ManagedServerDNSAliasListResultPage.Values() []ManagedServerDNSAlias
1. ManagedServerDNSAliasProperties.MarshalJSON() ([]byte, error)
1. ManagedServerDNSAliasesClient.Acquire(context.Context, string, string, string, ManagedServerDNSAliasAcquisition) (ManagedServerDNSAliasesAcquireFuture, error)
1. ManagedServerDNSAliasesClient.AcquirePreparer(context.Context, string, string, string, ManagedServerDNSAliasAcquisition) (*http.Request, error)
1. ManagedServerDNSAliasesClient.AcquireResponder(*http.Response) (ManagedServerDNSAlias, error)
1. ManagedServerDNSAliasesClient.AcquireSender(*http.Request) (ManagedServerDNSAliasesAcquireFuture, error)
1. ManagedServerDNSAliasesClient.CreateOrUpdate(context.Context, string, string, string, ManagedServerDNSAliasCreation) (ManagedServerDNSAliasesCreateOrUpdateFuture, error)
1. ManagedServerDNSAliasesClient.CreateOrUpdatePreparer(context.Context, string, string, string, ManagedServerDNSAliasCreation) (*http.Request, error)
1. ManagedServerDNSAliasesClient.CreateOrUpdateResponder(*http.Response) (ManagedServerDNSAlias, error)
1. ManagedServerDNSAliasesClient.CreateOrUpdateSender(*http.Request) (ManagedServerDNSAliasesCreateOrUpdateFuture, error)
1. ManagedServerDNSAliasesClient.Delete(context.Context, string, string, string) (ManagedServerDNSAliasesDeleteFuture, error)
1. ManagedServerDNSAliasesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ManagedServerDNSAliasesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ManagedServerDNSAliasesClient.DeleteSender(*http.Request) (ManagedServerDNSAliasesDeleteFuture, error)
1. ManagedServerDNSAliasesClient.Get(context.Context, string, string, string) (ManagedServerDNSAlias, error)
1. ManagedServerDNSAliasesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ManagedServerDNSAliasesClient.GetResponder(*http.Response) (ManagedServerDNSAlias, error)
1. ManagedServerDNSAliasesClient.GetSender(*http.Request) (*http.Response, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstance(context.Context, string, string) (ManagedServerDNSAliasListResultPage, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstanceComplete(context.Context, string, string) (ManagedServerDNSAliasListResultIterator, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstancePreparer(context.Context, string, string) (*http.Request, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstanceResponder(*http.Response) (ManagedServerDNSAliasListResult, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstanceSender(*http.Request) (*http.Response, error)
1. NewDatabaseAdvancedThreatProtectionListResultIterator(DatabaseAdvancedThreatProtectionListResultPage) DatabaseAdvancedThreatProtectionListResultIterator
1. NewDatabaseAdvancedThreatProtectionListResultPage(DatabaseAdvancedThreatProtectionListResult, func(context.Context, DatabaseAdvancedThreatProtectionListResult) (DatabaseAdvancedThreatProtectionListResult, error)) DatabaseAdvancedThreatProtectionListResultPage
1. NewDatabaseAdvancedThreatProtectionSettingsClient(string) DatabaseAdvancedThreatProtectionSettingsClient
1. NewDatabaseAdvancedThreatProtectionSettingsClientWithBaseURI(string, string) DatabaseAdvancedThreatProtectionSettingsClient
1. NewDatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator(DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage) DatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator
1. NewDatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage(DatabaseSQLVulnerabilityAssessmentBaselineSetListResult, func(context.Context, DatabaseSQLVulnerabilityAssessmentBaselineSetListResult) (DatabaseSQLVulnerabilityAssessmentBaselineSetListResult, error)) DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage
1. NewDatabaseSQLVulnerabilityAssessmentBaselinesClient(string) DatabaseSQLVulnerabilityAssessmentBaselinesClient
1. NewDatabaseSQLVulnerabilityAssessmentBaselinesClientWithBaseURI(string, string) DatabaseSQLVulnerabilityAssessmentBaselinesClient
1. NewDatabaseSQLVulnerabilityAssessmentExecuteScanClient(string) DatabaseSQLVulnerabilityAssessmentExecuteScanClient
1. NewDatabaseSQLVulnerabilityAssessmentExecuteScanClientWithBaseURI(string, string) DatabaseSQLVulnerabilityAssessmentExecuteScanClient
1. NewDatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator(DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage) DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator
1. NewDatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage(DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult, func(context.Context, DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult) (DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult, error)) DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage
1. NewDatabaseSQLVulnerabilityAssessmentRuleBaselinesClient(string) DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient
1. NewDatabaseSQLVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(string, string) DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient
1. NewDatabaseSQLVulnerabilityAssessmentScanResultClient(string) DatabaseSQLVulnerabilityAssessmentScanResultClient
1. NewDatabaseSQLVulnerabilityAssessmentScanResultClientWithBaseURI(string, string) DatabaseSQLVulnerabilityAssessmentScanResultClient
1. NewDatabaseSQLVulnerabilityAssessmentScansClient(string) DatabaseSQLVulnerabilityAssessmentScansClient
1. NewDatabaseSQLVulnerabilityAssessmentScansClientWithBaseURI(string, string) DatabaseSQLVulnerabilityAssessmentScansClient
1. NewDatabaseSQLVulnerabilityAssessmentsSettingsClient(string) DatabaseSQLVulnerabilityAssessmentsSettingsClient
1. NewDatabaseSQLVulnerabilityAssessmentsSettingsClientWithBaseURI(string, string) DatabaseSQLVulnerabilityAssessmentsSettingsClient
1. NewDistributedAvailabilityGroupsClient(string) DistributedAvailabilityGroupsClient
1. NewDistributedAvailabilityGroupsClientWithBaseURI(string, string) DistributedAvailabilityGroupsClient
1. NewDistributedAvailabilityGroupsListResultIterator(DistributedAvailabilityGroupsListResultPage) DistributedAvailabilityGroupsListResultIterator
1. NewDistributedAvailabilityGroupsListResultPage(DistributedAvailabilityGroupsListResult, func(context.Context, DistributedAvailabilityGroupsListResult) (DistributedAvailabilityGroupsListResult, error)) DistributedAvailabilityGroupsListResultPage
1. NewEndpointCertificateListResultIterator(EndpointCertificateListResultPage) EndpointCertificateListResultIterator
1. NewEndpointCertificateListResultPage(EndpointCertificateListResult, func(context.Context, EndpointCertificateListResult) (EndpointCertificateListResult, error)) EndpointCertificateListResultPage
1. NewEndpointCertificatesClient(string) EndpointCertificatesClient
1. NewEndpointCertificatesClientWithBaseURI(string, string) EndpointCertificatesClient
1. NewIPv6FirewallRuleListResultIterator(IPv6FirewallRuleListResultPage) IPv6FirewallRuleListResultIterator
1. NewIPv6FirewallRuleListResultPage(IPv6FirewallRuleListResult, func(context.Context, IPv6FirewallRuleListResult) (IPv6FirewallRuleListResult, error)) IPv6FirewallRuleListResultPage
1. NewIPv6FirewallRulesClient(string) IPv6FirewallRulesClient
1. NewIPv6FirewallRulesClientWithBaseURI(string, string) IPv6FirewallRulesClient
1. NewLogicalDatabaseTransparentDataEncryptionListResultIterator(LogicalDatabaseTransparentDataEncryptionListResultPage) LogicalDatabaseTransparentDataEncryptionListResultIterator
1. NewLogicalDatabaseTransparentDataEncryptionListResultPage(LogicalDatabaseTransparentDataEncryptionListResult, func(context.Context, LogicalDatabaseTransparentDataEncryptionListResult) (LogicalDatabaseTransparentDataEncryptionListResult, error)) LogicalDatabaseTransparentDataEncryptionListResultPage
1. NewLogicalServerAdvancedThreatProtectionListResultIterator(LogicalServerAdvancedThreatProtectionListResultPage) LogicalServerAdvancedThreatProtectionListResultIterator
1. NewLogicalServerAdvancedThreatProtectionListResultPage(LogicalServerAdvancedThreatProtectionListResult, func(context.Context, LogicalServerAdvancedThreatProtectionListResult) (LogicalServerAdvancedThreatProtectionListResult, error)) LogicalServerAdvancedThreatProtectionListResultPage
1. NewManagedDatabaseAdvancedThreatProtectionListResultIterator(ManagedDatabaseAdvancedThreatProtectionListResultPage) ManagedDatabaseAdvancedThreatProtectionListResultIterator
1. NewManagedDatabaseAdvancedThreatProtectionListResultPage(ManagedDatabaseAdvancedThreatProtectionListResult, func(context.Context, ManagedDatabaseAdvancedThreatProtectionListResult) (ManagedDatabaseAdvancedThreatProtectionListResult, error)) ManagedDatabaseAdvancedThreatProtectionListResultPage
1. NewManagedDatabaseAdvancedThreatProtectionSettingsClient(string) ManagedDatabaseAdvancedThreatProtectionSettingsClient
1. NewManagedDatabaseAdvancedThreatProtectionSettingsClientWithBaseURI(string, string) ManagedDatabaseAdvancedThreatProtectionSettingsClient
1. NewManagedDatabaseMoveOperationListResultIterator(ManagedDatabaseMoveOperationListResultPage) ManagedDatabaseMoveOperationListResultIterator
1. NewManagedDatabaseMoveOperationListResultPage(ManagedDatabaseMoveOperationListResult, func(context.Context, ManagedDatabaseMoveOperationListResult) (ManagedDatabaseMoveOperationListResult, error)) ManagedDatabaseMoveOperationListResultPage
1. NewManagedDatabaseMoveOperationsClient(string) ManagedDatabaseMoveOperationsClient
1. NewManagedDatabaseMoveOperationsClientWithBaseURI(string, string) ManagedDatabaseMoveOperationsClient
1. NewManagedInstanceAdvancedThreatProtectionListResultIterator(ManagedInstanceAdvancedThreatProtectionListResultPage) ManagedInstanceAdvancedThreatProtectionListResultIterator
1. NewManagedInstanceAdvancedThreatProtectionListResultPage(ManagedInstanceAdvancedThreatProtectionListResult, func(context.Context, ManagedInstanceAdvancedThreatProtectionListResult) (ManagedInstanceAdvancedThreatProtectionListResult, error)) ManagedInstanceAdvancedThreatProtectionListResultPage
1. NewManagedInstanceAdvancedThreatProtectionSettingsClient(string) ManagedInstanceAdvancedThreatProtectionSettingsClient
1. NewManagedInstanceAdvancedThreatProtectionSettingsClientWithBaseURI(string, string) ManagedInstanceAdvancedThreatProtectionSettingsClient
1. NewManagedInstanceDtcListResultIterator(ManagedInstanceDtcListResultPage) ManagedInstanceDtcListResultIterator
1. NewManagedInstanceDtcListResultPage(ManagedInstanceDtcListResult, func(context.Context, ManagedInstanceDtcListResult) (ManagedInstanceDtcListResult, error)) ManagedInstanceDtcListResultPage
1. NewManagedInstanceDtcsClient(string) ManagedInstanceDtcsClient
1. NewManagedInstanceDtcsClientWithBaseURI(string, string) ManagedInstanceDtcsClient
1. NewManagedServerDNSAliasListResultIterator(ManagedServerDNSAliasListResultPage) ManagedServerDNSAliasListResultIterator
1. NewManagedServerDNSAliasListResultPage(ManagedServerDNSAliasListResult, func(context.Context, ManagedServerDNSAliasListResult) (ManagedServerDNSAliasListResult, error)) ManagedServerDNSAliasListResultPage
1. NewManagedServerDNSAliasesClient(string) ManagedServerDNSAliasesClient
1. NewManagedServerDNSAliasesClientWithBaseURI(string, string) ManagedServerDNSAliasesClient
1. NewServerAdvancedThreatProtectionSettingsClient(string) ServerAdvancedThreatProtectionSettingsClient
1. NewServerAdvancedThreatProtectionSettingsClientWithBaseURI(string, string) ServerAdvancedThreatProtectionSettingsClient
1. NewServerConnectionPolicyListResultIterator(ServerConnectionPolicyListResultPage) ServerConnectionPolicyListResultIterator
1. NewServerConnectionPolicyListResultPage(ServerConnectionPolicyListResult, func(context.Context, ServerConnectionPolicyListResult) (ServerConnectionPolicyListResult, error)) ServerConnectionPolicyListResultPage
1. NewServerTrustCertificatesClient(string) ServerTrustCertificatesClient
1. NewServerTrustCertificatesClientWithBaseURI(string, string) ServerTrustCertificatesClient
1. NewServerTrustCertificatesListResultIterator(ServerTrustCertificatesListResultPage) ServerTrustCertificatesListResultIterator
1. NewServerTrustCertificatesListResultPage(ServerTrustCertificatesListResult, func(context.Context, ServerTrustCertificatesListResult) (ServerTrustCertificatesListResult, error)) ServerTrustCertificatesListResultPage
1. NewSynapseLinkWorkspaceListResultIterator(SynapseLinkWorkspaceListResultPage) SynapseLinkWorkspaceListResultIterator
1. NewSynapseLinkWorkspaceListResultPage(SynapseLinkWorkspaceListResult, func(context.Context, SynapseLinkWorkspaceListResult) (SynapseLinkWorkspaceListResult, error)) SynapseLinkWorkspaceListResultPage
1. NewSynapseLinkWorkspacesClient(string) SynapseLinkWorkspacesClient
1. NewSynapseLinkWorkspacesClientWithBaseURI(string, string) SynapseLinkWorkspacesClient
1. NewVulnerabilityAssessmentBaselineClient(string) VulnerabilityAssessmentBaselineClient
1. NewVulnerabilityAssessmentBaselineClientWithBaseURI(string, string) VulnerabilityAssessmentBaselineClient
1. NewVulnerabilityAssessmentBaselinesClient(string) VulnerabilityAssessmentBaselinesClient
1. NewVulnerabilityAssessmentBaselinesClientWithBaseURI(string, string) VulnerabilityAssessmentBaselinesClient
1. NewVulnerabilityAssessmentExecuteScanClient(string) VulnerabilityAssessmentExecuteScanClient
1. NewVulnerabilityAssessmentExecuteScanClientWithBaseURI(string, string) VulnerabilityAssessmentExecuteScanClient
1. NewVulnerabilityAssessmentListResultIterator(VulnerabilityAssessmentListResultPage) VulnerabilityAssessmentListResultIterator
1. NewVulnerabilityAssessmentListResultPage(VulnerabilityAssessmentListResult, func(context.Context, VulnerabilityAssessmentListResult) (VulnerabilityAssessmentListResult, error)) VulnerabilityAssessmentListResultPage
1. NewVulnerabilityAssessmentRuleBaselineClient(string) VulnerabilityAssessmentRuleBaselineClient
1. NewVulnerabilityAssessmentRuleBaselineClientWithBaseURI(string, string) VulnerabilityAssessmentRuleBaselineClient
1. NewVulnerabilityAssessmentRuleBaselinesClient(string) VulnerabilityAssessmentRuleBaselinesClient
1. NewVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(string, string) VulnerabilityAssessmentRuleBaselinesClient
1. NewVulnerabilityAssessmentScanListResultIterator(VulnerabilityAssessmentScanListResultPage) VulnerabilityAssessmentScanListResultIterator
1. NewVulnerabilityAssessmentScanListResultPage(VulnerabilityAssessmentScanListResult, func(context.Context, VulnerabilityAssessmentScanListResult) (VulnerabilityAssessmentScanListResult, error)) VulnerabilityAssessmentScanListResultPage
1. NewVulnerabilityAssessmentScanRecordListResultTypeIterator(VulnerabilityAssessmentScanRecordListResultTypePage) VulnerabilityAssessmentScanRecordListResultTypeIterator
1. NewVulnerabilityAssessmentScanRecordListResultTypePage(VulnerabilityAssessmentScanRecordListResultType, func(context.Context, VulnerabilityAssessmentScanRecordListResultType) (VulnerabilityAssessmentScanRecordListResultType, error)) VulnerabilityAssessmentScanRecordListResultTypePage
1. NewVulnerabilityAssessmentScanResultClient(string) VulnerabilityAssessmentScanResultClient
1. NewVulnerabilityAssessmentScanResultClientWithBaseURI(string, string) VulnerabilityAssessmentScanResultClient
1. NewVulnerabilityAssessmentScansClient(string) VulnerabilityAssessmentScansClient
1. NewVulnerabilityAssessmentScansClientWithBaseURI(string, string) VulnerabilityAssessmentScansClient
1. NewVulnerabilityAssessmentsClient(string) VulnerabilityAssessmentsClient
1. NewVulnerabilityAssessmentsClientWithBaseURI(string, string) VulnerabilityAssessmentsClient
1. NewVulnerabilityAssessmentsSettingsClient(string) VulnerabilityAssessmentsSettingsClient
1. NewVulnerabilityAssessmentsSettingsClientWithBaseURI(string, string) VulnerabilityAssessmentsSettingsClient
1. PossibleAdvancedThreatProtectionStateValues() []AdvancedThreatProtectionState
1. PossibleDNSRefreshOperationStatusValues() []DNSRefreshOperationStatus
1. PossibleDatabaseIdentityTypeValues() []DatabaseIdentityType
1. PossibleMoveOperationModeValues() []MoveOperationMode
1. PossibleReplicationModeValues() []ReplicationMode
1. PossibleRuleSeverityValues() []RuleSeverity
1. PossibleRuleStatusValues() []RuleStatus
1. PossibleRuleTypeValues() []RuleType
1. PossibleServicePrincipalTypeValues() []ServicePrincipalType
1. PossibleSyncGroupsTypeValues() []SyncGroupsType
1. PossibleVulnerabilityAssessmentStateValues() []VulnerabilityAssessmentState
1. QueryCheck.MarshalJSON() ([]byte, error)
1. Remediation.MarshalJSON() ([]byte, error)
1. ServerAdvancedThreatProtection.MarshalJSON() ([]byte, error)
1. ServerAdvancedThreatProtectionSettingsClient.CreateOrUpdate(context.Context, string, string, ServerAdvancedThreatProtection) (ServerAdvancedThreatProtectionSettingsCreateOrUpdateFuture, error)
1. ServerAdvancedThreatProtectionSettingsClient.CreateOrUpdatePreparer(context.Context, string, string, ServerAdvancedThreatProtection) (*http.Request, error)
1. ServerAdvancedThreatProtectionSettingsClient.CreateOrUpdateResponder(*http.Response) (ServerAdvancedThreatProtection, error)
1. ServerAdvancedThreatProtectionSettingsClient.CreateOrUpdateSender(*http.Request) (ServerAdvancedThreatProtectionSettingsCreateOrUpdateFuture, error)
1. ServerAdvancedThreatProtectionSettingsClient.Get(context.Context, string, string) (ServerAdvancedThreatProtection, error)
1. ServerAdvancedThreatProtectionSettingsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. ServerAdvancedThreatProtectionSettingsClient.GetResponder(*http.Response) (ServerAdvancedThreatProtection, error)
1. ServerAdvancedThreatProtectionSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. ServerAdvancedThreatProtectionSettingsClient.ListByServer(context.Context, string, string) (LogicalServerAdvancedThreatProtectionListResultPage, error)
1. ServerAdvancedThreatProtectionSettingsClient.ListByServerComplete(context.Context, string, string) (LogicalServerAdvancedThreatProtectionListResultIterator, error)
1. ServerAdvancedThreatProtectionSettingsClient.ListByServerPreparer(context.Context, string, string) (*http.Request, error)
1. ServerAdvancedThreatProtectionSettingsClient.ListByServerResponder(*http.Response) (LogicalServerAdvancedThreatProtectionListResult, error)
1. ServerAdvancedThreatProtectionSettingsClient.ListByServerSender(*http.Request) (*http.Response, error)
1. ServerConnectionPoliciesClient.ListByServer(context.Context, string, string) (ServerConnectionPolicyListResultPage, error)
1. ServerConnectionPoliciesClient.ListByServerComplete(context.Context, string, string) (ServerConnectionPolicyListResultIterator, error)
1. ServerConnectionPoliciesClient.ListByServerPreparer(context.Context, string, string) (*http.Request, error)
1. ServerConnectionPoliciesClient.ListByServerResponder(*http.Response) (ServerConnectionPolicyListResult, error)
1. ServerConnectionPoliciesClient.ListByServerSender(*http.Request) (*http.Response, error)
1. ServerConnectionPolicyListResult.IsEmpty() bool
1. ServerConnectionPolicyListResult.MarshalJSON() ([]byte, error)
1. ServerConnectionPolicyListResultIterator.NotDone() bool
1. ServerConnectionPolicyListResultIterator.Response() ServerConnectionPolicyListResult
1. ServerConnectionPolicyListResultIterator.Value() ServerConnectionPolicy
1. ServerConnectionPolicyListResultPage.NotDone() bool
1. ServerConnectionPolicyListResultPage.Response() ServerConnectionPolicyListResult
1. ServerConnectionPolicyListResultPage.Values() []ServerConnectionPolicy
1. ServerTrustCertificate.MarshalJSON() ([]byte, error)
1. ServerTrustCertificateProperties.MarshalJSON() ([]byte, error)
1. ServerTrustCertificatesClient.CreateOrUpdate(context.Context, string, string, string, ServerTrustCertificate) (ServerTrustCertificatesCreateOrUpdateFuture, error)
1. ServerTrustCertificatesClient.CreateOrUpdatePreparer(context.Context, string, string, string, ServerTrustCertificate) (*http.Request, error)
1. ServerTrustCertificatesClient.CreateOrUpdateResponder(*http.Response) (ServerTrustCertificate, error)
1. ServerTrustCertificatesClient.CreateOrUpdateSender(*http.Request) (ServerTrustCertificatesCreateOrUpdateFuture, error)
1. ServerTrustCertificatesClient.Delete(context.Context, string, string, string) (ServerTrustCertificatesDeleteFuture, error)
1. ServerTrustCertificatesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ServerTrustCertificatesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ServerTrustCertificatesClient.DeleteSender(*http.Request) (ServerTrustCertificatesDeleteFuture, error)
1. ServerTrustCertificatesClient.Get(context.Context, string, string, string) (ServerTrustCertificate, error)
1. ServerTrustCertificatesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ServerTrustCertificatesClient.GetResponder(*http.Response) (ServerTrustCertificate, error)
1. ServerTrustCertificatesClient.GetSender(*http.Request) (*http.Response, error)
1. ServerTrustCertificatesClient.ListByInstance(context.Context, string, string) (ServerTrustCertificatesListResultPage, error)
1. ServerTrustCertificatesClient.ListByInstanceComplete(context.Context, string, string) (ServerTrustCertificatesListResultIterator, error)
1. ServerTrustCertificatesClient.ListByInstancePreparer(context.Context, string, string) (*http.Request, error)
1. ServerTrustCertificatesClient.ListByInstanceResponder(*http.Response) (ServerTrustCertificatesListResult, error)
1. ServerTrustCertificatesClient.ListByInstanceSender(*http.Request) (*http.Response, error)
1. ServerTrustCertificatesListResult.IsEmpty() bool
1. ServerTrustCertificatesListResult.MarshalJSON() ([]byte, error)
1. ServerTrustCertificatesListResultIterator.NotDone() bool
1. ServerTrustCertificatesListResultIterator.Response() ServerTrustCertificatesListResult
1. ServerTrustCertificatesListResultIterator.Value() ServerTrustCertificate
1. ServerTrustCertificatesListResultPage.NotDone() bool
1. ServerTrustCertificatesListResultPage.Response() ServerTrustCertificatesListResult
1. ServerTrustCertificatesListResultPage.Values() []ServerTrustCertificate
1. ServicePrincipal.MarshalJSON() ([]byte, error)
1. SynapseLinkWorkspace.MarshalJSON() ([]byte, error)
1. SynapseLinkWorkspaceListResult.IsEmpty() bool
1. SynapseLinkWorkspaceListResult.MarshalJSON() ([]byte, error)
1. SynapseLinkWorkspaceListResultIterator.NotDone() bool
1. SynapseLinkWorkspaceListResultIterator.Response() SynapseLinkWorkspaceListResult
1. SynapseLinkWorkspaceListResultIterator.Value() SynapseLinkWorkspace
1. SynapseLinkWorkspaceListResultPage.NotDone() bool
1. SynapseLinkWorkspaceListResultPage.Response() SynapseLinkWorkspaceListResult
1. SynapseLinkWorkspaceListResultPage.Values() []SynapseLinkWorkspace
1. SynapseLinkWorkspacesClient.ListByDatabase(context.Context, string, string, string) (SynapseLinkWorkspaceListResultPage, error)
1. SynapseLinkWorkspacesClient.ListByDatabaseComplete(context.Context, string, string, string) (SynapseLinkWorkspaceListResultIterator, error)
1. SynapseLinkWorkspacesClient.ListByDatabasePreparer(context.Context, string, string, string) (*http.Request, error)
1. SynapseLinkWorkspacesClient.ListByDatabaseResponder(*http.Response) (SynapseLinkWorkspaceListResult, error)
1. SynapseLinkWorkspacesClient.ListByDatabaseSender(*http.Request) (*http.Response, error)
1. TransparentDataEncryptionsClient.ListByDatabase(context.Context, string, string, string) (LogicalDatabaseTransparentDataEncryptionListResultPage, error)
1. TransparentDataEncryptionsClient.ListByDatabaseComplete(context.Context, string, string, string) (LogicalDatabaseTransparentDataEncryptionListResultIterator, error)
1. TransparentDataEncryptionsClient.ListByDatabasePreparer(context.Context, string, string, string) (*http.Request, error)
1. TransparentDataEncryptionsClient.ListByDatabaseResponder(*http.Response) (LogicalDatabaseTransparentDataEncryptionListResult, error)
1. TransparentDataEncryptionsClient.ListByDatabaseSender(*http.Request) (*http.Response, error)
1. UpdateVirtualClusterDNSServersOperation.MarshalJSON() ([]byte, error)
1. VaRule.MarshalJSON() ([]byte, error)
1. VirtualClusterDNSServersProperties.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessment.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessmentBaselineClient.Get(context.Context, string, string) (DatabaseSQLVulnerabilityAssessmentBaselineSet, error)
1. VulnerabilityAssessmentBaselineClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. VulnerabilityAssessmentBaselineClient.GetResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentBaselineSet, error)
1. VulnerabilityAssessmentBaselineClient.GetSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentBaselineClient.ListBySQLVulnerabilityAssessment(context.Context, string, string) (DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage, error)
1. VulnerabilityAssessmentBaselineClient.ListBySQLVulnerabilityAssessmentComplete(context.Context, string, string) (DatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator, error)
1. VulnerabilityAssessmentBaselineClient.ListBySQLVulnerabilityAssessmentPreparer(context.Context, string, string) (*http.Request, error)
1. VulnerabilityAssessmentBaselineClient.ListBySQLVulnerabilityAssessmentResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentBaselineSetListResult, error)
1. VulnerabilityAssessmentBaselineClient.ListBySQLVulnerabilityAssessmentSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentBaselinesClient.CreateOrUpdate(context.Context, string, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput) (DatabaseSQLVulnerabilityAssessmentBaselineSet, error)
1. VulnerabilityAssessmentBaselinesClient.CreateOrUpdatePreparer(context.Context, string, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput) (*http.Request, error)
1. VulnerabilityAssessmentBaselinesClient.CreateOrUpdateResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentBaselineSet, error)
1. VulnerabilityAssessmentBaselinesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentExecuteScanClient.Execute(context.Context, string, string) (VulnerabilityAssessmentExecuteScanExecuteFuture, error)
1. VulnerabilityAssessmentExecuteScanClient.ExecutePreparer(context.Context, string, string) (*http.Request, error)
1. VulnerabilityAssessmentExecuteScanClient.ExecuteResponder(*http.Response) (autorest.Response, error)
1. VulnerabilityAssessmentExecuteScanClient.ExecuteSender(*http.Request) (VulnerabilityAssessmentExecuteScanExecuteFuture, error)
1. VulnerabilityAssessmentListResult.IsEmpty() bool
1. VulnerabilityAssessmentListResult.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessmentListResultIterator.NotDone() bool
1. VulnerabilityAssessmentListResultIterator.Response() VulnerabilityAssessmentListResult
1. VulnerabilityAssessmentListResultIterator.Value() VulnerabilityAssessment
1. VulnerabilityAssessmentListResultPage.NotDone() bool
1. VulnerabilityAssessmentListResultPage.Response() VulnerabilityAssessmentListResult
1. VulnerabilityAssessmentListResultPage.Values() []VulnerabilityAssessment
1. VulnerabilityAssessmentRuleBaselineClient.CreateOrUpdate(context.Context, string, string, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineInput) (DatabaseSQLVulnerabilityAssessmentRuleBaseline, error)
1. VulnerabilityAssessmentRuleBaselineClient.CreateOrUpdatePreparer(context.Context, string, string, string, DatabaseSQLVulnerabilityAssessmentRuleBaselineInput) (*http.Request, error)
1. VulnerabilityAssessmentRuleBaselineClient.CreateOrUpdateResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentRuleBaseline, error)
1. VulnerabilityAssessmentRuleBaselineClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentRuleBaselineClient.Get(context.Context, string, string, string) (DatabaseSQLVulnerabilityAssessmentRuleBaseline, error)
1. VulnerabilityAssessmentRuleBaselineClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. VulnerabilityAssessmentRuleBaselineClient.GetResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentRuleBaseline, error)
1. VulnerabilityAssessmentRuleBaselineClient.GetSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentRuleBaselineClient.ListByBaseline(context.Context, string, string) (DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage, error)
1. VulnerabilityAssessmentRuleBaselineClient.ListByBaselineComplete(context.Context, string, string) (DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator, error)
1. VulnerabilityAssessmentRuleBaselineClient.ListByBaselinePreparer(context.Context, string, string) (*http.Request, error)
1. VulnerabilityAssessmentRuleBaselineClient.ListByBaselineResponder(*http.Response) (DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult, error)
1. VulnerabilityAssessmentRuleBaselineClient.ListByBaselineSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentRuleBaselinesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. VulnerabilityAssessmentRuleBaselinesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. VulnerabilityAssessmentRuleBaselinesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. VulnerabilityAssessmentRuleBaselinesClient.DeleteSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentScanErrorType.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessmentScanListResult.IsEmpty() bool
1. VulnerabilityAssessmentScanListResult.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessmentScanListResultIterator.NotDone() bool
1. VulnerabilityAssessmentScanListResultIterator.Response() VulnerabilityAssessmentScanListResult
1. VulnerabilityAssessmentScanListResultIterator.Value() VulnerabilityAssessmentScanResults
1. VulnerabilityAssessmentScanListResultPage.NotDone() bool
1. VulnerabilityAssessmentScanListResultPage.Response() VulnerabilityAssessmentScanListResult
1. VulnerabilityAssessmentScanListResultPage.Values() []VulnerabilityAssessmentScanResults
1. VulnerabilityAssessmentScanRecordListResultType.IsEmpty() bool
1. VulnerabilityAssessmentScanRecordListResultType.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessmentScanRecordListResultTypeIterator.NotDone() bool
1. VulnerabilityAssessmentScanRecordListResultTypeIterator.Response() VulnerabilityAssessmentScanRecordListResultType
1. VulnerabilityAssessmentScanRecordListResultTypeIterator.Value() VulnerabilityAssessmentScanRecordType
1. VulnerabilityAssessmentScanRecordListResultTypePage.NotDone() bool
1. VulnerabilityAssessmentScanRecordListResultTypePage.Response() VulnerabilityAssessmentScanRecordListResultType
1. VulnerabilityAssessmentScanRecordListResultTypePage.Values() []VulnerabilityAssessmentScanRecordType
1. VulnerabilityAssessmentScanRecordPropertiesType.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessmentScanRecordType.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessmentScanResultClient.Get(context.Context, string, string, string, string) (VulnerabilityAssessmentScanResults, error)
1. VulnerabilityAssessmentScanResultClient.GetPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. VulnerabilityAssessmentScanResultClient.GetResponder(*http.Response) (VulnerabilityAssessmentScanResults, error)
1. VulnerabilityAssessmentScanResultClient.GetSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentScanResultClient.ListByScan(context.Context, string, string, string) (VulnerabilityAssessmentScanListResultPage, error)
1. VulnerabilityAssessmentScanResultClient.ListByScanComplete(context.Context, string, string, string) (VulnerabilityAssessmentScanListResultIterator, error)
1. VulnerabilityAssessmentScanResultClient.ListByScanPreparer(context.Context, string, string, string) (*http.Request, error)
1. VulnerabilityAssessmentScanResultClient.ListByScanResponder(*http.Response) (VulnerabilityAssessmentScanListResult, error)
1. VulnerabilityAssessmentScanResultClient.ListByScanSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentScanResultProperties.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessmentScanResults.MarshalJSON() ([]byte, error)
1. VulnerabilityAssessmentScansClient.Get(context.Context, string, string, string) (VulnerabilityAssessmentScanRecordType, error)
1. VulnerabilityAssessmentScansClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. VulnerabilityAssessmentScansClient.GetResponder(*http.Response) (VulnerabilityAssessmentScanRecordType, error)
1. VulnerabilityAssessmentScansClient.GetSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessments(context.Context, string, string) (VulnerabilityAssessmentScanRecordListResultTypePage, error)
1. VulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessmentsComplete(context.Context, string, string) (VulnerabilityAssessmentScanRecordListResultTypeIterator, error)
1. VulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessmentsPreparer(context.Context, string, string) (*http.Request, error)
1. VulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessmentsResponder(*http.Response) (VulnerabilityAssessmentScanRecordListResultType, error)
1. VulnerabilityAssessmentScansClient.ListBySQLVulnerabilityAssessmentsSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentsClient.Delete(context.Context, string, string) (autorest.Response, error)
1. VulnerabilityAssessmentsClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. VulnerabilityAssessmentsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. VulnerabilityAssessmentsClient.DeleteSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentsSettingsClient.CreateOrUpdate(context.Context, string, string, VulnerabilityAssessment) (VulnerabilityAssessment, error)
1. VulnerabilityAssessmentsSettingsClient.CreateOrUpdatePreparer(context.Context, string, string, VulnerabilityAssessment) (*http.Request, error)
1. VulnerabilityAssessmentsSettingsClient.CreateOrUpdateResponder(*http.Response) (VulnerabilityAssessment, error)
1. VulnerabilityAssessmentsSettingsClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentsSettingsClient.Get(context.Context, string, string) (VulnerabilityAssessment, error)
1. VulnerabilityAssessmentsSettingsClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. VulnerabilityAssessmentsSettingsClient.GetResponder(*http.Response) (VulnerabilityAssessment, error)
1. VulnerabilityAssessmentsSettingsClient.GetSender(*http.Request) (*http.Response, error)
1. VulnerabilityAssessmentsSettingsClient.ListByServer(context.Context, string, string) (VulnerabilityAssessmentListResultPage, error)
1. VulnerabilityAssessmentsSettingsClient.ListByServerComplete(context.Context, string, string) (VulnerabilityAssessmentListResultIterator, error)
1. VulnerabilityAssessmentsSettingsClient.ListByServerPreparer(context.Context, string, string) (*http.Request, error)
1. VulnerabilityAssessmentsSettingsClient.ListByServerResponder(*http.Response) (VulnerabilityAssessmentListResult, error)
1. VulnerabilityAssessmentsSettingsClient.ListByServerSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. AdvancedThreatProtectionProperties
1. Baseline
1. BaselineAdjustedResult
1. BenchmarkReference
1. DatabaseAdvancedThreatProtection
1. DatabaseAdvancedThreatProtectionListResult
1. DatabaseAdvancedThreatProtectionListResultIterator
1. DatabaseAdvancedThreatProtectionListResultPage
1. DatabaseAdvancedThreatProtectionSettingsClient
1. DatabaseIdentity
1. DatabaseSQLVulnerabilityAssessmentBaselineSet
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResult
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResultIterator
1. DatabaseSQLVulnerabilityAssessmentBaselineSetListResultPage
1. DatabaseSQLVulnerabilityAssessmentBaselineSetProperties
1. DatabaseSQLVulnerabilityAssessmentBaselinesClient
1. DatabaseSQLVulnerabilityAssessmentExecuteScanClient
1. DatabaseSQLVulnerabilityAssessmentExecuteScanExecuteFuture
1. DatabaseSQLVulnerabilityAssessmentRuleBaseline
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineInput
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineInputProperties
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListInput
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListInputProperties
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResult
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultIterator
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineListResultPage
1. DatabaseSQLVulnerabilityAssessmentRuleBaselineProperties
1. DatabaseSQLVulnerabilityAssessmentRuleBaselinesClient
1. DatabaseSQLVulnerabilityAssessmentScanResultClient
1. DatabaseSQLVulnerabilityAssessmentScansClient
1. DatabaseSQLVulnerabilityAssessmentsSettingsClient
1. DatabaseUpdateProperties
1. DatabaseUserIdentity
1. DistributedAvailabilityGroup
1. DistributedAvailabilityGroupProperties
1. DistributedAvailabilityGroupsClient
1. DistributedAvailabilityGroupsCreateOrUpdateFuture
1. DistributedAvailabilityGroupsDeleteFuture
1. DistributedAvailabilityGroupsListResult
1. DistributedAvailabilityGroupsListResultIterator
1. DistributedAvailabilityGroupsListResultPage
1. DistributedAvailabilityGroupsUpdateFuture
1. EndpointCertificate
1. EndpointCertificateListResult
1. EndpointCertificateListResultIterator
1. EndpointCertificateListResultPage
1. EndpointCertificateProperties
1. EndpointCertificatesClient
1. IPv6FirewallRule
1. IPv6FirewallRuleListResult
1. IPv6FirewallRuleListResultIterator
1. IPv6FirewallRuleListResultPage
1. IPv6FirewallRulesClient
1. IPv6ServerFirewallRuleProperties
1. LedgerDigestUploadsCreateOrUpdateFuture
1. LedgerDigestUploadsDisableFuture
1. LogicalDatabaseTransparentDataEncryption
1. LogicalDatabaseTransparentDataEncryptionListResult
1. LogicalDatabaseTransparentDataEncryptionListResultIterator
1. LogicalDatabaseTransparentDataEncryptionListResultPage
1. LogicalServerAdvancedThreatProtectionListResult
1. LogicalServerAdvancedThreatProtectionListResultIterator
1. LogicalServerAdvancedThreatProtectionListResultPage
1. ManagedDatabaseAdvancedThreatProtection
1. ManagedDatabaseAdvancedThreatProtectionListResult
1. ManagedDatabaseAdvancedThreatProtectionListResultIterator
1. ManagedDatabaseAdvancedThreatProtectionListResultPage
1. ManagedDatabaseAdvancedThreatProtectionSettingsClient
1. ManagedDatabaseMoveDefinition
1. ManagedDatabaseMoveOperationListResult
1. ManagedDatabaseMoveOperationListResultIterator
1. ManagedDatabaseMoveOperationListResultPage
1. ManagedDatabaseMoveOperationResult
1. ManagedDatabaseMoveOperationResultProperties
1. ManagedDatabaseMoveOperationsClient
1. ManagedDatabaseRestoreDetailsBackupSetProperties
1. ManagedDatabaseRestoreDetailsUnrestorableFileProperties
1. ManagedDatabaseStartMoveDefinition
1. ManagedDatabasesCancelMoveFuture
1. ManagedDatabasesCompleteMoveFuture
1. ManagedDatabasesStartMoveFuture
1. ManagedInstanceAdvancedThreatProtection
1. ManagedInstanceAdvancedThreatProtectionListResult
1. ManagedInstanceAdvancedThreatProtectionListResultIterator
1. ManagedInstanceAdvancedThreatProtectionListResultPage
1. ManagedInstanceAdvancedThreatProtectionSettingsClient
1. ManagedInstanceAdvancedThreatProtectionSettingsCreateOrUpdateFuture
1. ManagedInstanceDtc
1. ManagedInstanceDtcListResult
1. ManagedInstanceDtcListResultIterator
1. ManagedInstanceDtcListResultPage
1. ManagedInstanceDtcProperties
1. ManagedInstanceDtcSecuritySettings
1. ManagedInstanceDtcTransactionManagerCommunicationSettings
1. ManagedInstanceDtcsClient
1. ManagedInstanceDtcsCreateOrUpdateFuture
1. ManagedServerDNSAlias
1. ManagedServerDNSAliasAcquisition
1. ManagedServerDNSAliasCreation
1. ManagedServerDNSAliasListResult
1. ManagedServerDNSAliasListResultIterator
1. ManagedServerDNSAliasListResultPage
1. ManagedServerDNSAliasProperties
1. ManagedServerDNSAliasesAcquireFuture
1. ManagedServerDNSAliasesClient
1. ManagedServerDNSAliasesCreateOrUpdateFuture
1. ManagedServerDNSAliasesDeleteFuture
1. QueryCheck
1. Remediation
1. ReplicationLinksDeleteFuture
1. ServerAdvancedThreatProtection
1. ServerAdvancedThreatProtectionSettingsClient
1. ServerAdvancedThreatProtectionSettingsCreateOrUpdateFuture
1. ServerConnectionPoliciesCreateOrUpdateFuture
1. ServerConnectionPolicyListResult
1. ServerConnectionPolicyListResultIterator
1. ServerConnectionPolicyListResultPage
1. ServerTrustCertificate
1. ServerTrustCertificateProperties
1. ServerTrustCertificatesClient
1. ServerTrustCertificatesCreateOrUpdateFuture
1. ServerTrustCertificatesDeleteFuture
1. ServerTrustCertificatesListResult
1. ServerTrustCertificatesListResultIterator
1. ServerTrustCertificatesListResultPage
1. ServicePrincipal
1. SynapseLinkWorkspace
1. SynapseLinkWorkspaceInfoProperties
1. SynapseLinkWorkspaceListResult
1. SynapseLinkWorkspaceListResultIterator
1. SynapseLinkWorkspaceListResultPage
1. SynapseLinkWorkspaceProperties
1. SynapseLinkWorkspacesClient
1. UpdateVirtualClusterDNSServersOperation
1. VaRule
1. VirtualClusterDNSServersProperties
1. VirtualClustersUpdateDNSServersFuture
1. VulnerabilityAssessment
1. VulnerabilityAssessmentBaselineClient
1. VulnerabilityAssessmentBaselinesClient
1. VulnerabilityAssessmentExecuteScanClient
1. VulnerabilityAssessmentExecuteScanExecuteFuture
1. VulnerabilityAssessmentListResult
1. VulnerabilityAssessmentListResultIterator
1. VulnerabilityAssessmentListResultPage
1. VulnerabilityAssessmentPolicyProperties
1. VulnerabilityAssessmentRuleBaselineClient
1. VulnerabilityAssessmentRuleBaselinesClient
1. VulnerabilityAssessmentScanErrorType
1. VulnerabilityAssessmentScanListResult
1. VulnerabilityAssessmentScanListResultIterator
1. VulnerabilityAssessmentScanListResultPage
1. VulnerabilityAssessmentScanRecordListResultType
1. VulnerabilityAssessmentScanRecordListResultTypeIterator
1. VulnerabilityAssessmentScanRecordListResultTypePage
1. VulnerabilityAssessmentScanRecordPropertiesType
1. VulnerabilityAssessmentScanRecordType
1. VulnerabilityAssessmentScanResultClient
1. VulnerabilityAssessmentScanResultProperties
1. VulnerabilityAssessmentScanResults
1. VulnerabilityAssessmentScansClient
1. VulnerabilityAssessmentsClient
1. VulnerabilityAssessmentsSettingsClient

#### New Struct Fields

1. Database.Identity
1. DatabaseBlobAuditingPolicyProperties.IsManagedIdentityInUse
1. DatabaseProperties.FederatedClientID
1. DatabaseProperties.SourceResourceID
1. DatabaseUpdate.*DatabaseUpdateProperties
1. DatabaseUpdate.Identity
1. ElasticPoolProperties.HighAvailabilityReplicaCount
1. ElasticPoolUpdateProperties.HighAvailabilityReplicaCount
1. ExtendedDatabaseBlobAuditingPolicyProperties.IsManagedIdentityInUse
1. ExtendedServerBlobAuditingPolicyProperties.IsManagedIdentityInUse
1. ManagedDatabaseProperties.StorageContainerIdentity
1. ManagedDatabaseRestoreDetailsProperties.CurrentBackupType
1. ManagedDatabaseRestoreDetailsProperties.CurrentRestorePlanSizeMB
1. ManagedDatabaseRestoreDetailsProperties.CurrentRestoredSizeMB
1. ManagedDatabaseRestoreDetailsProperties.DiffBackupSets
1. ManagedDatabaseRestoreDetailsProperties.FullBackupSets
1. ManagedDatabaseRestoreDetailsProperties.LogBackupSets
1. ManagedDatabaseRestoreDetailsProperties.NumberOfFilesQueued
1. ManagedDatabaseRestoreDetailsProperties.NumberOfFilesRestored
1. ManagedDatabaseRestoreDetailsProperties.NumberOfFilesRestoring
1. ManagedDatabaseRestoreDetailsProperties.NumberOfFilesSkipped
1. ManagedDatabaseRestoreDetailsProperties.NumberOfFilesUnrestorable
1. ManagedDatabaseRestoreDetailsProperties.Type
1. ManagedInstanceProperties.CurrentBackupStorageRedundancy
1. ManagedInstanceProperties.RequestedBackupStorageRedundancy
1. ManagedInstanceProperties.ServicePrincipal
1. ServerBlobAuditingPolicyProperties.IsManagedIdentityInUse
1. TransparentDataEncryptionProperties.State
1. VirtualClusterProperties.Version
