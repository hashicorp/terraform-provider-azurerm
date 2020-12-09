package loganalytics

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogAnalyticsWorkspace -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogAnalyticsStorageInsights -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/storageInsightConfigs/storageInsight1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogAnalyticsSavedSearch -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/savedSearches/search1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogAnalyticsCluster -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/clusters/cluster1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogAnalyticsSolution -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationsManagement/solutions/solution1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogAnalyticsClusterCustomerManagedKey -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.OperationalInsights/clusters/cluster1/key1
