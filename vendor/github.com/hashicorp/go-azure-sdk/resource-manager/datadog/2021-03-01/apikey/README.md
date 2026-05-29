
## `github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/apikey` Documentation

The `apikey` SDK allows for interaction with Azure Resource Manager `datadog` (API Version `2021-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/apikey"
```


### Client Initialization

```go
client := apikey.NewApiKeyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiKeyClient.MonitorsGetDefaultKey`

```go
ctx := context.TODO()
id := apikey.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsGetDefaultKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiKeyClient.MonitorsListApiKeys`

```go
ctx := context.TODO()
id := apikey.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.MonitorsListApiKeys(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListApiKeysComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiKeyClient.MonitorsSetDefaultKey`

```go
ctx := context.TODO()
id := apikey.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := apikey.DatadogApiKey{
	// ...
}


read, err := client.MonitorsSetDefaultKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
