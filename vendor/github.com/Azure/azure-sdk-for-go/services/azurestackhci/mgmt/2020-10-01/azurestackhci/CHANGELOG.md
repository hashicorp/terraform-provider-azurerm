Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewClusterListPage` parameter(s) have been changed from `(func(context.Context, ClusterList) (ClusterList, error))` to `(ClusterList, func(context.Context, ClusterList) (ClusterList, error))`
- Type of `ErrorResponse.Error` has been changed from `*ErrorResponseError` to `*ErrorDetail`
- Struct `ErrorResponseError` has been removed

## New Content

- New struct `ErrorDetail`
