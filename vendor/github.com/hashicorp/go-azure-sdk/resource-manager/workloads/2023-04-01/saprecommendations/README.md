
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saprecommendations` Documentation

The `saprecommendations` SDK allows for interaction with the Azure Resource Manager Service `workloads` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saprecommendations"
```


### Client Initialization

```go
client := saprecommendations.NewSAPRecommendationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SAPRecommendationsClient.SAPSizingRecommendations`

```go
ctx := context.TODO()
id := saprecommendations.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := saprecommendations.SAPSizingRecommendationRequest{
	// ...
}


read, err := client.SAPSizingRecommendations(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
