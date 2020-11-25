package dns

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DnsZone -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SrvRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/SRV/srv1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=TxtRecord -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/dnszones/zone1/TXT/txt1
