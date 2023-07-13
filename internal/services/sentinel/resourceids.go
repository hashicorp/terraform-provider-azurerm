// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SentinelAlertRuleTemplate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRuleTemplates/template1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DataConnector -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AutomationRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/automationRules/rule1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MLAnalyticsSettings -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/securityMLAnalyticsSettings/setting1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ThreatIntelligenceIndicator -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/threatIntelligence/main/indicators/indicator1
