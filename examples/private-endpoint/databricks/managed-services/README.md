## Example: Databricks Workspace with Private Endpoint, Managed-Services and DBFS Customer Managed Keys Enabled

This example provisions a Private Endpoint which connects to a Databricks Workspace within Azure with Managed-Services and DBFS Customer Managed Keys Enabled.

To find the correct Object ID to use for the `azurerm_key_vault_access_policy.managed` `object_id` field in your configuration file you will need to go to [portal](https://portal.azure.com/) -> `Azure Active Directory` and in the `search your tenant` bar enter the value `2ff814a6-3304-4ab8-85cb-cd0e6f879c1d`. You will see under `Enterprise application` results `AzureDatabricks`, click on the `AzureDatabricks` search result. This will open the `Enterprise Application` overview blade where you will see three values, the name of the application, the application ID, and the object ID. The value you want is the object ID, copy this value and paste it into the `object_id` field for your `azurerm_key_vault_access_policy.managed` configuration block.

### Variables

* `prefix` - (Required) The prefix used for all resources in this example.
