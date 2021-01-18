Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewPolicyTrackedResourcesQueryResultsPage` parameter(s) have been changed from `(func(context.Context, PolicyTrackedResourcesQueryResults) (PolicyTrackedResourcesQueryResults, error))` to `(PolicyTrackedResourcesQueryResults, func(context.Context, PolicyTrackedResourcesQueryResults) (PolicyTrackedResourcesQueryResults, error))`
- Function `NewPolicyStatesQueryResultsPage` parameter(s) have been changed from `(func(context.Context, PolicyStatesQueryResults) (PolicyStatesQueryResults, error))` to `(PolicyStatesQueryResults, func(context.Context, PolicyStatesQueryResults) (PolicyStatesQueryResults, error))`
- Function `NewPolicyMetadataCollectionPage` parameter(s) have been changed from `(func(context.Context, PolicyMetadataCollection) (PolicyMetadataCollection, error))` to `(PolicyMetadataCollection, func(context.Context, PolicyMetadataCollection) (PolicyMetadataCollection, error))`
- Function `NewRemediationDeploymentsListResultPage` parameter(s) have been changed from `(func(context.Context, RemediationDeploymentsListResult) (RemediationDeploymentsListResult, error))` to `(RemediationDeploymentsListResult, func(context.Context, RemediationDeploymentsListResult) (RemediationDeploymentsListResult, error))`
- Function `NewRemediationListResultPage` parameter(s) have been changed from `(func(context.Context, RemediationListResult) (RemediationListResult, error))` to `(RemediationListResult, func(context.Context, RemediationListResult) (RemediationListResult, error))`
- Function `NewPolicyEventsQueryResultsPage` parameter(s) have been changed from `(func(context.Context, PolicyEventsQueryResults) (PolicyEventsQueryResults, error))` to `(PolicyEventsQueryResults, func(context.Context, PolicyEventsQueryResults) (PolicyEventsQueryResults, error))`

## New Content

- New function `ExpressionEvaluationDetails.MarshalJSON() ([]byte, error)`
- New field `ExpressionKind` in struct `ExpressionEvaluationDetails`
