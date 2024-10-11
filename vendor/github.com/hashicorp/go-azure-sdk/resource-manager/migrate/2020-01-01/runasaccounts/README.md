
## `github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/runasaccounts` Documentation

The `runasaccounts` SDK allows for interaction with Azure Resource Manager `migrate` (API Version `2020-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/runasaccounts"
```


### Client Initialization

```go
client := runasaccounts.NewRunAsAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RunAsAccountsClient.GetAllRunAsAccountsInSite`

```go
ctx := context.TODO()
id := runasaccounts.NewVMwareSiteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmwareSiteName")

// alternatively `client.GetAllRunAsAccountsInSite(ctx, id)` can be used to do batched pagination
items, err := client.GetAllRunAsAccountsInSiteComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RunAsAccountsClient.GetRunAsAccount`

```go
ctx := context.TODO()
id := commonids.NewVMwareSiteRunAsAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmwareSiteName", "runAsAccountName")

read, err := client.GetRunAsAccount(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
