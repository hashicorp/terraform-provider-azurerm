change "resource-fix" {
  body = "`azurerm_windows_virtual_machine` - correctly populate `bypass_platform_safety_checks_on_user_schedule_enabled` in state when the OS profile is absent"
}
