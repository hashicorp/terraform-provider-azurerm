
## `github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-11-15/mongorbacs` Documentation

The `mongorbacs` SDK allows for interaction with Azure Resource Manager `cosmosdb` (API Version `2022-11-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-11-15/mongorbacs"
```


### Client Initialization

```go
client := mongorbacs.NewMongorbacsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MongorbacsClient.MongoDBResourcesCreateUpdateMongoRoleDefinition`

```go
ctx := context.TODO()
id := mongorbacs.NewMongodbRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoRoleDefinitionId")

payload := mongorbacs.MongoRoleDefinitionCreateUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesCreateUpdateMongoRoleDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MongorbacsClient.MongoDBResourcesCreateUpdateMongoUserDefinition`

```go
ctx := context.TODO()
id := mongorbacs.NewMongodbUserDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoUserDefinitionId")

payload := mongorbacs.MongoUserDefinitionCreateUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesCreateUpdateMongoUserDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MongorbacsClient.MongoDBResourcesDeleteMongoRoleDefinition`

```go
ctx := context.TODO()
id := mongorbacs.NewMongodbRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoRoleDefinitionId")

if err := client.MongoDBResourcesDeleteMongoRoleDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MongorbacsClient.MongoDBResourcesDeleteMongoUserDefinition`

```go
ctx := context.TODO()
id := mongorbacs.NewMongodbUserDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoUserDefinitionId")

if err := client.MongoDBResourcesDeleteMongoUserDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MongorbacsClient.MongoDBResourcesGetMongoRoleDefinition`

```go
ctx := context.TODO()
id := mongorbacs.NewMongodbRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoRoleDefinitionId")

read, err := client.MongoDBResourcesGetMongoRoleDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MongorbacsClient.MongoDBResourcesGetMongoUserDefinition`

```go
ctx := context.TODO()
id := mongorbacs.NewMongodbUserDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoUserDefinitionId")

read, err := client.MongoDBResourcesGetMongoUserDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MongorbacsClient.MongoDBResourcesListMongoRoleDefinitions`

```go
ctx := context.TODO()
id := mongorbacs.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

read, err := client.MongoDBResourcesListMongoRoleDefinitions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MongorbacsClient.MongoDBResourcesListMongoUserDefinitions`

```go
ctx := context.TODO()
id := mongorbacs.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

read, err := client.MongoDBResourcesListMongoUserDefinitions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
