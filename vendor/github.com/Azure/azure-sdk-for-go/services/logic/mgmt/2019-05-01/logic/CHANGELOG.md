# Change History

## Breaking Changes

### Removed Constants

1. AgreementType.AS2
1. AgreementType.Edifact
1. AgreementType.NotSpecified
1. AgreementType.X12
1. AzureAsyncOperationState.Canceled
1. AzureAsyncOperationState.Failed
1. AzureAsyncOperationState.Pending
1. AzureAsyncOperationState.Succeeded
1. DayOfWeek.Friday
1. DayOfWeek.Monday
1. DayOfWeek.Saturday
1. DayOfWeek.Sunday
1. DayOfWeek.Thursday
1. DayOfWeek.Tuesday
1. DayOfWeek.Wednesday
1. EventLevel.Critical
1. EventLevel.Error
1. EventLevel.Informational
1. EventLevel.LogAlways
1. EventLevel.Verbose
1. EventLevel.Warning
1. IntegrationServiceEnvironmentSkuScaleType.Automatic
1. IntegrationServiceEnvironmentSkuScaleType.Manual
1. IntegrationServiceEnvironmentSkuScaleType.None
1. OpenAuthenticationProviderType.AAD
1. SwaggerSchemaType.Array
1. SwaggerSchemaType.Boolean
1. SwaggerSchemaType.File
1. SwaggerSchemaType.Integer
1. SwaggerSchemaType.Null
1. SwaggerSchemaType.Number
1. SwaggerSchemaType.Object
1. SwaggerSchemaType.String

### Removed Funcs

1. *ManagedAPIListResultIterator.Next() error
1. *ManagedAPIListResultIterator.NextWithContext(context.Context) error
1. *ManagedAPIListResultPage.Next() error
1. *ManagedAPIListResultPage.NextWithContext(context.Context) error
1. ManagedAPIListResult.IsEmpty() bool
1. ManagedAPIListResultIterator.NotDone() bool
1. ManagedAPIListResultIterator.Response() ManagedAPIListResult
1. ManagedAPIListResultIterator.Value() ManagedAPI
1. ManagedAPIListResultPage.NotDone() bool
1. ManagedAPIListResultPage.Response() ManagedAPIListResult
1. ManagedAPIListResultPage.Values() []ManagedAPI
1. NewManagedAPIListResultIterator(ManagedAPIListResultPage) ManagedAPIListResultIterator
1. NewManagedAPIListResultPage(ManagedAPIListResult, func(context.Context, ManagedAPIListResult) (ManagedAPIListResult, error)) ManagedAPIListResultPage
1. OpenAuthenticationAccessPolicy.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. ManagedAPIListResultIterator
1. ManagedAPIListResultPage

#### Removed Struct Fields

1. ManagedAPI.autorest.Response
1. ManagedAPIListResult.autorest.Response

### Signature Changes

#### Funcs

1. IntegrationServiceEnvironmentManagedApisClient.Get
	- Returns
		- From: ManagedAPI, error
		- To: IntegrationServiceEnvironmentManagedAPI, error
1. IntegrationServiceEnvironmentManagedApisClient.GetResponder
	- Returns
		- From: ManagedAPI, error
		- To: IntegrationServiceEnvironmentManagedAPI, error
1. IntegrationServiceEnvironmentManagedApisClient.List
	- Returns
		- From: ManagedAPIListResultPage, error
		- To: IntegrationServiceEnvironmentManagedAPIListResultPage, error
1. IntegrationServiceEnvironmentManagedApisClient.ListComplete
	- Returns
		- From: ManagedAPIListResultIterator, error
		- To: IntegrationServiceEnvironmentManagedAPIListResultIterator, error
1. IntegrationServiceEnvironmentManagedApisClient.ListResponder
	- Returns
		- From: ManagedAPIListResult, error
		- To: IntegrationServiceEnvironmentManagedAPIListResult, error
1. IntegrationServiceEnvironmentManagedApisClient.Put
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, IntegrationServiceEnvironmentManagedAPI
1. IntegrationServiceEnvironmentManagedApisClient.PutPreparer
	- Params
		- From: context.Context, string, string, string
		- To: context.Context, string, string, string, IntegrationServiceEnvironmentManagedAPI
1. IntegrationServiceEnvironmentManagedApisClient.PutResponder
	- Returns
		- From: ManagedAPI, error
		- To: IntegrationServiceEnvironmentManagedAPI, error

#### Struct Fields

1. IntegrationAccountProperties.IntegrationServiceEnvironment changed type from *IntegrationServiceEnvironment to *ResourceReference
1. IntegrationServiceEnvironmentManagedApisPutFuture.Result changed type from func(IntegrationServiceEnvironmentManagedApisClient) (ManagedAPI, error) to func(IntegrationServiceEnvironmentManagedApisClient) (IntegrationServiceEnvironmentManagedAPI, error)

## Additive Changes

### New Constants

1. AgreementType.AgreementTypeAS2
1. AgreementType.AgreementTypeEdifact
1. AgreementType.AgreementTypeNotSpecified
1. AgreementType.AgreementTypeX12
1. AzureAsyncOperationState.AzureAsyncOperationStateCanceled
1. AzureAsyncOperationState.AzureAsyncOperationStateFailed
1. AzureAsyncOperationState.AzureAsyncOperationStatePending
1. AzureAsyncOperationState.AzureAsyncOperationStateSucceeded
1. DayOfWeek.DayOfWeekFriday
1. DayOfWeek.DayOfWeekMonday
1. DayOfWeek.DayOfWeekSaturday
1. DayOfWeek.DayOfWeekSunday
1. DayOfWeek.DayOfWeekThursday
1. DayOfWeek.DayOfWeekTuesday
1. DayOfWeek.DayOfWeekWednesday
1. EventLevel.EventLevelCritical
1. EventLevel.EventLevelError
1. EventLevel.EventLevelInformational
1. EventLevel.EventLevelLogAlways
1. EventLevel.EventLevelVerbose
1. EventLevel.EventLevelWarning
1. IntegrationServiceEnvironmentSkuScaleType.IntegrationServiceEnvironmentSkuScaleTypeAutomatic
1. IntegrationServiceEnvironmentSkuScaleType.IntegrationServiceEnvironmentSkuScaleTypeManual
1. IntegrationServiceEnvironmentSkuScaleType.IntegrationServiceEnvironmentSkuScaleTypeNone
1. ManagedServiceIdentityType.ManagedServiceIdentityTypeNone
1. ManagedServiceIdentityType.ManagedServiceIdentityTypeSystemAssigned
1. ManagedServiceIdentityType.ManagedServiceIdentityTypeUserAssigned
1. OpenAuthenticationProviderType.OpenAuthenticationProviderTypeAAD
1. SwaggerSchemaType.SwaggerSchemaTypeArray
1. SwaggerSchemaType.SwaggerSchemaTypeBoolean
1. SwaggerSchemaType.SwaggerSchemaTypeFile
1. SwaggerSchemaType.SwaggerSchemaTypeInteger
1. SwaggerSchemaType.SwaggerSchemaTypeNull
1. SwaggerSchemaType.SwaggerSchemaTypeNumber
1. SwaggerSchemaType.SwaggerSchemaTypeObject
1. SwaggerSchemaType.SwaggerSchemaTypeString

### New Funcs

1. *IntegrationServiceEnvironmentManagedAPI.UnmarshalJSON([]byte) error
1. *IntegrationServiceEnvironmentManagedAPIListResultIterator.Next() error
1. *IntegrationServiceEnvironmentManagedAPIListResultIterator.NextWithContext(context.Context) error
1. *IntegrationServiceEnvironmentManagedAPIListResultPage.Next() error
1. *IntegrationServiceEnvironmentManagedAPIListResultPage.NextWithContext(context.Context) error
1. ContentLink.MarshalJSON() ([]byte, error)
1. IntegrationServiceEnvironmentManagedAPI.MarshalJSON() ([]byte, error)
1. IntegrationServiceEnvironmentManagedAPIListResult.IsEmpty() bool
1. IntegrationServiceEnvironmentManagedAPIListResultIterator.NotDone() bool
1. IntegrationServiceEnvironmentManagedAPIListResultIterator.Response() IntegrationServiceEnvironmentManagedAPIListResult
1. IntegrationServiceEnvironmentManagedAPIListResultIterator.Value() IntegrationServiceEnvironmentManagedAPI
1. IntegrationServiceEnvironmentManagedAPIListResultPage.NotDone() bool
1. IntegrationServiceEnvironmentManagedAPIListResultPage.Response() IntegrationServiceEnvironmentManagedAPIListResult
1. IntegrationServiceEnvironmentManagedAPIListResultPage.Values() []IntegrationServiceEnvironmentManagedAPI
1. IntegrationServiceEnvironmentManagedAPIProperties.MarshalJSON() ([]byte, error)
1. ManagedServiceIdentity.MarshalJSON() ([]byte, error)
1. NewIntegrationServiceEnvironmentManagedAPIListResultIterator(IntegrationServiceEnvironmentManagedAPIListResultPage) IntegrationServiceEnvironmentManagedAPIListResultIterator
1. NewIntegrationServiceEnvironmentManagedAPIListResultPage(IntegrationServiceEnvironmentManagedAPIListResult, func(context.Context, IntegrationServiceEnvironmentManagedAPIListResult) (IntegrationServiceEnvironmentManagedAPIListResult, error)) IntegrationServiceEnvironmentManagedAPIListResultPage
1. PossibleManagedServiceIdentityTypeValues() []ManagedServiceIdentityType
1. UserAssignedIdentity.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. IntegrationServiceEnvironmentManagedAPI
1. IntegrationServiceEnvironmentManagedAPIDeploymentParameters
1. IntegrationServiceEnvironmentManagedAPIListResult
1. IntegrationServiceEnvironmentManagedAPIListResultIterator
1. IntegrationServiceEnvironmentManagedAPIListResultPage
1. IntegrationServiceEnvironmentManagedAPIProperties
1. ManagedServiceIdentity
1. UserAssignedIdentity

#### New Struct Fields

1. IntegrationServiceEnvironment.Identity
1. Workflow.Identity
