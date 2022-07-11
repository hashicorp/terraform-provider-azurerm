package logz

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogzMonitor -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Logz/monitors/monitor1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogzTagRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Logz/monitors/monitor1/tagRules/ruleSet1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogzSubAccount -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Logz/monitors/monitor1/accounts/subAccount1
