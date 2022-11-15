package devtestlabs

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Schedule -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/schedules/schedule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DevTestLab -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DevTestLabPolicy -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/policySets/policyset1/policies/policy1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DevTestVirtualMachine -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/virtualMachines/vm1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DevTestVirtualNetwork -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/virtualNetworks/vnet1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DevTestLabSchedule -rewrite=true -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/schedules/schedule1
