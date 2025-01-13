
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapsupportedsku` Documentation

The `sapsupportedsku` SDK allows for interaction with Azure Resource Manager `workloads` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapsupportedsku"
```


### Client Initialization

```go
client := sapsupportedsku.NewSAPSupportedSkuClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SAPSupportedSkuClient.SAPSupportedSku`

```go
ctx := context.TODO()
id := sapsupportedsku.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := sapsupportedsku.SAPSupportedSkusRequest{
	// ...
}


read, err := client.SAPSupportedSku(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
