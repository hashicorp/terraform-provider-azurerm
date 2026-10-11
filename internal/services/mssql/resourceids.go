// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabaseExtendedAuditingPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/databases/database1/extendedAuditingSettings/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=DatabaseVulnerabilityAssessmentRuleBaseline -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/databases/database1/vulnerabilityAssessments/default/rules/rule1/baselines/baseline1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EncryptionProtector -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/encryptionProtector/current
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerExtendedAuditingPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/extendedAuditingSettings/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerMicrosoftSupportAuditingPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/devOpsAuditingSettings/default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerSecurityAlertPolicy -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/securityAlertPolicies/Default
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ServerVulnerabilityAssessment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/vulnerabilityAssessments/default

// @tombuildsstuff: this resource id going to need a state migration prior to migrating to `hashicorp/go-azure-sdk`
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=RecoverableDatabase -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/recoverabledatabases/database1
