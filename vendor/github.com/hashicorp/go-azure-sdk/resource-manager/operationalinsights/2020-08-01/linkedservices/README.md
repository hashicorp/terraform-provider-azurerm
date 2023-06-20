
## `github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices` Documentation

The `linkedservices` SDK allows for interaction with the Azure Resource Manager Service `operationalinsights` (API Version `2020-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
```


### Client Initialization

```go
client := linkedservices.NewLinkedServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LinkedServicesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "linkedServiceValue")

payload := linkedservices.LinkedService{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LinkedServicesClient.Delete`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "linkedServiceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LinkedServicesClient.Get`

```go
ctx := context.TODO()
id := linkedservices.NewLinkedServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "linkedServiceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LinkedServicesClient.ListByWorkspace`

```go
ctx := context.TODO()
id := linkedservices.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

read, err := client.ListByWorkspace(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
