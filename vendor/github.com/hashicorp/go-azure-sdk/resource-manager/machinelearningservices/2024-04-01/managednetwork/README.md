
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/managednetwork` Documentation

The `managednetwork` SDK allows for interaction with Azure Resource Manager `machinelearningservices` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/managednetwork"
```


### Client Initialization

```go
client := managednetwork.NewManagedNetworkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedNetworkClient.ProvisionsProvisionManagedNetwork`

```go
ctx := context.TODO()
id := managednetwork.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

payload := managednetwork.ManagedNetworkProvisionOptions{
	// ...
}


if err := client.ProvisionsProvisionManagedNetworkThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedNetworkClient.SettingsRuleCreateOrUpdate`

```go
ctx := context.TODO()
id := managednetwork.NewOutboundRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "outboundRuleName")

payload := managednetwork.OutboundRuleBasicResource{
	// ...
}


if err := client.SettingsRuleCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedNetworkClient.SettingsRuleDelete`

```go
ctx := context.TODO()
id := managednetwork.NewOutboundRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "outboundRuleName")

if err := client.SettingsRuleDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedNetworkClient.SettingsRuleGet`

```go
ctx := context.TODO()
id := managednetwork.NewOutboundRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "outboundRuleName")

read, err := client.SettingsRuleGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedNetworkClient.SettingsRuleList`

```go
ctx := context.TODO()
id := managednetwork.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.SettingsRuleList(ctx, id)` can be used to do batched pagination
items, err := client.SettingsRuleListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
