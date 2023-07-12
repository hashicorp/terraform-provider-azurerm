// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoor

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=BackendPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/backendPools/pool1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontDoor -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FrontendEndpoint -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/frontendEndpoints/endpoint1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HealthProbe -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/healthProbeSettings/probe1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancing -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/loadBalancingSettings/setting1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RoutingRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/routingRules/rule1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=WebApplicationFirewallPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/policy1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CustomHttpsConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontDoors/frontdoor1/customHttpsConfiguration/endpoint1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LoadBalancingRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/loadBalancers/loadBalancer1/loadBalancingRules/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RulesEngine -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/frontdoors/frontdoor1/rulesEngines/rule1 -rewrite=true
