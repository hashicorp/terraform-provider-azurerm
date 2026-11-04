change "resource-enhancement" {
  body = "`azurerm_container_app` - add `azure-durabletask-scheduler` as a valid value for the `template.custom_scale_rule.custom_rule_type` property"
}
change "resource-enhancement" {
  body = "`azurerm_container_app_job` - add `azure-durabletask-scheduler` as a valid value for the `event_trigger_config.scale.rules.custom_rule_type` property"
}
