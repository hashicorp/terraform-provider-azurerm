
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/quotabycounterkeys` Documentation

The `quotabycounterkeys` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/quotabycounterkeys"
```


### Client Initialization

```go
client := quotabycounterkeys.NewQuotaByCounterKeysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `QuotaByCounterKeysClient.ListByService`

```go
ctx := context.TODO()
id := quotabycounterkeys.NewQuotaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "quotaCounterKey")

// alternatively `client.ListByService(ctx, id)` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `QuotaByCounterKeysClient.Update`

```go
ctx := context.TODO()
id := quotabycounterkeys.NewQuotaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "quotaCounterKey")

payload := quotabycounterkeys.QuotaCounterValueUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
