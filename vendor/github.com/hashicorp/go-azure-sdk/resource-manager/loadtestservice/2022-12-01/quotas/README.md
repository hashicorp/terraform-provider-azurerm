
## `github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/quotas` Documentation

The `quotas` SDK allows for interaction with Azure Resource Manager `loadtestservice` (API Version `2022-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/quotas"
```


### Client Initialization

```go
client := quotas.NewQuotasClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `QuotasClient.CheckAvailability`

```go
ctx := context.TODO()
id := quotas.NewQuotaID("12345678-1234-9876-4563-123456789012", "locationName", "quotaName")

payload := quotas.QuotaBucketRequest{
	// ...
}


read, err := client.CheckAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QuotasClient.Get`

```go
ctx := context.TODO()
id := quotas.NewQuotaID("12345678-1234-9876-4563-123456789012", "locationName", "quotaName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QuotasClient.List`

```go
ctx := context.TODO()
id := quotas.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
