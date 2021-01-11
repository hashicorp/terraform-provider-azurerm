Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewResourceSkusResultPage` parameter(s) have been changed from `(func(context.Context, ResourceSkusResult) (ResourceSkusResult, error))` to `(ResourceSkusResult, func(context.Context, ResourceSkusResult) (ResourceSkusResult, error))`
- Function `NewAccountListResultPage` parameter(s) have been changed from `(func(context.Context, AccountListResult) (AccountListResult, error))` to `(AccountListResult, func(context.Context, AccountListResult) (AccountListResult, error))`
- Function `NewOperationEntityListResultPage` parameter(s) have been changed from `(func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))` to `(OperationEntityListResult, func(context.Context, OperationEntityListResult) (OperationEntityListResult, error))`

## New Content

- New const `Enterprise`
- New field `DateCreated` in struct `AccountProperties`
