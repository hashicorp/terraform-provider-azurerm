# Sample Serverless Application on Azure: Integration Between Function App and API Management with Terraform

This example demonstrates how to integrate Azure API Management with an Azure Linux Function App.

# Requirements

The following versions were used and confirmed to work. While other versions may also work, they have not been tested with this setup:

- **Azure CLI**: Version 2.65.0  
- **Terraform CLI**: Version 1.9.7  
- **Terraform AzureRM Provider**: Version 4.0 or later  
- **Azure Functions Core Tools**: Version 4.0.6543  
- **Node.js**: Version 22.9.0  
- **Node Package Manager (NPM)**: Version 10.8.3  

# How to Run

## Deploy Infrastructure

1. Log in to your Azure account:
   ```
   az login
   ```

2. Retrieve the subscription ID and resource group name from your Azure environment. Update the `subscription_id` and `resource_group_name` values in the `terraform_configs/variables.tf` file accordingly. To obtain the subscription ID and resource group name, you can use:
   ```
   az group list
   ```

3. Navigate to the `terraform_configs` folder:
   ```
   cd terraform_configs
   ```

4. Initialize and deploy the Terraform configuration:
   ```
   terraform init
   terraform plan -out main.tfplan
   terraform apply main.tfplan
   ```

   **Note:** Deploying the Azure API Management instance can take up to 90 minutes. Please be patient.

5. Take note of the following outputs:
   - **`function_app_name`**: The name of your Function App, such as `myfuncappsbigwbgyzdync`, which you will need to deploy your function code.
   - **`frontend_url`**: The frontend URL of your API Management instance.

## Deploy the Code to the Function App

Azure Functions Core Tools has been used to package and deploy the Function App code. While other tools can be used, Terraform currently does not fully integrate with such deployment tools, requiring the functions to be deployed separately after each Terraform deployment. Automation can simplify this process using scripts.

1. Navigate to the `function_code` folder:
   ```
   cd ../function_code
   ```

2. Install the code dependencies:
   ```
   npm install
   ```

3. Deploy the code:
   ```
   func azure functionapp publish your_function_app_name
   ```
   Replace `your_function_app_name` with the `function_app_name` from the Terraform outputs.

   **Note:** The `func` CLI may occasionally encounter errors, even if it reports `Deployment completed successfully`. If the `func azure functionapp publish` command does not return the `Invoke URL`, rerun the command.

4. Verify the deployment:
   - **Backend URL (`Invoke URL`)**: Send a GET request to this URL. The response should return a 200 status code with the following message:
     ```
     "Hello world, this is coming from Function App!"
     ```

   - **Frontend URL (`frontend_url`)**: Send a GET request to this URL. If everything is correctly configured, the frontend should return the same 200 status code and message.

5. Remember: After each Terraform deployment, you may need to redeploy the Function App code.

> [!NOTE]  
> For easier development and debugging, CORS restrictions have been disabled by setting `Access-Control-Allow-Origin: *` within the function code. Once the application is running successfully, ensure CORS is re-enabled and properly configured to secure the application.

# Additional Notes

Integration between Azure Function App and Azure API Management can be challenging and prone to various issues. Due to the limited documentation, I have automated variables to reduce potential errors. If you plan to modify the code, keep the following considerations in mind:

- Ensure the `url_template` property in the `azurerm_api_management_api_operation` resource aligns with the function name specified in the `Invoke URL`. Avoid appending slashes to the `url_template` value.
- Ensure that `https://your-function-app-name.azurewebsites.net/api` is consistently used in both the policy and the `Service URL` in the backend configuration, and ensure it overrides any existing value in the backend.

---

This sample aims to bridge gaps in existing documentation and provide a working, reproducible example of Azure Function App and API Management integration. Contributions and feedback are welcome!
