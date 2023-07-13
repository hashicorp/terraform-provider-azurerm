
## `github.com/hashicorp/go-azure-sdk/resource-manager/insights/2016-03-01/logprofiles` Documentation

The `logprofiles` SDK allows for interaction with the Azure Resource Manager Service `insights` (API Version `2016-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2016-03-01/logprofiles"
```


### Client Initialization

```go
client := logprofiles.NewLogProfilesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LogProfilesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := logprofiles.NewLogProfileID("12345678-1234-9876-4563-123456789012", "logProfileValue")

payload := logprofiles.LogProfileResource{
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


### Example Usage: `LogProfilesClient.Delete`

```go
ctx := context.TODO()
id := logprofiles.NewLogProfileID("12345678-1234-9876-4563-123456789012", "logProfileValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LogProfilesClient.Get`

```go
ctx := context.TODO()
id := logprofiles.NewLogProfileID("12345678-1234-9876-4563-123456789012", "logProfileValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LogProfilesClient.List`

```go
ctx := context.TODO()
id := logprofiles.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
