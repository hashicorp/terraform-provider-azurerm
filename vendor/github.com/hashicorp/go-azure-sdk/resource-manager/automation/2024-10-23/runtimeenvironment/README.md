
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment` Documentation

The `runtimeenvironment` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2024-10-23`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
```


### Client Initialization

```go
client := runtimeenvironment.NewRuntimeEnvironmentClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RuntimeEnvironmentClient.Create`

```go
ctx := context.TODO()
id := runtimeenvironment.NewRuntimeEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runtimeEnvironmentName")

payload := runtimeenvironment.RuntimeEnvironment{
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


### Example Usage: `RuntimeEnvironmentClient.Delete`

```go
ctx := context.TODO()
id := runtimeenvironment.NewRuntimeEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runtimeEnvironmentName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RuntimeEnvironmentClient.Get`

```go
ctx := context.TODO()
id := runtimeenvironment.NewRuntimeEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runtimeEnvironmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RuntimeEnvironmentClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := runtimeenvironment.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RuntimeEnvironmentClient.Update`

```go
ctx := context.TODO()
id := runtimeenvironment.NewRuntimeEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "runtimeEnvironmentName")

payload := runtimeenvironment.RuntimeEnvironmentUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
