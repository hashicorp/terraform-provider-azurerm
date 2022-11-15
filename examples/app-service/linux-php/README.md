# Example: App Service configured for PHP

This example provisions an App Service inside an App Service Plan which is configured for PHP.

When provisioning this example - the `repository_url` is output - however since the username includes `$` sign (as shown below) this will need to be escaped:

## Variables

- `prefix` - (Required) The prefix used for all resources in this example.
- `location` - (Required) Azure Region in which all resources in this example should be provisioned.

## Outputs

- `app_name` - The name of the app.
- `repository_url` - The URL of the SCM site for the app.

## Example use

```bash
$ terraform apply
var.location
  The Azure location where all resources in this example should be created

  Enter a value: westeurope

var.prefix
  The prefix used for all resources in this example

  Enter a value: tomdev099

azurerm_resource_group.example: Refreshing state... (ID: /subscriptions/00000000-0000-0000-0000-...000/resourceGroups/tomdev099-resources)
azurerm_app_service_plan.example: Refreshing state... (ID: /subscriptions/00000000-0000-0000-0000-...icrosoft.Web/serverfarms/tomdev099-asp)
azurerm_app_service.example: Refreshing state... (ID: /subscriptions/00000000-0000-0000-0000-...crosoft.Web/sites/tomdev099-appservice)

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

repository_url = https://$tomdev099-appservice:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX@tomdev099-appservice.scm.azurewebsites.net/tomdev099-appservice.git
```

You can escape the `$` character in Bash by using a Backslash, for example:

```
$ git remote add origin	https://\$tomdev099-appservice:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX@tomdev099-appservice.scm.azurewebsites.net/tomdev099-appservice.git
```
