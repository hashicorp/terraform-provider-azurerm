resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_template_deployment" "example" {
  name                = "${var.prefix}-deployment"
  resource_group_name = "${azurerm_resource_group.example.name}"

  template_body = <<DEPLOY
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
  },
  "variables": {
  },
  "resources": [
  ],
  "outputs": {
  }
}
DEPLOY

  deployment_mode = "Incremental"
}
