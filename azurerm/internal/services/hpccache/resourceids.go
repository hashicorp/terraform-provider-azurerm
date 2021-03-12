package hpccache

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Cache -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.StorageCache/caches/cache1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=StorageTarget -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.StorageCache/caches/cache1/storageTargets/target1
