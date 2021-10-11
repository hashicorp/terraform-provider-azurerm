package monitor

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ActionGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/actionGroups/actionGroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ActionRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.AlertsManagement/actionRules/actionRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SmartDetectorAlertRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.AlertsManagement/smartdetectoralertrules/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ActivityLogAlert -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/activityLogAlerts/alert1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AutoscaleSetting -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/autoscaleSettings/setting1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogProfile -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Insights/logProfiles/profile1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MetricAlert -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/metricAlerts/alert1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ScheduledQueryRules -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/scheduledQueryRules/rule1
