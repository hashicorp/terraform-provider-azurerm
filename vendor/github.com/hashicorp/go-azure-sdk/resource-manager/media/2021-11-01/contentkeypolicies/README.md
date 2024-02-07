
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/contentkeypolicies` Documentation

The `contentkeypolicies` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2021-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/contentkeypolicies"
```


### Client Initialization

```go
client := contentkeypolicies.NewContentKeyPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContentKeyPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

payload := contentkeypolicies.ContentKeyPolicy{
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


### Example Usage: `ContentKeyPoliciesClient.Delete`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentKeyPoliciesClient.Get`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentKeyPoliciesClient.GetPolicyPropertiesWithSecrets`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

read, err := client.GetPolicyPropertiesWithSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentKeyPoliciesClient.List`

```go
ctx := context.TODO()
id := contentkeypolicies.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

// alternatively `client.List(ctx, id, contentkeypolicies.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, contentkeypolicies.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContentKeyPoliciesClient.Update`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

payload := contentkeypolicies.ContentKeyPolicy{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
