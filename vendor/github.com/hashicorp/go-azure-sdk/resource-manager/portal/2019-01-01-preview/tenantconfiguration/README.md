
## `github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/tenantconfiguration` Documentation

The `tenantconfiguration` SDK allows for interaction with the Azure Resource Manager Service `portal` (API Version `2019-01-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/tenantconfiguration"
```


### Client Initialization

```go
client := tenantconfiguration.NewTenantConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TenantConfigurationClient.TenantConfigurationsCreate`

```go
ctx := context.TODO()

payload := tenantconfiguration.Configuration{
	// ...
}


read, err := client.TenantConfigurationsCreate(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantConfigurationClient.TenantConfigurationsDelete`

```go
ctx := context.TODO()


read, err := client.TenantConfigurationsDelete(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantConfigurationClient.TenantConfigurationsGet`

```go
ctx := context.TODO()


read, err := client.TenantConfigurationsGet(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TenantConfigurationClient.TenantConfigurationsList`

```go
ctx := context.TODO()


read, err := client.TenantConfigurationsList(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
