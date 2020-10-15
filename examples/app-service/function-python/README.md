## Example: Azure Function App - Python

This example is used to demonstrate how to provision an python functions app. It also contains sample python app that can be deployed.

This example provisions:
- A storage account
- App service plan
- A Functions App

To deploy Azure Function app, refer to these references:
- [Tutorial: Create and deploy serverless Azure Functions in Python with Visual Studio Code](https://docs.microsoft.com/en-us/azure/developer/python/tutorial-vs-code-serverless-python-01)
- [Quickstart: Create a function in Azure that responds to HTTP requests](https://docs.microsoft.com/en-us/azure/azure-functions/functions-create-first-azure-function-azure-cli?pivots=programming-language-python&tabs=bash%2Cbrowser)
- [Continuous deployment for Azure Functions](https://docs.microsoft.com/en-us/azure/azure-functions/functions-continuous-deployment)
- [Continuous delivery by using Azure DevOps](https://docs.microsoft.com/en-us/azure/azure-functions/functions-how-to-azure-devops?tabs=python)
- [Continuous delivery by using GitHub Action](https://docs.microsoft.com/en-us/azure/azure-functions/functions-how-to-github-actions?tabs=javascript)

### Variables

- `storageaccount` - (Required) Name of the storage account
- `location` - (Required) Azure Region in which all resources in this example should be provisioned