Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewUsageModelsResultPage` parameter(s) have been changed from `(func(context.Context, UsageModelsResult) (UsageModelsResult, error))` to `(UsageModelsResult, func(context.Context, UsageModelsResult) (UsageModelsResult, error))`
- Function `NewStorageTargetsResultPage` parameter(s) have been changed from `(func(context.Context, StorageTargetsResult) (StorageTargetsResult, error))` to `(StorageTargetsResult, func(context.Context, StorageTargetsResult) (StorageTargetsResult, error))`
- Function `NewResourceSkusResultPage` parameter(s) have been changed from `(func(context.Context, ResourceSkusResult) (ResourceSkusResult, error))` to `(ResourceSkusResult, func(context.Context, ResourceSkusResult) (ResourceSkusResult, error))`
- Function `NewAPIOperationListResultPage` parameter(s) have been changed from `(func(context.Context, APIOperationListResult) (APIOperationListResult, error))` to `(APIOperationListResult, func(context.Context, APIOperationListResult) (APIOperationListResult, error))`
- Function `NewCachesListResultPage` parameter(s) have been changed from `(func(context.Context, CachesListResult) (CachesListResult, error))` to `(CachesListResult, func(context.Context, CachesListResult) (CachesListResult, error))`

## New Content

- New const `MetricAggregationTypeAverage`
- New const `MetricAggregationTypeMaximum`
- New const `MetricAggregationTypeTotal`
- New const `Key`
- New const `MetricAggregationTypeNone`
- New const `ManagedIdentity`
- New const `Application`
- New const `MetricAggregationTypeNotSpecified`
- New const `MetricAggregationTypeMinimum`
- New const `MetricAggregationTypeCount`
- New const `User`
- New function `APIOperation.MarshalJSON() ([]byte, error)`
- New function `PossibleCreatedByTypeValues() []CreatedByType`
- New function `*APIOperation.UnmarshalJSON([]byte) error`
- New function `PossibleMetricAggregationTypeValues() []MetricAggregationType`
- New struct `APIOperationProperties`
- New struct `APIOperationPropertiesServiceSpecification`
- New struct `MetricDimension`
- New struct `MetricSpecification`
- New struct `SystemData`
- New field `Location` in struct `StorageTargetResource`
- New field `SystemData` in struct `StorageTargetResource`
- New field `SystemData` in struct `Cache`
- New anonymous field `*APIOperationProperties` in struct `APIOperation`
- New field `Origin` in struct `APIOperation`
- New field `IsDataAction` in struct `APIOperation`
- New field `Description` in struct `APIOperationDisplay`
- New field `Location` in struct `StorageTarget`
- New field `SystemData` in struct `StorageTarget`
