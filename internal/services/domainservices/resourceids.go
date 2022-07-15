package domainservices

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DomainService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AAD/domainServices/DomainService1/initialReplicaSetId/replicaSetID
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DomainServiceReplicaSet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AAD/domainServices/DomainService1/replicaSets/replicaSetID
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DomainServiceTrust -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AAD/domainServices/DomainService1/trusts/trust1
