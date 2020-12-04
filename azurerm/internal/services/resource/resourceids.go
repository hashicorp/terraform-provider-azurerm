package resource

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ResourceGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ResourceGroupTemplateDeployment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Resources/deployments/deploy1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SubscriptionTemplateDeployment -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Resources/deployments/deploy1
