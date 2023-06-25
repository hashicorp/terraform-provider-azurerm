
## `github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2020-01-01/getrecommendations` Documentation

The `getrecommendations` SDK allows for interaction with the Azure Resource Manager Service `advisor` (API Version `2020-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2020-01-01/getrecommendations"
```


### Client Initialization

```go
client := getrecommendations.NewGetRecommendationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GetRecommendationsClient.RecommendationsGet`

```go
ctx := context.TODO()
id := getrecommendations.NewScopedRecommendationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "recommendationIdValue")

read, err := client.RecommendationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GetRecommendationsClient.RecommendationsList`

```go
ctx := context.TODO()
id := getrecommendations.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.RecommendationsList(ctx, id, getrecommendations.DefaultRecommendationsListOperationOptions())` can be used to do batched pagination
items, err := client.RecommendationsListComplete(ctx, id, getrecommendations.DefaultRecommendationsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
