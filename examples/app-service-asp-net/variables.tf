variable "groupName" {
  type    = "string"
  default = "yourGroupName"
}

variable "location" {
  type    = "string"
  default = "westeurope"
}

variable "webName" {
  type = "string"

  # will create -> mySite123456.azurewebsites.net  #should be unique
  # default = "mySite123456"  
}

variable "subscription_id" {
  type = "string"
}

variable "client_id" {
  type = "string"
}

variable "client_secret" {
  type = "string"
}

variable "tenant_id" {
  type = "string"
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
