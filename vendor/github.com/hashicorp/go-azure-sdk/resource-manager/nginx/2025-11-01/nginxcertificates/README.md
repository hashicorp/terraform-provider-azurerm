
## `github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxcertificates` Documentation

The `nginxcertificates` SDK allows for interaction with Azure Resource Manager `nginx` (API Version `2025-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2025-11-01/nginxcertificates"
```


### Client Initialization

```go
client := nginxcertificates.NewNginxCertificatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NginxCertificatesClient.CertificatesCreateOrUpdate`

```go
ctx := context.TODO()
id := nginxcertificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "certificateName")

payload := nginxcertificates.NginxCertificate{
	// ...
}


if err := client.CertificatesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NginxCertificatesClient.CertificatesDelete`

```go
ctx := context.TODO()
id := nginxcertificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "certificateName")

if err := client.CertificatesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NginxCertificatesClient.CertificatesGet`

```go
ctx := context.TODO()
id := nginxcertificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName", "certificateName")

read, err := client.CertificatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxCertificatesClient.CertificatesList`

```go
ctx := context.TODO()
id := nginxcertificates.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentName")

// alternatively `client.CertificatesList(ctx, id)` can be used to do batched pagination
items, err := client.CertificatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
