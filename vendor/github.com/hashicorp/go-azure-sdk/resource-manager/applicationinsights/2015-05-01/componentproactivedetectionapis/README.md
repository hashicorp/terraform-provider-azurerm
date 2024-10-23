
## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis` Documentation

The `componentproactivedetectionapis` SDK allows for interaction with Azure Resource Manager `applicationinsights` (API Version `2015-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
```


### Client Initialization

```go
client := componentproactivedetectionapis.NewComponentProactiveDetectionAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ComponentProactiveDetectionAPIsClient.ProactiveDetectionConfigurationsGet`

```go
ctx := context.TODO()
id := componentproactivedetectionapis.NewProactiveDetectionConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName", "configurationId")

read, err := client.ProactiveDetectionConfigurationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentProactiveDetectionAPIsClient.ProactiveDetectionConfigurationsList`

```go
ctx := context.TODO()
id := componentproactivedetectionapis.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName")

read, err := client.ProactiveDetectionConfigurationsList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentProactiveDetectionAPIsClient.ProactiveDetectionConfigurationsUpdate`

```go
ctx := context.TODO()
id := componentproactivedetectionapis.NewProactiveDetectionConfigID("12345678-1234-9876-4563-123456789012", "example-resource-group", "componentName", "configurationId")

payload := componentproactivedetectionapis.ApplicationInsightsComponentProactiveDetectionConfiguration{
	// ...
}


read, err := client.ProactiveDetectionConfigurationsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
