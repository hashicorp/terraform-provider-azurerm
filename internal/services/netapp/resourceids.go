package netapp

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Account -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CapacityPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Snapshot -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1/snapshots/snapshot1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Volume -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1
