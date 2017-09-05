# Create a Virtual Machine from a User Image

> Prerequisite - The Storage Account with the User Image VHD should already exist

This template allows you to create a Virtual Machines from a User image using managed disk. This template also deploys a Virtual Network, Public IP addresses and a Network Interface.

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

### Generate an ssh key

Generate an ssh key as follows:

```
ssh-keygen -t rsa -b 2048 
```

Copy the contents of the following and place into the `admin_ssh_publickey` variable in `terraform.tfvars`:

```
cat ~/.ssh/id_rsa.pub
```

There are instructions for using PuTTY on Windows to generate your ssh keys [here](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ssh-from-windows).

More information on using ssh with VMs in Azure:

- [How to create and use an SSH public and private key pair for Linux VMs in Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/mac-create-ssh-keys)
- [How to Use SSH keys with Windows on Azure](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ssh-from-windows)

## Running the sample

Once you complete the pre-requisites and fill in all the variables in `terraform.tfvars`, you are ready to provision your infrastructure with Terraform. If this is your first time running Terraform for Azure, run the command

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

This is based on the [101-vm-from-user-image](https://github.com/Azure/azure-quickstart-templates/tree/master/101-vm-from-user-image) Azure Quick Start Template.