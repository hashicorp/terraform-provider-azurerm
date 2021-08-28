package databricks

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Workspace -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/workspaces/workspace1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CustomerManagedKey -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Databricks/customerMangagedKey/workspace1
