package cdn

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Endpoint -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Profile -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CustomDomain -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1/customDomains/domain1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AfdOriginGroups -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/origingroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AfdOrigins -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/origingroup1/origins/origin1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AfdCustomDomain -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/custom-domain-com
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AfdEndpoints -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/afdEndpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AfdRuleSets -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AfdRuleRules -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1/rules/rule1
