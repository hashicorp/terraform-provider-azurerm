package machinelearning

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ComputeCluster -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=InferenceCluster -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=KubernetesCluster -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Workspace -id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1
