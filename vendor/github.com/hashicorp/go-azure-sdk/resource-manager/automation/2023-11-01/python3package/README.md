
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/python3package` Documentation

The `python3package` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/python3package"
```


### Client Initialization

```go
client := python3package.NewPython3PackageClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `Python3PackageClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := python3package.NewPython3PackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "python3PackageName")

payload := python3package.PythonPackageCreateParameters{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `Python3PackageClient.Delete`

```go
ctx := context.TODO()
id := python3package.NewPython3PackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "python3PackageName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `Python3PackageClient.Get`

```go
ctx := context.TODO()
id := python3package.NewPython3PackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "python3PackageName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `Python3PackageClient.ListByAutomationAccount`

```go
ctx := context.TODO()
id := python3package.NewAutomationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName")

// alternatively `client.ListByAutomationAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByAutomationAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `Python3PackageClient.Update`

```go
ctx := context.TODO()
id := python3package.NewPython3PackageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "python3PackageName")

payload := python3package.PythonPackageUpdateParameters{
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
