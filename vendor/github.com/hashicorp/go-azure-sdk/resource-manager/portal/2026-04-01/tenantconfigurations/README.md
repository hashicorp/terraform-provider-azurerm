
## `github.com/hashicorp/go-azure-sdk/resource-manager/portal/2026-04-01/tenantconfigurations` Documentation

The `tenantconfigurations` SDK allows for interaction with Azure Resource Manager `portal` (API Version `2026-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/portal/2026-04-01/tenantconfigurations"
```


### Client Initialization

```go
client := tenantconfigurations.NewTenantConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TenantConfigurationsClient.Create`

```go
ctx := context.TODO()
id := tenantconfigurations.NewTenantConfigurationID("tenantConfigurationName")

payload := tenantconfigurations.Configuration{
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


### Example Usage: `TenantConfigurationsClient.Delete`

```go
ctx := context.TODO()
id := tenantconfigurations.NewTenantConfigurationID("tenantConfigurationName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantConfigurationsClient.Get`

```go
ctx := context.TODO()
id := tenantconfigurations.NewTenantConfigurationID("tenantConfigurationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantConfigurationsClient.List`

```go
ctx := context.TODO()


// alternatively `client.List(ctx)` can be used to do batched pagination
items, err := client.ListComplete(ctx)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
