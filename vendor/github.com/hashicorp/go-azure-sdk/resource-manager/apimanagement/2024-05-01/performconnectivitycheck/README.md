
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/performconnectivitycheck` Documentation

The `performconnectivitycheck` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/performconnectivitycheck"
```


### Client Initialization

```go
client := performconnectivitycheck.NewPerformConnectivityCheckClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PerformConnectivityCheckClient.Async`

```go
ctx := context.TODO()
id := performconnectivitycheck.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

payload := performconnectivitycheck.ConnectivityCheckRequest{
	// ...
}


if err := client.AsyncThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
