package sentinel

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AlertRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRules/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SentinelAlertRuleTemplate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/AlertRuleTemplates/template1
