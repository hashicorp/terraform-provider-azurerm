// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Component -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SmartDetectionRule -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1/smartDetectionRule/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=WebTest -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/webTests/test1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApiKey -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1/apiKeys/apikey1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AnalyticsUserItem -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1/myAnalyticsItems/item1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AnalyticsSharedItem -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/components/component1/analyticsItems/item1
