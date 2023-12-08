
## `github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings` Documentation

The `diagnosticsettings` SDK allows for interaction with the Azure Resource Manager Service `azureactivedirectory` (API Version `2017-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
```


### Client Initialization

```go
client := diagnosticsettings.NewDiagnosticSettingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DiagnosticSettingsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := diagnosticsettings.NewDiagnosticSettingID("diagnosticSettingValue")

payload := diagnosticsettings.DiagnosticSettingsResource{
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


### Example Usage: `DiagnosticSettingsClient.Delete`

```go
ctx := context.TODO()
id := diagnosticsettings.NewDiagnosticSettingID("diagnosticSettingValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticSettingsClient.Get`

```go
ctx := context.TODO()
id := diagnosticsettings.NewDiagnosticSettingID("diagnosticSettingValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticSettingsClient.List`

```go
ctx := context.TODO()


read, err := client.List(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
