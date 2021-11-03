package mariadb

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Server -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMariaDB/servers/server1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MariaDBFirewallRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMariaDB/servers/server1/firewallRules/firewallRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MariaDBVirtualNetworkRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMariaDB/servers/server1/virtualNetworkRules/vnetrule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MariaDBConfiguration -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMariaDB/servers/server1/configurations/config1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MariaDBDatabase -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforMariaDB/servers/server1/databases/db1
