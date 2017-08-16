#Azure loadbalancer Module
resource "azurerm_resource_group" "rg" {
  name     = "${var.prefix}-rg"
  location = "${var.location}"
}

#TODO: Move to provision via module not within LB module itself
resource "azurerm_public_ip" "mypublicIP" {
  name                         = "${var.prefix}-publicIP"
  location                     = "${var.location}"
  resource_group_name          = "${azurerm_resource_group.rg.name}"
  public_ip_address_allocation = "${var.public_ip_address_allocation}"
}

resource "azurerm_lb" "mylb" {
  name                      = "${var.prefix}-lb"  
  resource_group_name       = "${azurerm_resource_group.rg.name}"
  location                  = "${var.location}"
  frontend_ip_configuration {
    name                 = "${var.frontend_name}"
    #TODO: Conditional values for private or public 
    subnet_id            = "${var.frontend_subnet_id}"
    private_ip_address   = "${var.frontend_private_ip}"
    public_ip_address_id = "${azurerm_public_ip.mypublicIP.id}"
  }
  #TODO: tags
}