description = "Scheduled query rule example resource with log query and schedule"
enabled = true
query = "requests | summarize AggregatedValue = count() by bin(TimeGenerated, 5m)"
frequency_in_minutes = 5
time_window_in_minutes = 30