description = "Scheduled query rule example resource with log query and schedule"
enabled = true
query = "exceptions | summarize AggregatedValue = count() by bin(TimeGenerated, 5m)"
schedule = {
  frequency_in_minutes = 5
  time_window_in_minutes = 30
}
data_source_id = "/subscriptions/b67f7fec-69fc-4974-9099-a26bd6ffeda3/resourceGroups/MyResourceGroup/providers/Microsoft.OperationalInsights/workspaces/MyWorkspace"
