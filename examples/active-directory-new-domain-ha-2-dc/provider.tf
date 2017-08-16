# Microsoft Azure Resource Manager Provider

#
# Uncomment this provider block if you have set the following environment variables: 
# ARM_SUBSCRIPTION_ID, ARM_CLIENT_ID, ARM_CLIENT_SECRET and ARM_TENANT_ID
#
provider "azurerm" {}

#
# Uncomment this provider block if you are using variables (NOT environment variables)
# to provide the azurerm provider requirements.
#
# provider "azurerm" {
#   subscription_id = "${var.subscription_id}"
#   client_id       = "${var.client_id}"
#   client_secret   = "${var.client_secret}"
#   tenant_id       = "${var.tenant_id}"
# }

