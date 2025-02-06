# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "automation_schedule_start_time" {
  value = azurerm_automation_schedule.one-time.start_time
}

output "automation_schedule_week_interval" {
  value = azurerm_automation_schedule.hour.interval
}
