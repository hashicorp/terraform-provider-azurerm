package trafficmanager

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AzureEndpoint -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/trafficManagerProfiles/trafficManagerProfile1/azureEndpoints/azureEndpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ExternalEndpoint -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/trafficManagerProfiles/trafficManagerProfile1/externalEndpoints/externalEndpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NestedEndpoint -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/trafficManagerProfiles/trafficManagerProfile1/nestedEndpoints/nestedEndpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=TrafficManagerProfile -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/trafficManagerProfiles/trafficManagerProfile1
