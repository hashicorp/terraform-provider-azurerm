
## `github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors` Documentation

The `frontdoors` SDK allows for interaction with Azure Resource Manager `frontdoor` (API Version `2020-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
```


### Client Initialization

```go
client := frontdoors.NewFrontDoorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FrontDoorsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := frontdoors.NewFrontDoorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName")

payload := frontdoors.FrontDoor{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FrontDoorsClient.Delete`

```go
ctx := context.TODO()
id := frontdoors.NewFrontDoorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FrontDoorsClient.EndpointsPurgeContent`

```go
ctx := context.TODO()
id := frontdoors.NewFrontDoorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName")

payload := frontdoors.PurgeParameters{
	// ...
}


if err := client.EndpointsPurgeContentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FrontDoorsClient.FrontendEndpointsDisableHTTPS`

```go
ctx := context.TODO()
id := frontdoors.NewFrontendEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName", "frontendEndpointName")

if err := client.FrontendEndpointsDisableHTTPSThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FrontDoorsClient.FrontendEndpointsEnableHTTPS`

```go
ctx := context.TODO()
id := frontdoors.NewFrontendEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName", "frontendEndpointName")

payload := frontdoors.CustomHTTPSConfiguration{
	// ...
}


if err := client.FrontendEndpointsEnableHTTPSThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FrontDoorsClient.FrontendEndpointsGet`

```go
ctx := context.TODO()
id := frontdoors.NewFrontendEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName", "frontendEndpointName")

read, err := client.FrontendEndpointsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FrontDoorsClient.FrontendEndpointsListByFrontDoor`

```go
ctx := context.TODO()
id := frontdoors.NewFrontDoorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName")

// alternatively `client.FrontendEndpointsListByFrontDoor(ctx, id)` can be used to do batched pagination
items, err := client.FrontendEndpointsListByFrontDoorComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FrontDoorsClient.Get`

```go
ctx := context.TODO()
id := frontdoors.NewFrontDoorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FrontDoorsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FrontDoorsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FrontDoorsClient.RulesEnginesCreateOrUpdate`

```go
ctx := context.TODO()
id := frontdoors.NewRulesEngineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName", "rulesEngineName")

payload := frontdoors.RulesEngine{
	// ...
}


if err := client.RulesEnginesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FrontDoorsClient.RulesEnginesDelete`

```go
ctx := context.TODO()
id := frontdoors.NewRulesEngineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName", "rulesEngineName")

if err := client.RulesEnginesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FrontDoorsClient.RulesEnginesGet`

```go
ctx := context.TODO()
id := frontdoors.NewRulesEngineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName", "rulesEngineName")

read, err := client.RulesEnginesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FrontDoorsClient.RulesEnginesListByFrontDoor`

```go
ctx := context.TODO()
id := frontdoors.NewFrontDoorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName")

// alternatively `client.RulesEnginesListByFrontDoor(ctx, id)` can be used to do batched pagination
items, err := client.RulesEnginesListByFrontDoorComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FrontDoorsClient.ValidateCustomDomain`

```go
ctx := context.TODO()
id := frontdoors.NewFrontDoorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "frontDoorName")

payload := frontdoors.ValidateCustomDomainInput{
	// ...
}


read, err := client.ValidateCustomDomain(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
