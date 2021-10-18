package firewall

// Firewall Policy
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Firewall -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Network/azureFirewalls/firewall1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallApplicationRuleCollection -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/applicationRuleCollections/applicationRuleCollection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallNatRuleCollection -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/natRuleCollections/natRuleCollection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallNetworkRuleCollection -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/azureFirewalls/myfirewall/networkRuleCollections/networkRuleCollection1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallPolicyRuleCollectionGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/firewallPolicies/policy1/ruleCollectionGroups/ruleCollectionGroup1
