package compute

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AvailabilitySet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/availabilitySets/set1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DedicatedHostGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/hostGroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DedicatedHost -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/hostGroup1/hosts/host1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DiskEncryptionSet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/diskEncryptionSets/set1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Image -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/images/image1
