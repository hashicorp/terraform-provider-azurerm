
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/serversecurityalertpolicies` Documentation

The `serversecurityalertpolicies` SDK allows for interaction with the Azure Resource Manager Service `mysql` (API Version `2017-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/serversecurityalertpolicies"
```


### Client Initialization

```go
client := serversecurityalertpolicies.NewServerSecurityAlertPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServerSecurityAlertPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := serversecurityalertpolicies.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

payload := serversecurityalertpolicies.ServerSecurityAlertPolicy{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServerSecurityAlertPoliciesClient.Get`

```go
ctx := context.TODO()
id := serversecurityalertpolicies.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServerSecurityAlertPoliciesClient.ListByServer`

```go
ctx := context.TODO()
id := serversecurityalertpolicies.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
