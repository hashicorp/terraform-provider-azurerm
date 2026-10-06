
## `github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules` Documentation

The `forwardingrules` SDK allows for interaction with Azure Resource Manager `dnsresolver` (API Version `2022-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
```


### Client Initialization

```go
client := forwardingrules.NewForwardingRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ForwardingRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := forwardingrules.NewForwardingRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName", "forwardingRuleName")

payload := forwardingrules.ForwardingRule{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, forwardingrules.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ForwardingRulesClient.Delete`

```go
ctx := context.TODO()
id := forwardingrules.NewForwardingRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName", "forwardingRuleName")

read, err := client.Delete(ctx, id, forwardingrules.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ForwardingRulesClient.Get`

```go
ctx := context.TODO()
id := forwardingrules.NewForwardingRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName", "forwardingRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ForwardingRulesClient.List`

```go
ctx := context.TODO()
id := forwardingrules.NewDnsForwardingRulesetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName")

// alternatively `client.List(ctx, id, forwardingrules.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, forwardingrules.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ForwardingRulesClient.Update`

```go
ctx := context.TODO()
id := forwardingrules.NewForwardingRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName", "forwardingRuleName")

payload := forwardingrules.ForwardingRulePatch{
	// ...
}


read, err := client.Update(ctx, id, payload, forwardingrules.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
