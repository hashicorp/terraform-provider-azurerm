
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/verifiedpartners` Documentation

The `verifiedpartners` SDK allows for interaction with Azure Resource Manager `eventgrid` (API Version `2022-06-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/verifiedpartners"
```


### Client Initialization

```go
client := verifiedpartners.NewVerifiedPartnersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VerifiedPartnersClient.Get`

```go
ctx := context.TODO()
id := verifiedpartners.NewVerifiedPartnerID("verifiedPartnerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VerifiedPartnersClient.List`

```go
ctx := context.TODO()


// alternatively `client.List(ctx, verifiedpartners.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, verifiedpartners.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
