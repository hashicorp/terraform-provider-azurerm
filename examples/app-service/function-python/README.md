# Example: Azure Function App - Python

This example is used to demonstrate how to provision a python functions app.

This example provisions:

- A storage account
- App service plan
- A Linux Function App
- A Python HTTP Trigger Function for "Hello {name}"

## Variables

- `prefix` - (Required) The prefix used for all resources in this example.
- `location` - (Required) Azure Region in which all resources in this example should be provisioned.

## Outputs

- `app_name` - The name of the app.
- `function_url` - The invocation URL of the function.

## Notes

The `azurerm_function_app_function` resource in `main.tf` is shown with two mechanisms of providing the JSON data to the resource, via the `file()` and `jsonencode()` functions, the latter is commented out.

The Function will be available at:

`https://{prefix}-python-example-app.azurewebsites.net/api/example-python-function?`

And will also be output by this example as `function_url`.

Since this is an anonymous auth function, no API key is required, and the function can be shown to be working by adding the `name` query parameter, e.g.

`https://{prefix}-python-example-app.azurewebsites.net/api/example-python-function?name=world`

*NOTE:* replace `{prefix}` with your `var.prefix` value in the URL's above.
