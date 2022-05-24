## Example: Linux Python Web App deployed from local ZIP

This example provisions a Linux Web App inside an App Service Plan which is configured for Python and deploys a basic Flask App from a local ZIP file.

**Note:** The sample app will deploy allowing access from anywhere.


```bash
$ terraform apply
var.location
  The Azure location where all resources in this example should be created

  Enter a value: westeurope

var.prefix
  The prefix used for all resources in this example

  Enter a value: exampleapp

azurerm_resource_group.example: Creating...
azurerm_resource_group.example: Creation complete after 1s [id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/exampleappRG-zipdeploy]
azurerm_service_plan.example: Creating...
azurerm_service_plan.example: Still creating... [10s elapsed]
azurerm_service_plan.example: Creation complete after 12s [id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/exampleappRG-zipdeploy/providers/Microsoft.Web/serverfarms/exampleapp-sp-zipdeploy]
azurerm_linux_web_app.example: Creating...
azurerm_linux_web_app.example: Still creating... [10s elapsed]
azurerm_linux_web_app.example: Still creating... [20s elapsed]
azurerm_linux_web_app.example: Still creating... [30s elapsed]
azurerm_linux_web_app.example: Still creating... [40s elapsed]
azurerm_linux_web_app.example: Still creating... [50s elapsed]
azurerm_linux_web_app.example: Still creating... [1m0s elapsed]
azurerm_linux_web_app.example: Still creating... [1m10s elapsed]
azurerm_linux_web_app.example: Creation complete after 1m15s [id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/exampleappRG-zipdeploy/providers/Microsoft.Web/sites/exampleapp-zipdeploy]

Apply complete! Resources: 3 added, 0 changed, 0 destroyed.

Outputs:

app_url = "https://exampleapp-zipdeploy.azurewebsites.net"
linux_web_app_name = "exampleapp-zipdeploy"
```

**NOTE:** The source for the example ZIP used here can be found at: [https://github.com/jackofallops/azure-app-service-python-flask-example](https://github.com/jackofallops/azure-app-service-python-flask-example) 