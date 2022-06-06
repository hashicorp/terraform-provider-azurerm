# Example: App Service configured for Windows Container

This example provisions an App Service inside an App Service Plan which is configured to run a Windows Container.

## Variables

- `prefix` - (Required) The prefix used for all resources in this example.
- `location` - (Required) Azure Region in which all resources in this example should be provisioned.

## Outputs

- `app_name` - The name of the app.
- `app_url` - The default URL to access the app.
