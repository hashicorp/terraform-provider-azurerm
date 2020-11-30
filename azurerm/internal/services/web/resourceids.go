package web

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedCertificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/customhost.contoso.com
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SlotVirtualNetworkSwiftConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1/networkconfig/virtualNetwork
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkSwiftConnection -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/networkconfig/virtualNetwork
