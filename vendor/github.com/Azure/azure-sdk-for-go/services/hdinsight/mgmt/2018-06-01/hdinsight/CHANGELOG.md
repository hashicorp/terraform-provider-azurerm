# Change History

## Breaking Changes

### Removed Constants

1. AsyncOperationState.Failed
1. AsyncOperationState.InProgress
1. AsyncOperationState.Succeeded
1. DaysOfWeek.Friday
1. DaysOfWeek.Monday
1. DaysOfWeek.Saturday
1. DaysOfWeek.Sunday
1. DaysOfWeek.Thursday
1. DaysOfWeek.Tuesday
1. DaysOfWeek.Wednesday
1. DirectoryType.ActiveDirectory
1. FilterMode.Default
1. FilterMode.Exclude
1. FilterMode.Include
1. FilterMode.Recommend
1. JSONWebKeyEncryptionAlgorithm.RSA15
1. JSONWebKeyEncryptionAlgorithm.RSAOAEP
1. JSONWebKeyEncryptionAlgorithm.RSAOAEP256
1. OSType.Linux
1. OSType.Windows
1. PrivateLink.Disabled
1. PrivateLink.Enabled
1. ResourceIdentityType.None
1. ResourceIdentityType.SystemAssigned
1. ResourceIdentityType.SystemAssignedUserAssigned
1. ResourceIdentityType.UserAssigned
1. ResourceProviderConnection.Inbound
1. ResourceProviderConnection.Outbound
1. Tier.Premium
1. Tier.Standard

### Struct Changes

#### Removed Structs

1. OperationResource

#### Removed Struct Fields

1. CapabilitiesResult.VMSizeFilters
1. CapabilitiesResult.VMSizes
1. Extension.autorest.Response
1. VMSizeCompatibilityFilter.Vmsizes

### Signature Changes

#### Funcs

1. ExtensionsClient.Get
	- Returns
		- From: Extension, error
		- To: ClusterMonitoringResponse, error
1. ExtensionsClient.GetResponder
	- Returns
		- From: Extension, error
		- To: ClusterMonitoringResponse, error

#### Struct Fields

1. Usage.CurrentValue changed type from *int32 to *int64
1. Usage.Limit changed type from *int32 to *int64
1. VersionSpec.IsDefault changed type from *string to *bool

## Additive Changes

### New Constants

1. AsyncOperationState.AsyncOperationStateFailed
1. AsyncOperationState.AsyncOperationStateInProgress
1. AsyncOperationState.AsyncOperationStateSucceeded
1. DaysOfWeek.DaysOfWeekFriday
1. DaysOfWeek.DaysOfWeekMonday
1. DaysOfWeek.DaysOfWeekSaturday
1. DaysOfWeek.DaysOfWeekSunday
1. DaysOfWeek.DaysOfWeekThursday
1. DaysOfWeek.DaysOfWeekTuesday
1. DaysOfWeek.DaysOfWeekWednesday
1. DirectoryType.DirectoryTypeActiveDirectory
1. FilterMode.FilterModeDefault
1. FilterMode.FilterModeExclude
1. FilterMode.FilterModeInclude
1. FilterMode.FilterModeRecommend
1. JSONWebKeyEncryptionAlgorithm.JSONWebKeyEncryptionAlgorithmRSA15
1. JSONWebKeyEncryptionAlgorithm.JSONWebKeyEncryptionAlgorithmRSAOAEP
1. JSONWebKeyEncryptionAlgorithm.JSONWebKeyEncryptionAlgorithmRSAOAEP256
1. OSType.OSTypeLinux
1. OSType.OSTypeWindows
1. PrivateLink.PrivateLinkDisabled
1. PrivateLink.PrivateLinkEnabled
1. ResourceIdentityType.ResourceIdentityTypeNone
1. ResourceIdentityType.ResourceIdentityTypeSystemAssigned
1. ResourceIdentityType.ResourceIdentityTypeSystemAssignedUserAssigned
1. ResourceIdentityType.ResourceIdentityTypeUserAssigned
1. ResourceProviderConnection.ResourceProviderConnectionInbound
1. ResourceProviderConnection.ResourceProviderConnectionOutbound
1. Tier.TierPremium
1. Tier.TierStandard

### New Funcs

1. *ClustersUpdateIdentityCertificateFuture.UnmarshalJSON([]byte) error
1. ApplicationGetHTTPSEndpoint.MarshalJSON() ([]byte, error)
1. ApplicationsClient.GetAzureAsyncOperationStatus(context.Context, string, string, string, string) (AsyncOperationResult, error)
1. ApplicationsClient.GetAzureAsyncOperationStatusPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ApplicationsClient.GetAzureAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. ApplicationsClient.GetAzureAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. ClusterCreateRequestValidationParameters.MarshalJSON() ([]byte, error)
1. ClustersClient.GetAzureAsyncOperationStatus(context.Context, string, string, string) (AsyncOperationResult, error)
1. ClustersClient.GetAzureAsyncOperationStatusPreparer(context.Context, string, string, string) (*http.Request, error)
1. ClustersClient.GetAzureAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. ClustersClient.GetAzureAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. ClustersClient.UpdateIdentityCertificate(context.Context, string, string, UpdateClusterIdentityCertificateParameters) (ClustersUpdateIdentityCertificateFuture, error)
1. ClustersClient.UpdateIdentityCertificatePreparer(context.Context, string, string, UpdateClusterIdentityCertificateParameters) (*http.Request, error)
1. ClustersClient.UpdateIdentityCertificateResponder(*http.Response) (autorest.Response, error)
1. ClustersClient.UpdateIdentityCertificateSender(*http.Request) (ClustersUpdateIdentityCertificateFuture, error)
1. ExtensionsClient.GetAzureAsyncOperationStatus(context.Context, string, string, string, string) (AsyncOperationResult, error)
1. ExtensionsClient.GetAzureAsyncOperationStatusPreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ExtensionsClient.GetAzureAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. ExtensionsClient.GetAzureAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. KafkaRestProperties.MarshalJSON() ([]byte, error)
1. LocationsClient.CheckNameAvailability(context.Context, string, NameAvailabilityCheckRequestParameters) (NameAvailabilityCheckResult, error)
1. LocationsClient.CheckNameAvailabilityPreparer(context.Context, string, NameAvailabilityCheckRequestParameters) (*http.Request, error)
1. LocationsClient.CheckNameAvailabilityResponder(*http.Response) (NameAvailabilityCheckResult, error)
1. LocationsClient.CheckNameAvailabilitySender(*http.Request) (*http.Response, error)
1. LocationsClient.GetAzureAsyncOperationStatus(context.Context, string, string) (AsyncOperationResult, error)
1. LocationsClient.GetAzureAsyncOperationStatusPreparer(context.Context, string, string) (*http.Request, error)
1. LocationsClient.GetAzureAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. LocationsClient.GetAzureAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. LocationsClient.ValidateClusterCreateRequest(context.Context, string, ClusterCreateRequestValidationParameters) (ClusterCreateValidationResult, error)
1. LocationsClient.ValidateClusterCreateRequestPreparer(context.Context, string, ClusterCreateRequestValidationParameters) (*http.Request, error)
1. LocationsClient.ValidateClusterCreateRequestResponder(*http.Response) (ClusterCreateValidationResult, error)
1. LocationsClient.ValidateClusterCreateRequestSender(*http.Request) (*http.Response, error)
1. NameAvailabilityCheckResult.MarshalJSON() ([]byte, error)
1. ScriptActionsClient.GetExecutionAsyncOperationStatus(context.Context, string, string, string) (AsyncOperationResult, error)
1. ScriptActionsClient.GetExecutionAsyncOperationStatusPreparer(context.Context, string, string, string) (*http.Request, error)
1. ScriptActionsClient.GetExecutionAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. ScriptActionsClient.GetExecutionAsyncOperationStatusSender(*http.Request) (*http.Response, error)
1. VirtualMachinesClient.GetAsyncOperationStatus(context.Context, string, string, string) (AsyncOperationResult, error)
1. VirtualMachinesClient.GetAsyncOperationStatusPreparer(context.Context, string, string, string) (*http.Request, error)
1. VirtualMachinesClient.GetAsyncOperationStatusResponder(*http.Response) (AsyncOperationResult, error)
1. VirtualMachinesClient.GetAsyncOperationStatusSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. AaddsResourceDetails
1. AsyncOperationResult
1. ClusterCreateRequestValidationParameters
1. ClusterCreateValidationResult
1. ClustersUpdateIdentityCertificateFuture
1. NameAvailabilityCheckRequestParameters
1. NameAvailabilityCheckResult
1. UpdateClusterIdentityCertificateParameters
1. ValidationErrorInfo

#### New Struct Fields

1. ApplicationGetHTTPSEndpoint.PrivateIPAddress
1. CapabilitiesResult.VmsizeFilters
1. CapabilitiesResult.Vmsizes
1. KafkaRestProperties.ConfigurationOverride
1. Role.VMGroupName
1. StorageAccount.Fileshare
1. StorageAccount.Saskey
1. VMSizeCompatibilityFilter.ComputeIsolationSupported
1. VMSizeCompatibilityFilter.ESPApplied
1. VMSizeCompatibilityFilter.OsType
1. VMSizeCompatibilityFilter.VMSizes
