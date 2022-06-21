package policy

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ResourceGroupAssignment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Authorization/policyAssignments/assignment1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SubscriptionAssignment -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/policyAssignments/assignment1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualMachineConfigurationAssignment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/assignment1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ResourceGroupPolicyRemediation -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.PolicyInsights/remediations/remediation1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SubscriptionPolicyRemediation -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.PolicyInsights/remediations/remediation1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ResourceGroupPolicyExemption -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Authorization/policyExemptions/exemption1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SubscriptionPolicyExemption -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Authorization/policyExemptions/exemption1
