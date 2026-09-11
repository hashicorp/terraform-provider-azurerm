
## `github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/threatintelligence` Documentation

The `threatintelligence` SDK allows for interaction with Azure Resource Manager `securityinsights` (API Version `2022-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/threatintelligence"
```


### Client Initialization

```go
client := threatintelligence.NewThreatIntelligenceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ThreatIntelligenceClient.IndicatorAppendTags`

```go
ctx := context.TODO()
id := threatintelligence.NewIndicatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "indicatorName")

payload := threatintelligence.ThreatIntelligenceAppendTags{
	// ...
}


read, err := client.IndicatorAppendTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ThreatIntelligenceClient.IndicatorCreate`

```go
ctx := context.TODO()
id := threatintelligence.NewIndicatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "indicatorName")

payload := threatintelligence.ThreatIntelligenceIndicatorModel{
	// ...
}


read, err := client.IndicatorCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ThreatIntelligenceClient.IndicatorCreateIndicator`

```go
ctx := context.TODO()
id := threatintelligence.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

payload := threatintelligence.ThreatIntelligenceIndicatorModel{
	// ...
}


read, err := client.IndicatorCreateIndicator(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ThreatIntelligenceClient.IndicatorDelete`

```go
ctx := context.TODO()
id := threatintelligence.NewIndicatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "indicatorName")

read, err := client.IndicatorDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ThreatIntelligenceClient.IndicatorGet`

```go
ctx := context.TODO()
id := threatintelligence.NewIndicatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "indicatorName")

read, err := client.IndicatorGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ThreatIntelligenceClient.IndicatorMetricsList`

```go
ctx := context.TODO()
id := threatintelligence.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

read, err := client.IndicatorMetricsList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ThreatIntelligenceClient.IndicatorQueryIndicators`

```go
ctx := context.TODO()
id := threatintelligence.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

payload := threatintelligence.ThreatIntelligenceFilteringCriteria{
	// ...
}


// alternatively `client.IndicatorQueryIndicators(ctx, id, payload)` can be used to do batched pagination
items, err := client.IndicatorQueryIndicatorsComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ThreatIntelligenceClient.IndicatorReplaceTags`

```go
ctx := context.TODO()
id := threatintelligence.NewIndicatorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "indicatorName")

payload := threatintelligence.ThreatIntelligenceIndicatorModel{
	// ...
}


read, err := client.IndicatorReplaceTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ThreatIntelligenceClient.IndicatorsList`

```go
ctx := context.TODO()
id := threatintelligence.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.IndicatorsList(ctx, id, threatintelligence.DefaultIndicatorsListOperationOptions())` can be used to do batched pagination
items, err := client.IndicatorsListComplete(ctx, id, threatintelligence.DefaultIndicatorsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
