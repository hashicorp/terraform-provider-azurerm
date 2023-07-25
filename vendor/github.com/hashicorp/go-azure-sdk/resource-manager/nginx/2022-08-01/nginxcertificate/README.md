
## `github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxcertificate` Documentation

The `nginxcertificate` SDK allows for interaction with the Azure Resource Manager Service `nginx` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxcertificate"
```


### Client Initialization

```go
client := nginxcertificate.NewNginxCertificateClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NginxCertificateClient.CertificatesCreateOrUpdate`

```go
ctx := context.TODO()
id := nginxcertificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentValue", "certificateValue")

payload := nginxcertificate.NginxCertificate{
	// ...
}


if err := client.CertificatesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NginxCertificateClient.CertificatesDelete`

```go
ctx := context.TODO()
id := nginxcertificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentValue", "certificateValue")

if err := client.CertificatesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NginxCertificateClient.CertificatesGet`

```go
ctx := context.TODO()
id := nginxcertificate.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentValue", "certificateValue")

read, err := client.CertificatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NginxCertificateClient.CertificatesList`

```go
ctx := context.TODO()
id := nginxcertificate.NewNginxDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "nginxDeploymentValue")

// alternatively `client.CertificatesList(ctx, id)` can be used to do batched pagination
items, err := client.CertificatesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
