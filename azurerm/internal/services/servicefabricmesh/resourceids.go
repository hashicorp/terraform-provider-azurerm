package servicefabricmesh

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Application -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceFabricMesh/applications/application1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Network -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceFabricMesh/networks/network1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Secret -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceFabricMesh/secrets/secret1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SecretValue -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ServiceFabricMesh/secrets/secret1/values/value1
