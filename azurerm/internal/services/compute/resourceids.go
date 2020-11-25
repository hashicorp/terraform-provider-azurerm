package compute

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AvailabilitySet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/availabilitySets/set1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DedicatedHostGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/hostGroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DedicatedHost -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/hostGroups/hostGroup1/hosts/host1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DiskEncryptionSet -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/diskEncryptionSets/set1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Image -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/images/image1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedDisk -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/disks/disk1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ProximityPlacementGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/proximityPlacementGroups/group1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SharedImage -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/galleries/gallery1/images/image1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SharedImageGallery -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/galleries/gallery1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SharedImageVersion -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/galleries/gallery1/images/image1/versions/version1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualMachine -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/machine1
