// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package securitycenter

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Contact -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/securityContacts/contact1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Workspace -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/workspaceSettings/workspace1

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IotSecuritySolution -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Security/iotSecuritySolutions/solution1 -rewrite=true

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VulnerabilityAssessmentVm -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm-name1/providers/Microsoft.Security/serverVulnerabilityAssessments/default1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VulnerabilityAssessmentsSetting -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/serverVulnerabilityAssessmentsSettings/AzureServersSetting
