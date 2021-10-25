package applicationinsights

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Component -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SmartDetectionRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1/smartDetectionRule/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=WebTest -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/webTests/test1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApiKey -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1/apiKeys/apikey1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AnalyticsUserItem -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1/myAnalyticsItems/item1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AnalyticsSharedItem -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1/analyticsItems/item1
