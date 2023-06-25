
## `github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway` Documentation

The `sqldedicatedgateway` SDK allows for interaction with the Azure Resource Manager Service `cosmosdb` (API Version `2022-05-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
```


### Client Initialization

```go
client := sqldedicatedgateway.NewSqlDedicatedGatewayClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SqlDedicatedGatewayClient.ServiceCreate`

```go
ctx := context.TODO()
id := sqldedicatedgateway.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "serviceValue")

payload := sqldedicatedgateway.ServiceResourceCreateUpdateParameters{
	// ...
}


if err := client.ServiceCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SqlDedicatedGatewayClient.ServiceDelete`

```go
ctx := context.TODO()
id := sqldedicatedgateway.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "serviceValue")

if err := client.ServiceDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SqlDedicatedGatewayClient.ServiceGet`

```go
ctx := context.TODO()
id := sqldedicatedgateway.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountValue", "serviceValue")

read, err := client.ServiceGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
