Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewListPage` parameter(s) have been changed from `(func(context.Context, List) (List, error))` to `(List, func(context.Context, List) (List, error))`
- Function `NewAssignmentListPage` parameter(s) have been changed from `(func(context.Context, AssignmentList) (AssignmentList, error))` to `(AssignmentList, func(context.Context, AssignmentList) (AssignmentList, error))`
- Function `NewArtifactListPage` parameter(s) have been changed from `(func(context.Context, ArtifactList) (ArtifactList, error))` to `(ArtifactList, func(context.Context, ArtifactList) (ArtifactList, error))`
- Function `NewAssignmentOperationListPage` parameter(s) have been changed from `(func(context.Context, AssignmentOperationList) (AssignmentOperationList, error))` to `(AssignmentOperationList, func(context.Context, AssignmentOperationList) (AssignmentOperationList, error))`
- Function `NewPublishedBlueprintListPage` parameter(s) have been changed from `(func(context.Context, PublishedBlueprintList) (PublishedBlueprintList, error))` to `(PublishedBlueprintList, func(context.Context, PublishedBlueprintList) (PublishedBlueprintList, error))`
