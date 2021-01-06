Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewDataSourceListResultPage` parameter(s) have been changed from `(func(context.Context, DataSourceListResult) (DataSourceListResult, error))` to `(DataSourceListResult, func(context.Context, DataSourceListResult) (DataSourceListResult, error))`
- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `NewClusterListResultPage` parameter(s) have been changed from `(func(context.Context, ClusterListResult) (ClusterListResult, error))` to `(ClusterListResult, func(context.Context, ClusterListResult) (ClusterListResult, error))`
- Function `NewStorageInsightListResultPage` parameter(s) have been changed from `(func(context.Context, StorageInsightListResult) (StorageInsightListResult, error))` to `(StorageInsightListResult, func(context.Context, StorageInsightListResult) (StorageInsightListResult, error))`
- Field `AllTables` of struct `DataExportProperties` has been removed
