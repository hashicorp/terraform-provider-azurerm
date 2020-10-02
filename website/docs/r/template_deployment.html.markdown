---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_template_deployment"
description: |-
  Manages a template deployment of resources.
---

# azurerm_template_deployment

Manages a template deployment of resources

~> **Note:** The `azurerm_template_deployment` resource has been superseded by the [`azurerm_resource_group_template_deployment`](resource_group_template_deployment.html) resource. The existing `azurerm_template_deployment` resource will be deprecated (but still available) in version 3.0 of the AzureRM Terraform Provider - we recommend using the `azurerm_resource_group_template_deployment` resource for new deployments.

-> **Note:** This resource will not clean up nested resources deployed by the ARM Template Deployment. We recommend using the [`azurerm_resource_group_template_deployment`](resource_group_template_deployment.html) resource for new deployments, which can do this.

## Example Usage

~> **Note:** This example uses [Storage Accounts](storage_account.html) and [Public IP's](public_ip.html) which are natively supported by Terraform - we'd highly recommend using the Native Resources where possible instead rather than an ARM Template, for the reasons outlined above.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_template_deployment" "example" {
  name                = "acctesttemplate-01"
  resource_group_name = azurerm_resource_group.example.name

  template_body = <<DEPLOY
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "storageAccountType": {
      "type": "string",
      "defaultValue": "Standard_LRS",
      "allowedValues": [
        "Standard_LRS",
        "Standard_GRS",
        "Standard_ZRS"
      ],
      "metadata": {
        "description": "Storage Account type"
      }
    }
  },
  "variables": {
    "location": "[resourceGroup().location]",
    "storageAccountName": "[concat(uniquestring(resourceGroup().id), 'storage')]",
    "publicIPAddressName": "[concat('myPublicIp', uniquestring(resourceGroup().id))]",
    "publicIPAddressType": "Dynamic",
    "apiVersion": "2015-06-15",
    "dnsLabelPrefix": "terraform-acctest"
  },
  "resources": [
    {
      "type": "Microsoft.Storage/storageAccounts",
      "name": "[variables('storageAccountName')]",
      "apiVersion": "[variables('apiVersion')]",
      "location": "[variables('location')]",
      "properties": {
        "accountType": "[parameters('storageAccountType')]"
      }
    },
    {
      "type": "Microsoft.Network/publicIPAddresses",
      "apiVersion": "[variables('apiVersion')]",
      "name": "[variables('publicIPAddressName')]",
      "location": "[variables('location')]",
      "properties": {
        "publicIPAllocationMethod": "[variables('publicIPAddressType')]",
        "dnsSettings": {
          "domainNameLabel": "[variables('dnsLabelPrefix')]"
        }
      }
    }
  ],
  "outputs": {
    "storageAccountName": {
      "type": "string",
      "value": "[variables('storageAccountName')]"
    }
  }
}
DEPLOY


  # these key-value pairs are passed into the ARM Template's `parameters` block
  parameters = {
    "storageAccountType" = "Standard_GRS"
  }

  deployment_mode = "Incremental"
}

output "storageAccountName" {
  value = azurerm_template_deployment.example.outputs["storageAccountName"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the template deployment. Changing this forces a
    new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to
    create the template deployment.
* `deployment_mode` - (Required) Specifies the mode that is used to deploy resources. This value could be either `Incremental` or `Complete`.
    Note that you will almost *always* want this to be set to `Incremental` otherwise the deployment will destroy all infrastructure not
    specified within the template, and Terraform will not be aware of this.
* `template_body` - (Optional) Specifies the JSON definition for the template.

~> **Note:** There's a [`file` function available](https://www.terraform.io/docs/configuration/functions/file.html) which allows you to read this from an external file, which helps makes this more resource more readable.

* `parameters` - (Optional) Specifies the name and value pairs that define the deployment parameters for the template.

* `parameters_body` - (Optional) Specifies a valid Azure JSON parameters file that define the deployment parameters. It can contain KeyVault references

~> **Note:** There's a [`file` function available](https://www.terraform.io/docs/configuration/functions/file.html) which allows you to read this from an external file, which helps makes this more resource more readable.

~> **Also Note:** This is NOT an Azure deployment parameters file, as defined in the [`Microsoft schema's`](https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#). It is effectively the object supplied to the "parameters" attribute in that schema. If you are providing, or generating via [`template_file`](https://www.terraform.io/docs/providers/template/d/file.html), this argument, do not provide a full deployment parameters JSON file with "$schema" and "contentVersion" attributes, just provide the object for the "parameters" attribute of that schema.

## Attributes Reference

The following attributes are exported:

* `id` - The Template Deployment ID.

* `outputs` - A map of supported scalar output types returned from the deployment (currently, Azure Template Deployment outputs of type String, Int and Bool are supported, and are converted to strings - others will be ignored) and can be accessed using `.outputs["name"]`.

## Note

Terraform does not know about the individual resources created by Azure using a deployment template and therefore cannot delete these resources during a destroy. Destroying a template deployment removes the associated deployment operations, but will not delete the Azure resources created by the deployment. In order to delete these resources, the containing resource group must also be destroyed. [More information](https://docs.microsoft.com/en-us/rest/api/resources/deployments#Deployments_Delete).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Template Deployment.
* `update` - (Defaults to 3 hours) Used when updating the Template Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Template Deployment.
* `delete` - (Defaults to 3 hours) Used when deleting the Template Deployment.
