package datafactory

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IntegrationRuntime -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DataFactory/factories/factory1/integrationruntimes/runtime1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LinkedService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DataFactory/factories/factory1/linkedservices/linkedService1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DataSet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DataFactory/factories/facName1/datasets/dataSet1
