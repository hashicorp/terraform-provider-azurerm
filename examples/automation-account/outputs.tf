
output "automation_schedule-start_time" {
  value = "${azurerm_automation_schedule.one-time.start_time}"
}

output "automation_schedule-week-interval" {
  value = "${azurerm_automation_schedule.hour.interval}"
}