package consumption

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ConsumptionBudgetResourceGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Consumption/budgets/budget1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ConsumptionBudgetSubscription -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Consumption/budgets/budget1
