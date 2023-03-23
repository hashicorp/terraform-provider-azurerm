# Example: Linux Python Web App deployed from local ZIP

This example provisions a Linux Web App inside an App Service Plan which is configured for Python and deploys a basic Flask App from a local ZIP file.

## Variables

- `prefix` - (Required) The prefix used for all resources in this example.
- `location` - (Required) Azure Region in which all resources in this example should be provisioned.

## Outputs

- `app_name` - The name of the app.
- `app_url` - The default URL to access the app.
- `app_uptime` - The "uptime" endpoint of the example app. Returns the number of seconds since the process started for the example app.
- `app_healthcheck_endpoint` - the "health" endpoint of the example app. Returns `200 OK` with a body of `OK` when the app is started.

**Note:** The sample app will deploy allowing access from anywhere.

**NOTE:** The source for the example ZIP used here can be found at: [https://github.com/jackofallops/azure-app-service-python-flask-example](https://github.com/jackofallops/azure-app-service-python-flask-example) 