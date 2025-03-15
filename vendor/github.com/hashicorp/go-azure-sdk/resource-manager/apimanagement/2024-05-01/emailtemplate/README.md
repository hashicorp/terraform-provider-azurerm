
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/emailtemplate` Documentation

The `emailtemplate` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/emailtemplate"
```


### Client Initialization

```go
client := emailtemplate.NewEmailTemplateClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EmailTemplateClient.ListByService`

```go
ctx := context.TODO()
id := emailtemplate.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, emailtemplate.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, emailtemplate.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
