package datalake

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Account -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DataLakeStore/accounts/account1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DataLakeStore/accounts/account1/virtualNetworkRules/virtualNetworkRule1
