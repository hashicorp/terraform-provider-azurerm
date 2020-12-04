Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewPrivateLinkResourcesPage` parameter(s) have been changed from `(func(context.Context, PrivateLinkResources) (PrivateLinkResources, error))` to `(PrivateLinkResources, func(context.Context, PrivateLinkResources) (PrivateLinkResources, error))`
- Function `NewVaultListPage` parameter(s) have been changed from `(func(context.Context, VaultList) (VaultList, error))` to `(VaultList, func(context.Context, VaultList) (VaultList, error))`
- Function `NewClientDiscoveryResponsePage` parameter(s) have been changed from `(func(context.Context, ClientDiscoveryResponse) (ClientDiscoveryResponse, error))` to `(ClientDiscoveryResponse, func(context.Context, ClientDiscoveryResponse) (ClientDiscoveryResponse, error))`

## New Content

- New field `Identity` in struct `PatchVault`
