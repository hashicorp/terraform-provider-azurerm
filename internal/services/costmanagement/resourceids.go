// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package costmanagement

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SubscriptionCostManagementExport -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/exports/export1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ResourceGroupCostManagementExport -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.CostManagement/exports/export1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AnomalyAlertView -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/views/ms:DailyAnomalyByResourceGroup
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SubscriptionCostManagementView -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.CostManagement/views/view1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ResourceGroupCostManagementView -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.CostManagement/views/view1
