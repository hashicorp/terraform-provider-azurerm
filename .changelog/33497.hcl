change "new-property" {
  body = "`azurerm_linux_function_app` - add support for the `container_app_environment_id` property"
}
change "new-property" {
  body = "Data Source: `azurerm_linux_function_app` - export the `container_app_environment_id` property"
}
change "resource-fix" {
  body = "`azurerm_linux_function_app_slot` - fix a panic when the parent `azurerm_linux_function_app` is hosted on a Container Apps Environment"
}
change "resource-fix" {
  body = "`azurerm_container_app_environment` - fix an issue where `Delete` could fail to correctly wait for the resource to be deleted"
}
