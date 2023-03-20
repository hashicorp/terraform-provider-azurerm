package databricks

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabricksVirtualNetworkPeering -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1/virtualNetworkPeerings/peer1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabricksWorkspaces -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1 -rewrite=true
