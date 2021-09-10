package iothub

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Enrichment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IotHubs/hub1/Enrichments/enrichment1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IotHub -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IotHubs/hub1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EndpointStorageContainer -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IotHubs/hub1/Endpoints/storageContainerEndpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EndpointServicebusTopic -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IotHubs/hub1/Endpoints/serviceBusTopicEndpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EndpointServicebusQueue -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IotHubs/hub1/Endpoints/serviceBusQueueEndpoint1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EndpointEventhub -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IotHubs/hub1/Endpoints/eventHubEndpoint1
