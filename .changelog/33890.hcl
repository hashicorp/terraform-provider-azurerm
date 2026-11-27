change "resource-fix" {
  body = "`azurerm_linux_web_app` - correctly decode `site_config.application_stack.docker_image_name` when the fx version has a lowercase `docker|` prefix"
}
change "resource-fix" {
  body = "`azurerm_linux_web_app_slot` - correctly decode `site_config.application_stack.docker_image_name` when the fx version has a lowercase `docker|` prefix"
}
change "resource-fix" {
  body = "`azurerm_windows_web_app` - correctly decode `site_config.application_stack.docker_image_name` when the fx version has a lowercase `docker|` prefix"
}
change "resource-fix" {
  body = "`azurerm_windows_web_app_slot` - correctly decode `site_config.application_stack.docker_image_name` when the fx version has a lowercase `docker|` prefix"
}
