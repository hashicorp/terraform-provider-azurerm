// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

// CDN
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Endpoint -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Profile -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CustomDomain -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/endpoints/endpoint1/customDomains/domain1

// CDN FrontDoor
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorCustomDomain -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/customDomains/customDomain1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorEndpoint -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorFirewallPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/policy1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorOrigin -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1/origins/origin1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorOriginGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/originGroups/originGroup1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorProfile -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorRoute -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/routes/route1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorRuleSet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1/rules/rule1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorSecret -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/secrets/secret1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorSecurityPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/securityPolicies/securityPolicy1 -rewrite=true

// CDN FrontDoor "Associations"
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorRouteDisableLinkToDefaultDomain -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/routes/route1/disableLinkToDefaultDomain/disableLinkToDefaultDomain1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoorCustomDomainAssociation -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/associations/assoc1
