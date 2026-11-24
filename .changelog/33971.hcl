change "resource-fix" {
  body = "`azurerm_dev_test_schedule` - preserve `notification_settings.0.webhook_url` in state since the API no longer returns it"
}
change "resource-fix" {
  body = "`azurerm_dev_test_global_vm_shutdown_schedule` - preserve `notification_settings.0.webhook_url` in state since the API no longer returns it"
}
