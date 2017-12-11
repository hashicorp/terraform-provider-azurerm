# Azure App Service Sample

Sample to deploy an App Service within an App Service Plan.

## Creates

1. A Resource Group
2. An [App Service Plan](https://docs.microsoft.com/en-us/azure/app-service/azure-web-sites-web-hosting-plans-in-depth-overview)
3. An [App Service](https://azure.microsoft.com/en-gb/services/app-service/) configured for usage with .NET 4.x Application

## Usage

- Provide values to all variables (credentials and names).
- Create with `terraform apply`
- Destroy all with `terraform destroy --force`
