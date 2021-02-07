Generated from https://github.com/Azure/azure-rest-api-specs/tree/3a3a9452f965a227ce43e6b545035b99dd175f23

Code generator @microsoft.azure/autorest.go@~2.1.165

## Breaking Changes

## Signature Changes

### Funcs

1. WorkspaceManagedIdentitySQLControlSettingsClient.CreateOrUpdate
	- Returns
		- From: ManagedIdentitySQLControlSettingsModel, error
		- To: WorkspaceManagedIdentitySQLControlSettingsCreateOrUpdateFuture, error
1. WorkspaceManagedIdentitySQLControlSettingsClient.CreateOrUpdateSender
	- Returns
		- From: *http.Response, error
		- To: WorkspaceManagedIdentitySQLControlSettingsCreateOrUpdateFuture, error

### New Funcs

1. *WorkspaceManagedIdentitySQLControlSettingsCreateOrUpdateFuture.Result(WorkspaceManagedIdentitySQLControlSettingsClient) (ManagedIdentitySQLControlSettingsModel, error)

## Struct Changes

### New Structs

1. WorkspaceManagedIdentitySQLControlSettingsCreateOrUpdateFuture

### New Struct Fields

1. WorkspacePatchProperties.Encryption
