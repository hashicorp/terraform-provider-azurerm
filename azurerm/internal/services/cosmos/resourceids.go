package cosmos

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CassandraKeyspace -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/DatabaseAccounts/acc1/cassandraKeyspaces/keyspace1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabaseAccount -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/DatabaseAccounts/acc1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=GremlinDatabase -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DocumentDB/DatabaseAccounts/acc1/gremlinDatabases/database1
