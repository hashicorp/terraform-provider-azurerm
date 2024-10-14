
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabaseversions` Documentation

The `autonomousdatabaseversions` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabaseversions"
```


### Client Initialization

```go
client := autonomousdatabaseversions.NewAutonomousDatabaseVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutonomousDatabaseVersionsClient.Get`

```go
ctx := context.TODO()
id := autonomousdatabaseversions.NewAutonomousDbVersionID("12345678-1234-9876-4563-123456789012", "locationName", "autonomousDbVersionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutonomousDatabaseVersionsClient.ListByLocation`

```go
ctx := context.TODO()
id := autonomousdatabaseversions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
