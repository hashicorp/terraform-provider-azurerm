
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettingscategories` Documentation

The `diagnosticsettingscategories` SDK allows for interaction with Azure Resource Manager `insights` (API Version `2021-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettingscategories"
```


### Client Initialization

```go
client := diagnosticsettingscategories.NewDiagnosticSettingsCategoriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DiagnosticSettingsCategoriesClient.DiagnosticSettingsCategoryGet`

```go
ctx := context.TODO()
id := diagnosticsettingscategories.NewScopedDiagnosticSettingsCategoryID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "diagnosticSettingsCategoryName")

read, err := client.DiagnosticSettingsCategoryGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticSettingsCategoriesClient.DiagnosticSettingsCategoryList`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.DiagnosticSettingsCategoryList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
