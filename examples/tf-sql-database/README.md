# Create a SQL Server Database in Azure

You can use this module to create a basic SQL Server Database in Azure. Below is an example of how you would use this module in your code:

```
module "sql-database" {
  source             = "./path/to/module"
  resource_group     = "my-resource-group"
  location           = "west us"
  db_name            = "mydatabase"
  sql_admin_username = "adminaccount"
  sql_password       = "adminpassword"
}
```
