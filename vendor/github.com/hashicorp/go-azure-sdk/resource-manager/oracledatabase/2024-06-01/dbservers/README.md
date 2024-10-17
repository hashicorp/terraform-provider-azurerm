
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbservers` Documentation

The `dbservers` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbservers"
```


### Client Initialization

```go
client := dbservers.NewDbServersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DbServersClient.Get`

```go
ctx := context.TODO()
id := dbservers.NewDbServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudExadataInfrastructureName", "dbServerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DbServersClient.ListByCloudExadataInfrastructure`

```go
ctx := context.TODO()
id := dbservers.NewCloudExadataInfrastructureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudExadataInfrastructureName")

// alternatively `client.ListByCloudExadataInfrastructure(ctx, id)` can be used to do batched pagination
items, err := client.ListByCloudExadataInfrastructureComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
