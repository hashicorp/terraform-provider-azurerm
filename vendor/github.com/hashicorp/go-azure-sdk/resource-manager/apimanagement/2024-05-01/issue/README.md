
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/issue` Documentation

The `issue` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/issue"
```


### Client Initialization

```go
client := issue.NewIssueClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IssueClient.Get`

```go
ctx := context.TODO()
id := issue.NewIssueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "issueId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IssueClient.ListByService`

```go
ctx := context.TODO()
id := issue.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, issue.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, issue.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
