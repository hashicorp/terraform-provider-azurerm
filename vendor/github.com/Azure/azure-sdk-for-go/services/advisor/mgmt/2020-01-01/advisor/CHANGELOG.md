Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewMetadataEntityListResultPage` parameter(s) have been changed from `(func(context.Context, MetadataEntityListResult) (MetadataEntityListResult, error))` to `(MetadataEntityListResult, func(context.Context, MetadataEntityListResult) (MetadataEntityListResult, error))`
- Function `SuppressionsClient.Get` return value(s) have been changed from `(SuppressionContract, error)` to `(SetObject, error)`
- Function `NewOperationEntityListResultPage` parameter(s) have been changed from `(func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))` to `(OperationEntityListResult, func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))`
- Function `SuppressionsClient.GetResponder` return value(s) have been changed from `(SuppressionContract, error)` to `(SetObject, error)`
- Function `NewResourceRecommendationBaseListResultPage` parameter(s) have been changed from `(func(context.Context, ResourceRecommendationBaseListResult) (ResourceRecommendationBaseListResult, error))` to `(ResourceRecommendationBaseListResult, func(context.Context, ResourceRecommendationBaseListResult) (ResourceRecommendationBaseListResult, error))`
- Function `NewSuppressionContractListResultPage` parameter(s) have been changed from `(func(context.Context, SuppressionContractListResult) (SuppressionContractListResult, error))` to `(SuppressionContractListResult, func(context.Context, SuppressionContractListResult) (SuppressionContractListResult, error))`
- Function `NewConfigurationListResultPage` parameter(s) have been changed from `(func(context.Context, ConfigurationListResult) (ConfigurationListResult, error))` to `(ConfigurationListResult, func(context.Context, ConfigurationListResult) (ConfigurationListResult, error))`

## New Content

- New function `SuppressionProperties.MarshalJSON() ([]byte, error)`
- New field `ExpirationTimeStamp` in struct `SuppressionProperties`
