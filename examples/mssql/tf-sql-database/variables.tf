variable "resource_group" {
  description = "The name of the resource group in which to create the database."
}

variable "location" {
  description = "The location/region where the database and server are created. Changing this forces a new resource to be created."
}

variable "db_name" {
  description = "The name of the database to be created."
}

variable "db_edition" {
  description = "The edition of the database to be created."
  default     = "Basic"
}

variable "sql_admin_username" {
  description = "The administrator username of the SQL Server."
}

variable "sql_password" {
  description = "The administrator password of the SQL Server."
}

variable start_ip_address {
  description = "Defines the start IP address used in your database firewall rule."
  default     = "0.0.0.0"
}

variable end_ip_address {
  description = "Defines the end IP address used in your database firewall rule."
  default     = "0.0.0.0"
}
