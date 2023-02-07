
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/sentinelonboardingstates` Documentation

The `sentinelonboardingstates` SDK allows for interaction with the Azure Resource Manager Service `securityinsights` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/sentinelonboardingstates"
```


### Client Initialization

```go
client := sentinelonboardingstates.NewSentinelOnboardingStatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SentinelOnboardingStatesClient.Create`

```go
ctx := context.TODO()
id := sentinelonboardingstates.NewOnboardingStateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "onboardingStateValue")

payload := sentinelonboardingstates.SentinelOnboardingState{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SentinelOnboardingStatesClient.Delete`

```go
ctx := context.TODO()
id := sentinelonboardingstates.NewOnboardingStateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "onboardingStateValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SentinelOnboardingStatesClient.Get`

```go
ctx := context.TODO()
id := sentinelonboardingstates.NewOnboardingStateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "onboardingStateValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SentinelOnboardingStatesClient.List`

```go
ctx := context.TODO()
id := sentinelonboardingstates.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
