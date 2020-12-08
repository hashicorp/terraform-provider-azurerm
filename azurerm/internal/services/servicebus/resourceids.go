package servicebus

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Namespace -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NamespaceNetworkRuleSet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/networkrulesets/networkRuleSet1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Topic -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceBus/namespaces/namespace1/topics/topic1
