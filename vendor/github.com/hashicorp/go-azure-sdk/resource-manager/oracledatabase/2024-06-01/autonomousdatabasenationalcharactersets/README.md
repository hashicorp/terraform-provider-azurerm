
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasenationalcharactersets` Documentation

The `autonomousdatabasenationalcharactersets` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasenationalcharactersets"
```


### Client Initialization

```go
client := autonomousdatabasenationalcharactersets.NewAutonomousDatabaseNationalCharacterSetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutonomousDatabaseNationalCharacterSetsClient.Get`

```go
ctx := context.TODO()
id := autonomousdatabasenationalcharactersets.NewAutonomousDatabaseNationalCharacterSetID("12345678-1234-9876-4563-123456789012", "locationName", "autonomousDatabaseNationalCharacterSetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutonomousDatabaseNationalCharacterSetsClient.ListByLocation`

```go
ctx := context.TODO()
id := autonomousdatabasenationalcharactersets.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
