
## `github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotassubscriptions` Documentation

The `groupquotassubscriptions` SDK allows for interaction with Azure Resource Manager `quota` (API Version `2025-07-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/quota/2025-07-15/groupquotassubscriptions"
```


### Client Initialization

```go
client := groupquotassubscriptions.NewGroupQuotasSubscriptionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GroupQuotasSubscriptionsClient.GroupQuotaSubscriptionsCreateOrUpdate`

```go
ctx := context.TODO()
id := groupquotassubscriptions.NewSubscriptionID("managementGroupId", "groupQuotaName", "12345678-1234-9876-4563-123456789012")

if err := client.GroupQuotaSubscriptionsCreateOrUpdateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GroupQuotasSubscriptionsClient.GroupQuotaSubscriptionsDelete`

```go
ctx := context.TODO()
id := groupquotassubscriptions.NewSubscriptionID("managementGroupId", "groupQuotaName", "12345678-1234-9876-4563-123456789012")

if err := client.GroupQuotaSubscriptionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GroupQuotasSubscriptionsClient.GroupQuotaSubscriptionsGet`

```go
ctx := context.TODO()
id := groupquotassubscriptions.NewSubscriptionID("managementGroupId", "groupQuotaName", "12345678-1234-9876-4563-123456789012")

read, err := client.GroupQuotaSubscriptionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupQuotasSubscriptionsClient.GroupQuotaSubscriptionsList`

```go
ctx := context.TODO()
id := groupquotassubscriptions.NewGroupQuotaID("managementGroupId", "groupQuotaName")

// alternatively `client.GroupQuotaSubscriptionsList(ctx, id)` can be used to do batched pagination
items, err := client.GroupQuotaSubscriptionsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GroupQuotasSubscriptionsClient.GroupQuotaSubscriptionsUpdate`

```go
ctx := context.TODO()
id := groupquotassubscriptions.NewSubscriptionID("managementGroupId", "groupQuotaName", "12345678-1234-9876-4563-123456789012")

if err := client.GroupQuotaSubscriptionsUpdateThenPoll(ctx, id); err != nil {
	// handle the error
}
```
