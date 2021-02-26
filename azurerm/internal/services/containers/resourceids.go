package containers

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Cluster -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=NodePool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ContainerGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerInstance/containerGroups/containerGroup1
