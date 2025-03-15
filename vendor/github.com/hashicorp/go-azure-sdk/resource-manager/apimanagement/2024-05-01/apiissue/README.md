
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiissue` Documentation

The `apiissue` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiissue"
```


### Client Initialization

```go
client := apiissue.NewApiIssueClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiIssueClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apiissue.NewApiIssueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId")

payload := apiissue.IssueContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apiissue.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiIssueClient.Delete`

```go
ctx := context.TODO()
id := apiissue.NewApiIssueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId")

read, err := client.Delete(ctx, id, apiissue.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiIssueClient.Get`

```go
ctx := context.TODO()
id := apiissue.NewApiIssueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId")

read, err := client.Get(ctx, id, apiissue.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiIssueClient.GetEntityTag`

```go
ctx := context.TODO()
id := apiissue.NewApiIssueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiIssueClient.ListByService`

```go
ctx := context.TODO()
id := apiissue.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

// alternatively `client.ListByService(ctx, id, apiissue.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, apiissue.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiIssueClient.Update`

```go
ctx := context.TODO()
id := apiissue.NewApiIssueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId")

payload := apiissue.IssueUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, apiissue.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
