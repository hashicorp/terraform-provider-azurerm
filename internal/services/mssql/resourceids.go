// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabaseExtendedAuditingPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/databases/database1/extendedAuditingSettings/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabaseVulnerabilityAssessmentRuleBaseline -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/databases/database1/vulnerabilityAssessments/default/rules/rule1/baselines/baseline1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ElasticPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/elasticPools/pool1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EncryptionProtector -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/encryptionProtector/current
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FailoverGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/failoverGroups/failoverGroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=FirewallRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/firewallRules/rule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=JobAgent -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/jobAgents/jobagent1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=JobCredential -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/jobAgents/jobagent1/credentials/credential1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=OutboundFirewallRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/outboundFirewallRules/fqdn1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Server -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerDNSAlias -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/dnsAliases/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerExtendedAuditingPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/extendedAuditingSettings/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerMicrosoftSupportAuditingPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/devOpsAuditingSettings/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerSecurityAlertPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/securityAlertPolicies/Default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerVulnerabilityAssessment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/vulnerabilityAssessments/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SqlVirtualMachine -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/virtualMachine1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=VirtualNetworkRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/virtualNetworkRules/virtualNetworkRule1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ManagedInstancesSecurityAlertPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/managedInstances/instance1/securityAlertPolicies/Default

// @tombuildsstuff: this resource id going to need a state migration prior to migrating to `hashicorp/go-azure-sdk`
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RecoverableDatabase -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/recoverabledatabases/database1
