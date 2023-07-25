// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Automation -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Security/automations/testAutomation1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Contact -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/securityContacts/contact1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Pricing -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/pricings/pricing1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Setting -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/settings/setting1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Workspace -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/workspaceSettings/workspace1

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AssessmentMetadata -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/assessmentMetadata/metadata1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=AutoProvisioningSetting -id=/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/autoProvisioningSettings/default -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IotSecuritySolution -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Security/iotSecuritySolutions/solution1 -rewrite=true

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VulnerabilityAssessmentVm -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm-name1/providers/Microsoft.Security/serverVulnerabilityAssessments/default1
