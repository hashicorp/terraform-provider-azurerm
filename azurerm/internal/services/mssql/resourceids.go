package mssql

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Database -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/databases/database1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabaseExtendedAuditingPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/databases/database1/extendedAuditingSettings/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ElasticPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/elasticPools/pool1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RecoverableDatabase -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/recoverabledatabases/database1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Server -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerExtendedAuditingPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/extendedAuditingSettings/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Server -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/virtualMachine1
