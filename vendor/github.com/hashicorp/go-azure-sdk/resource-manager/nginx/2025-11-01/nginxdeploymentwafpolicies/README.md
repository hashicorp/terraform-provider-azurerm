
## `github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeploymentwafpolicies` Documentation

The `nginxdeploymentwafpolicies` SDK allows for interaction with Azure Resource Manager `nginx` (API Version `2025-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxdeploymentwafpolicies"
```


### Client Initialization

```go
client := nginxdeploymentwafpolicies.NewNginxDeploymentWafPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NginxDeploymentWafPoliciesClient.Analysis`

```go
ctx := context.TODO()
id := nginxdeploymentwafpolicies.NewWafPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "wafPolicyName")

payload := nginxdeploymentwafpolicies.NginxDeploymentWafPolicyAnalysisCreateRequest{
	// ...
}


read, err := client.Analysis(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxDeploymentWafPoliciesClient.WafPolicyCreate`

```go
ctx := context.TODO()
id := nginxdeploymentwafpolicies.NewWafPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "wafPolicyName")

payload := nginxdeploymentwafpolicies.NginxDeploymentWafPolicy{
	// ...
}


if err := client.WafPolicyCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NginxDeploymentWafPoliciesClient.WafPolicyDelete`

```go
ctx := context.TODO()
id := nginxdeploymentwafpolicies.NewWafPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "wafPolicyName")

if err := client.WafPolicyDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NginxDeploymentWafPoliciesClient.WafPolicyGet`

```go
ctx := context.TODO()
id := nginxdeploymentwafpolicies.NewWafPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "wafPolicyName")

read, err := client.WafPolicyGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
