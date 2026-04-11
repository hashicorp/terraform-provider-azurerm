
## `github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/service` Documentation

The `service` SDK allows for interaction with <unknown source data type> `search` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/service"
```


### Client Initialization

```go
client := service.NewServiceClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServiceClient.GetServiceStatistics`

```go
ctx := context.TODO()


read, err := client.GetServiceStatistics(ctx, service.DefaultGetServiceStatisticsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
