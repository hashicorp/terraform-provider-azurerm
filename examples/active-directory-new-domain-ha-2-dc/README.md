# Create 2 new Windows VMs, create a new AD Forest, Domain and 2 DCs in an availability set

This template will deploy 2 new VMs (along with a new VNet, Storage Account and Load Balancer) and create a new AD forest and domain. Each VM will be created as a DC for the new domain and will be placed in an availability set. Each VM will also have an RDP endpoint added with a public load balanced IP address.

## Pre-requisites

### Setting up Terraform Access to Azure

To enable Terraform to provision resources into Azure, you need to create two entities in Azure Active Directory (AAD) - AAD Application and AAD Service Principal. [Azure CLI 2.0](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) allows you to quickly provision both by following the instructions below. 

First, login to administer your azure subscription by issuing the following command

```
az login
```

NOTE: If you're using the China, German or Government Azure Clouds, you need to first configure the Azure CLI to work with that Cloud. You can do this by running:

```
az cloud set --name AzureChinaCloud|AzureGermanCloud|AzureUSGovernment
```

If you have multiple Azure Subscriptions, their details are returned by the az login command. 
Set the Subscription that you want to use for this session.

```
az account set --subscription="${SUBSCRIPTION_ID}"
```

Query the account to get the Subscription Id and Tenant Id values.

```
az account show --query "{subscriptionId:id, tenantId:tenantId}"
```

Next, create separate credentials for Terraform.

```
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/${SUBSCRIPTION_ID}"
```

This outputs your client_id (appId), client_secret (password), sp_name, and tenant. Take note of all these variables. Use the returned `appId` value for the `service_principal_client_id` variable in `terraform.tfvars`. Use the password value for the `service_principal_client_secret` variable in `terraform.tfvars`.

NOTE: instead of inserting these values into a `terraform.tfvars` file, you can set corresponding environment variables as described in detail on [docs.microsoft.com](https://docs.microsoft.com/en-us/azure/virtual-machines/terraform-install-configure).

## Running the sample

Once you complete the pre-requisites and fill in all the variables in `terraform.tfvars`, you are ready to provision your infrastructure with Terraform. Start off by running

```
terraform init
```

to initialize Azure provider. 

To see the changes that will be made to your infrastructure (without actually applying them), run the following command

```
terraform plan
```

To apply changes to your infrastructure, run the following command:

```
terraform apply
```

## Optional: configure automated tests/CI environment with Travis CI
In the samples folder, files deploy.ci.sh and deploy.mach.sh are not part of Terraform deployments but rather a part of the continuous integration environment and automated tests that have been setup to validate successful resource deployment. You can setup your own continuous integration environment on [Travis CI](https://travis-ci.org) by modifying travis.yaml to run deploy.ci.sh scripts upon successful code push.


## Further information

For more information on Azure Virtual Machines and Virtual Machine Extensions:

- [Virtual Machine Documentation](https://docs.microsoft.com/en-us/azure/virtual-machines/)
- [Virtual Machine REST API Reference](https://docs.microsoft.com/en-us/rest/api/compute/virtualmachines)
- [Virtual Machine Extension REST API Reference](https://docs.microsoft.com/en-us/rest/api/compute/extensions)

---

This is based on the [active-directory-new-domain-ha-2-dc](https://github.com/Azure/azure-quickstart-templates/tree/master/active-directory-new-domain-ha-2-dc) Azure Quick Start Template.
