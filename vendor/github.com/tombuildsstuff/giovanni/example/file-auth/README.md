## Example: using a Shared Key

This example provisions a Storage Container using a Shared Key for authentication.

To run this example you need the following Environment Variables set:

* `ARM_ENVIRONMENT` - The Azure Environment (`public`, `germany` etc)

You also need to update `main.go` to set the variable `storageAccountName` and `storageAccountKey` to an existing Storage Account (since we don't provision one for you).

Assuming you've got Go installed - you can then run this using:

```bash
$ go run main.go
```