
## `github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/analysisservices` Documentation

The `analysisservices` SDK allows for interaction with the Azure Resource Manager Service `analysisservices` (API Version `2017-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/analysisservices"
```


### Client Initialization

```go
client := analysisservices.NewAnalysisServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AnalysisServicesClient.ServersListSkusForNew`

```go
ctx := context.TODO()
id := analysisservices.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ServersListSkusForNew(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
