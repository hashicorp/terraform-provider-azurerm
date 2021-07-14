package appservice

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=WebApp -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServicePlan -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/serverfarms/farm1
