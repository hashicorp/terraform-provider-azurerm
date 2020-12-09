package redis

// TODO: fix the generator bug making this `Redi` rather than `Redis`

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Cache -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cache/Redis/redis1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cache/Redis/redis1/firewallRules/firewallRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LinkedServer -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cache/Redis/redis1/linkedServers/linkedServer1
