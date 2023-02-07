
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/contentkeypolicies` Documentation

The `contentkeypolicies` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/contentkeypolicies"
```


### Client Initialization

```go
client := contentkeypolicies.NewContentKeyPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContentKeyPoliciesClient.ContentKeyPoliciesCreateOrUpdate`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

payload := contentkeypolicies.ContentKeyPolicy{
	// ...
}


read, err := client.ContentKeyPoliciesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentKeyPoliciesClient.ContentKeyPoliciesDelete`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

read, err := client.ContentKeyPoliciesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentKeyPoliciesClient.ContentKeyPoliciesGet`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

read, err := client.ContentKeyPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentKeyPoliciesClient.ContentKeyPoliciesGetPolicyPropertiesWithSecrets`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

read, err := client.ContentKeyPoliciesGetPolicyPropertiesWithSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentKeyPoliciesClient.ContentKeyPoliciesList`

```go
ctx := context.TODO()
id := contentkeypolicies.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

// alternatively `client.ContentKeyPoliciesList(ctx, id, contentkeypolicies.DefaultContentKeyPoliciesListOperationOptions())` can be used to do batched pagination
items, err := client.ContentKeyPoliciesListComplete(ctx, id, contentkeypolicies.DefaultContentKeyPoliciesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContentKeyPoliciesClient.ContentKeyPoliciesUpdate`

```go
ctx := context.TODO()
id := contentkeypolicies.NewContentKeyPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "contentKeyPolicyValue")

payload := contentkeypolicies.ContentKeyPolicy{
	// ...
}


read, err := client.ContentKeyPoliciesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
