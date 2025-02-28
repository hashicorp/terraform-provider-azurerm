
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiissueattachment` Documentation

The `apiissueattachment` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiissueattachment"
```


### Client Initialization

```go
client := apiissueattachment.NewApiIssueAttachmentClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiIssueAttachmentClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apiissueattachment.NewAttachmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId", "attachmentId")

payload := apiissueattachment.IssueAttachmentContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apiissueattachment.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiIssueAttachmentClient.Delete`

```go
ctx := context.TODO()
id := apiissueattachment.NewAttachmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId", "attachmentId")

read, err := client.Delete(ctx, id, apiissueattachment.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiIssueAttachmentClient.Get`

```go
ctx := context.TODO()
id := apiissueattachment.NewAttachmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId", "attachmentId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiIssueAttachmentClient.GetEntityTag`

```go
ctx := context.TODO()
id := apiissueattachment.NewAttachmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId", "attachmentId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiIssueAttachmentClient.ListByService`

```go
ctx := context.TODO()
id := apiissueattachment.NewApiIssueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "issueId")

// alternatively `client.ListByService(ctx, id, apiissueattachment.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, apiissueattachment.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
