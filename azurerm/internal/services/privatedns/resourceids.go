package privatedns

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PrivateDnsZone -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkLink -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/virtualNetworkLinks/virtualNetworkLink1

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ARecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/A/eh1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AaaaRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/AAAA/eheh1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CnameRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/CNAME/name1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MxRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/MX/mx1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PtrRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/PTR/ptr1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SrvRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/SRV/srv1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=TxtRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/privateDnsZones/privateDnsZone1/TXT/txt1
