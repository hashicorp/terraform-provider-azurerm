package datadog

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatadogMonitor -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatadogTagRules -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1/tagRules/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatadogSingleSignOnConfigurations -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Datadog/monitors/monitor1/singleSignOnConfigurations/default
