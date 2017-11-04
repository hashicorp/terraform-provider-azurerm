variable "resource_group_name" {
  type    = "string"
  default = "yourGroupName"
}

variable "location" {
  type    = "string"
  default = "westeurope"
}

variable "app_service_name" {
  type = "string"

  # will create -> mySite123456.azurewebsites.net  #should be unique
  # default = "mySite123456"  
}

variable "deploy_user" {
  type = "string"
}

variable "deploy_pass" {
  type = "string"
}

variable "deployZipFile" {
  type    = "string"
  default = "web-package.zip"
}
