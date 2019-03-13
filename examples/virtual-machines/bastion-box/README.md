# A FreeBSD Jumpbox 

This template allows you to deploy a simple FreeBSD jumpbox VM using the latest patched version. This will deploy an A0 size VM in the resource group location and return the FQDN of the VM.

This template takes a minimum amount of parameters and deploys FreeBSD as a jumpbox VM on an isolated subnet (management subnet). A second subnet named Web is also created as a placeholder.

## main.tf
The `main.tf` file contains the actual resources that will be deployed. It also contains the Azure Resource Group definition and any defined variables.

Azure requires that an application is added to Azure Active Directory to generate the `client_id`, `client_secret`, and `tenant_id` needed by Terraform (`subscription_id` can be recovered from your Azure account details). Please go [here](https://www.terraform.io/docs/providers/azurerm/) for full instructions on how to create this to populate the azurerm provider block or environment variables.

## outputs.tf
This data is outputted when `terraform apply` is called, and can be queried using the `terraform output` command.

## terraform.tfvars
If a `terraform.tfvars` or any `.auto.tfvars` files are present in the current directory, Terraform automatically loads them to populate variables. We don't recommend saving usernames and password to version control, but you can create a local secret variables file and use the `-var-file` flag or the `.auto.tfvars` extension to load it.

## variables.tf
The `variables.tf` file contains all of the input parameters that the user can specify when deploying this Terraform template.

## Post-Deployment

1. The FreeBSD jumpbox will have a public IP and can be accessed through SSH using public keys only. The default username is `vmadmin` but that can be changed on the variables.tf files.

![graph](graph.png)
