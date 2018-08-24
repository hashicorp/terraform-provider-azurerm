# Configure the Microsoft Azure Provider
provider "azurerm" {
  subscription_id = "${var.subscription_id}"
  client_id       = "${var.client_id}"
  client_secret   = "${var.client_secret}"
  tenant_id       = "${var.tenant_id}"
}

##########################################################
## Create Resource group Network & subnets
##########################################################
module "network" {
  source              = "..\\modules\\network"
  address_space       = "${var.address_space}"
  dns_servers         = ["${var.dns_servers}"]
  environment_name    = "${var.environment_name}"
  resource_group_name = "${var.resource_group_name}"
  location            = "${var.location}"
  dcsubnet_name       = "${var.dcsubnet_name}"
  dcsubnet_prefix     = "${var.dcsubnet_prefix}"
  wafsubnet_name      = "${var.wafsubnet_name}"
  wafsubnet_prefix    = "${var.wafsubnet_prefix}"
  rpsubnet_name       = "${var.rpsubnet_name}"
  rpsubnet_prefix     = "${var.rpsubnet_prefix}"
  issubnet_name       = "${var.issubnet_name}"
  issubnet_prefix     = "${var.issubnet_prefix}"
  dbsubnet_name       = "${var.dbsubnet_name}"
  dbsubnet_prefix     = "${var.dbsubnet_prefix}"
}

##########################################################
## Create DC VM & AD Forest
##########################################################

module "active-directory" {
  source                        = "..\\modules\\active-directory"
  resource_group_name           = "${module.network.out_resource_group_name}"
  location                      = "${var.location}"
  prefix                        = "${var.prefix}"
  subnet_id                     = "${module.network.dc_subnet_subnet_id}"
  active_directory_domain       = "${var.prefix}.local"
  active_directory_netbios_name = "${var.prefix}"
  private_ip_address            = "${var.private_ip_address}"
  admin_username                = "${var.admin_username}"
  admin_password                = "${var.admin_password}"
}

##########################################################
## Create IIS VM's & Join domain
##########################################################

module "iis-vm" {
  source                    = "..\\modules\\iis-vm"
  resource_group_name       = "${module.active-directory.out_resource_group_name}"
  location                  = "${module.active-directory.out_dc_location}"
  prefix                    = "${var.prefix}"
  subnet_id                 = "${module.network.is_subnet_subnet_id}"
  active_directory_domain   = "${var.prefix}.local"
  active_directory_username = "${var.admin_username}"
  active_directory_password = "${var.admin_password}"
  admin_username            = "${var.admin_username}"
  admin_password            = "${var.admin_password}"
  vmcount                   = "${var.vmcount}"
}

##########################################################
## Create Secondary Domain Controller VM & Join domain
##########################################################
module "dc2-vm" {
  source                        =  "..\\modules\\dc2-vm"
  resource_group_name           = "${module.active-directory.out_resource_group_name}"
  location                      = "${module.active-directory.out_dc_location}"
  dcavailability_set_id         = "${module.active-directory.out_dcavailabilityset}"
  prefix                        = "${var.prefix}"
  subnet_id                     = "${module.network.dc_subnet_subnet_id}"
  active_directory_domain       = "${var.prefix}.local"
  active_directory_username     = "${var.admin_username}"
  active_directory_password     = "${var.admin_password}"
  active_directory_netbios_name = "${var.prefix}"
  dc2private_ip_address         = "${var.dc2private_ip_address}"
  admin_username                = "${var.admin_username}"
  admin_password                = "${var.admin_password}"
  domainadmin_username          = "${var.domainadmin_username}"
}

##########################################################
## Create SQL Server VM Join domain
##########################################################
module "sql-vm" {
  source                    = "..\\modules\\sql-vm"
  resource_group_name       = "${module.active-directory.out_resource_group_name}"
  location                  = "${module.active-directory.out_dc_location}"
  prefix                    = "${var.prefix}"
  subnet_id                 = "${module.network.db_subnet_subnet_id}"
  active_directory_domain   = "${var.prefix}.local"
  active_directory_username = "${var.admin_username}"
  active_directory_password = "${var.admin_password}"
  admin_username            = "${var.admin_username}"
  admin_password            = "${var.admin_password}"
  sqlvmcount                = "${var.sqlvmcount}"
  lbprivate_ip_address      = "${var.lbprivate_ip_address}"
}