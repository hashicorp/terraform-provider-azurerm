package monitor

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ActionGroup -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/actionGroups/actionGroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ActionRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.AlertsManagement/actionRules/actionRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SmartDetectorAlertRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.AlertsManagement/smartdetectoralertrules/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ActivityLogAlert -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/activityLogAlerts/alert1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AutoscaleSetting -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/autoscaleSettings/setting1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogProfile -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Insights/logProfiles/profile1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MetricAlert -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/metricAlerts/alert1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateLinkScope -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/privateLinkScopes/pls1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateLinkScopedService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/privateLinkScopes/pls1/scopedResources/sr1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ScheduledQueryRules -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/scheduledQueryRules/rule1
