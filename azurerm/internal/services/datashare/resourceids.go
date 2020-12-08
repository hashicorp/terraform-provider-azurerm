package datashare

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Account -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DataSet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/shares/share1/dataSets/dataSet1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Share -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DataShare/accounts/account1/shares/share1
