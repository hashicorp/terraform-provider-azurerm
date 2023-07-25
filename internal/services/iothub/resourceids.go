// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Enrichment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/enrichments/enrichment1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IotHub -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ConsumerGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/eventHubEndpoints/events/consumerGroups/group1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FallbackRoute -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/fallbackRoute/default -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Route -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/routes/route1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SharedAccessPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/iotHubKeys/sharedAccessPolicy1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EndpointStorageContainer -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/endpoints/storageContainerEndpoint1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EndpointServiceBusTopic -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/endpoints/serviceBusTopicEndpoint1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EndpointServiceBusQueue -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/endpoints/serviceBusQueueEndpoint1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EndpointEventhub -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/endpoints/eventHubEndpoint1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IotHubCertificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/certificates/cert1 -rewrite=true
