variable "resource_group" {
  description = "The name of the resource group in which to create the database."
}

variable "location" {
  description = "The location/region where the database and server are created. Changing this forces a new resource to be created."
  default     = "southcentralus"
}

variable "db_name" {
  description = "The name of the database to be created."
}

variable "sql_admin" {
  description = "The administrator username of the SQL Server."
}

variable "sql_password" {
  description = "The administrator password of the SQL Server."
}
