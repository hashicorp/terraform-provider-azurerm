
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/advancedthreatprotectionsettings` Documentation

The `advancedthreatprotectionsettings` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2024-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/advancedthreatprotectionsettings"
```


### Client Initialization

```go
client := advancedthreatprotectionsettings.NewAdvancedThreatProtectionSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AdvancedThreatProtectionSettingsClient.Get`

```go
ctx := context.TODO()
id := advancedthreatprotectionsettings.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AdvancedThreatProtectionSettingsClient.List`

```go
ctx := context.TODO()
id := advancedthreatprotectionsettings.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AdvancedThreatProtectionSettingsClient.Update`

```go
ctx := context.TODO()
id := advancedthreatprotectionsettings.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := advancedthreatprotectionsettings.AdvancedThreatProtectionForUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AdvancedThreatProtectionSettingsClient.UpdatePut`

```go
ctx := context.TODO()
id := advancedthreatprotectionsettings.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

payload := advancedthreatprotectionsettings.AdvancedThreatProtection{
	// ...
}


if err := client.UpdatePutThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
