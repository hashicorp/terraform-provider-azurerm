package dns

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DnsZone -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=CnameRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/CNAME/name1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=MxRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/MX/mx1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NsRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/NS/ns1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=PtrRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/PTR/ptr1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SrvRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/SRV/srv1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=TxtRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/TXT/txt1
