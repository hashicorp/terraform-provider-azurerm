
## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountassemblies` Documentation

The `integrationaccountassemblies` SDK allows for interaction with the Azure Resource Manager Service `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountassemblies"
```


### Client Initialization

```go
client := integrationaccountassemblies.NewIntegrationAccountAssembliesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationAccountAssembliesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationaccountassemblies.NewAssemblyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "assemblyValue")

payload := integrationaccountassemblies.AssemblyDefinition{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountAssembliesClient.Delete`

```go
ctx := context.TODO()
id := integrationaccountassemblies.NewAssemblyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "assemblyValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountAssembliesClient.Get`

```go
ctx := context.TODO()
id := integrationaccountassemblies.NewAssemblyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "assemblyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountAssembliesClient.List`

```go
ctx := context.TODO()
id := integrationaccountassemblies.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountAssembliesClient.ListContentCallbackUrl`

```go
ctx := context.TODO()
id := integrationaccountassemblies.NewAssemblyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "assemblyValue")

read, err := client.ListContentCallbackUrl(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
