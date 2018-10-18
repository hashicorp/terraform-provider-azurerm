variable "resource_group" {
  description = "The name of the resource group in which to create the SQL Server."
}

variable "location" {
  description = "The location/region where the SQL Server is created. Changing this forces a new resource to be created."
  default     = "southcentralus"
}

variable "sql_admin" {
  description = "The administrator username of the SQL Server."
}

variable "sql_password" {
  description = "The administrator password of the SQL Server."
}
