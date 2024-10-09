
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/oraclesubscriptions` Documentation

The `oraclesubscriptions` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/oraclesubscriptions"
```


### Client Initialization

```go
client := oraclesubscriptions.NewOracleSubscriptionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OracleSubscriptionsClient.AddAzureSubscriptions`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := oraclesubscriptions.AzureSubscriptions{
	// ...
}


if err := client.AddAzureSubscriptionsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OracleSubscriptionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := oraclesubscriptions.OracleSubscription{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OracleSubscriptionsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OracleSubscriptionsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OracleSubscriptionsClient.ListActivationLinks`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

if err := client.ListActivationLinksThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OracleSubscriptionsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OracleSubscriptionsClient.ListCloudAccountDetails`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

if err := client.ListCloudAccountDetailsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OracleSubscriptionsClient.ListSaasSubscriptionDetails`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

if err := client.ListSaasSubscriptionDetailsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OracleSubscriptionsClient.Update`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := oraclesubscriptions.OracleSubscriptionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
