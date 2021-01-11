Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewSessionHostListPage` parameter(s) have been changed from `(func(context.Context, SessionHostList) (SessionHostList, error))` to `(SessionHostList, func(context.Context, SessionHostList) (SessionHostList, error))`
- Function `NewHostPoolListPage` parameter(s) have been changed from `(func(context.Context, HostPoolList) (HostPoolList, error))` to `(HostPoolList, func(context.Context, HostPoolList) (HostPoolList, error))`
- Function `NewApplicationGroupListPage` parameter(s) have been changed from `(func(context.Context, ApplicationGroupList) (ApplicationGroupList, error))` to `(ApplicationGroupList, func(context.Context, ApplicationGroupList) (ApplicationGroupList, error))`
- Function `NewUserSessionListPage` parameter(s) have been changed from `(func(context.Context, UserSessionList) (UserSessionList, error))` to `(UserSessionList, func(context.Context, UserSessionList) (UserSessionList, error))`
- Function `NewWorkspaceListPage` parameter(s) have been changed from `(func(context.Context, WorkspaceList) (WorkspaceList, error))` to `(WorkspaceList, func(context.Context, WorkspaceList) (WorkspaceList, error))`
- Function `NewApplicationListPage` parameter(s) have been changed from `(func(context.Context, ApplicationList) (ApplicationList, error))` to `(ApplicationList, func(context.Context, ApplicationList) (ApplicationList, error))`
- Function `NewStartMenuItemListPage` parameter(s) have been changed from `(func(context.Context, StartMenuItemList) (StartMenuItemList, error))` to `(StartMenuItemList, func(context.Context, StartMenuItemList) (StartMenuItemList, error))`

## New Content

- New field `VMTemplate` in struct `HostPoolPatchProperties`
