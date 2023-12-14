## Example: using Azure Active Directory authentication

This example provisions a Storage Container using Azure Active Directory for authentication.

To run this example you need the following Environment Variables set:

* `ARM_CLIENT_ID` - The UUID of the Service Principal/Application
* `ARM_CLIENT_SECRET` - The Secret associated with the Service Principal
* `ARM_ENVIRONMENT` - The Azure Environment (`public`, `germany` etc)
* `ARM_SUBSCRIPTION_ID` - The UUID of the Azure Subscription
* `ARM_TENANT_ID` - The UUID of the Azure Tenant

You also need to update `main.go` to set the variable `storageAccountName` to an existing Storage Account (since we don't provision one for you).

Assuming you've got Go installed - you can then run this using:

```bash
$ go run main.go
```