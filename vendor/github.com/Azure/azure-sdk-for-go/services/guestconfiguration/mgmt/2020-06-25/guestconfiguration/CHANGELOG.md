Generated from https://github.com/Azure/azure-rest-api-specs/tree/e0f8b9ab0f5fe5e71b7429ebfea8a33c19ec9d8d/specification/guestconfiguration/resource-manager/readme.md tag: `package-2020-06-25`

Code generator @microsoft.azure/autorest.go@2.1.178


## Breaking Changes

### Removed Constants

1. AllowModuleOverwrite.False
1. AllowModuleOverwrite.True
1. RebootIfNeeded.RebootIfNeededFalse
1. RebootIfNeeded.RebootIfNeededTrue

### Removed Funcs

1. PossibleAllowModuleOverwriteValues() []AllowModuleOverwrite
1. PossibleRebootIfNeededValues() []RebootIfNeeded

## Struct Changes

### Removed Structs

1. AssignmentsCreateOrUpdateFuture
1. AssignmentsDeleteFuture
1. HCRPAssignmentsCreateOrUpdateFuture
1. HCRPAssignmentsDeleteFuture

## Signature Changes

### Funcs

1. AssignmentsClient.CreateOrUpdate
	- Returns
		- From: AssignmentsCreateOrUpdateFuture, error
		- To: Assignment, error
1. AssignmentsClient.CreateOrUpdateSender
	- Returns
		- From: AssignmentsCreateOrUpdateFuture, error
		- To: *http.Response, error
1. AssignmentsClient.Delete
	- Returns
		- From: AssignmentsDeleteFuture, error
		- To: autorest.Response, error
1. AssignmentsClient.DeleteSender
	- Returns
		- From: AssignmentsDeleteFuture, error
		- To: *http.Response, error
1. HCRPAssignmentsClient.CreateOrUpdate
	- Returns
		- From: HCRPAssignmentsCreateOrUpdateFuture, error
		- To: Assignment, error
1. HCRPAssignmentsClient.CreateOrUpdateSender
	- Returns
		- From: HCRPAssignmentsCreateOrUpdateFuture, error
		- To: *http.Response, error
1. HCRPAssignmentsClient.Delete
	- Returns
		- From: HCRPAssignmentsDeleteFuture, error
		- To: autorest.Response, error
1. HCRPAssignmentsClient.DeleteSender
	- Returns
		- From: HCRPAssignmentsDeleteFuture, error
		- To: *http.Response, error

### Struct Fields

1. ConfigurationSetting.AllowModuleOverwrite changed type from AllowModuleOverwrite to *bool
1. ConfigurationSetting.RebootIfNeeded changed type from RebootIfNeeded to *bool
