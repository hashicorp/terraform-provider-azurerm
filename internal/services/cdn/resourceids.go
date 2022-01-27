package cdn

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontdoorProfileOriginGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontdoorProfileCustomDomain -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontdoorProfileOrigin -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1/origins/origin1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontdoorProfileEndpoint -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Endpoint -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Profile -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CustomDomain -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1/customDomains/domain1
