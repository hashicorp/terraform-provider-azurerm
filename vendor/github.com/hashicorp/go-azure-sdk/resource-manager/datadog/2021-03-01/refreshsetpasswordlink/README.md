
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/refreshsetpasswordlink` Documentation

The `refreshsetpasswordlink` SDK allows for interaction with Azure Resource Manager `datadog` (API Version `2021-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/refreshsetpasswordlink"
```


### Client Initialization

```go
client := refreshsetpasswordlink.NewRefreshSetPasswordLinkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RefreshSetPasswordLinkClient.MonitorsRefreshSetPasswordLink`

```go
ctx := context.TODO()
id := refreshsetpasswordlink.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsRefreshSetPasswordLink(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
