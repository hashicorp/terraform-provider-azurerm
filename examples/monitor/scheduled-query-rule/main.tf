resource "azurerm_resource_group" "example" {
  name     = format("%s-resources", var.prefix)
  location = var.location
}

resource "azurerm_application_insights" "example" {
  name                = format("%s-insights", var.prefix)
  location            = var.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_scheduled_query_rule" "example" {
  name                   = format("%s-queryrule", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name

  enabled                = var.enabled
  description            = var.description
  frequency_in_minutes   = var.schedule.frequency_in_minutes
  time_window_in_minutes = var.schedule.time_window_in_minutes
  query                  = var.query
  data_source_id         = azurerm_application_insights.example.id
  authorized_resources   = [azurerm_application_insights.example.id]
  query_type             = "ResultCount"
  action                 = ""
}

# Maybe add result count trigger to example
# Creates a Log Alert Rule (Scheduled Query Rule type)
# $triggerCondition = New-AzScheduledQueryRuleTriggerCondition -ThresholdOperator "GreaterThan" -Threshold 3 -MetricTrigger $metricTrigger
# $alertingAction = New-AzScheduledQueryRuleAlertingAction -AznsAction $aznsActionGroup -Severity "1" -Trigger $triggerCondition
