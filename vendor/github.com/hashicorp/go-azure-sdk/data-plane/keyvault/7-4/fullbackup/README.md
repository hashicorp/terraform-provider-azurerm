
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/fullbackup` Documentation

The `fullbackup` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `7.4`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/fullbackup"
```


### Client Initialization

```go
client := fullbackup.NewFullBackupClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `FullBackupClient.FullBackup`

```go
ctx := context.TODO()

payload := fullbackup.SASTokenParameter{
	// ...
}


if err := client.FullBackupThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FullBackupClient.Status`

```go
ctx := context.TODO()
id := fullbackup.NewBackupID("jobId")

read, err := client.Status(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
