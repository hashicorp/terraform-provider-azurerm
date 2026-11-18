change "breaking" {
  body = "`azurerm_windows_function_app` - reject `site_config.use_32_bit_worker = true` when `site_config.application_stack.powershell_core_version` is `7.6`; leave `use_32_bit_worker` unset or set it to `false`"
}

change "breaking" {
  body = "`azurerm_windows_function_app_slot` - reject `site_config.use_32_bit_worker = true` when `site_config.application_stack.powershell_core_version` is `7.6`; leave `use_32_bit_worker` unset or set it to `false`"
}

change "resource-fix" {
  body = "`azurerm_windows_function_app` - apply the PowerShell `7.6` stack settings (.NET 10 host and 64-bit worker) when `site_config.application_stack.powershell_core_version` is `7.6`"
}

change "resource-fix" {
  body = "`azurerm_windows_function_app_slot` - apply the PowerShell `7.6` stack settings (.NET 10 host and 64-bit worker) when `site_config.application_stack.powershell_core_version` is `7.6`"
}
