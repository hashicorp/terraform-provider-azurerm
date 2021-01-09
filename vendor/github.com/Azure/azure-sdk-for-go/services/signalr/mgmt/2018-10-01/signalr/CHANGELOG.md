Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewResourceListPage` parameter(s) have been changed from `(func(context.Context, ResourceList) (ResourceList, error))` to `(ResourceList, func(context.Context, ResourceList) (ResourceList, error))`
- Function `NewOperationListPage` parameter(s) have been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList, func(context.Context, OperationList) (OperationList, error))`
- Function `NewUsageListPage` parameter(s) have been changed from `(func(context.Context, UsageList) (UsageList, error))` to `(UsageList, func(context.Context, UsageList) (UsageList, error))`
