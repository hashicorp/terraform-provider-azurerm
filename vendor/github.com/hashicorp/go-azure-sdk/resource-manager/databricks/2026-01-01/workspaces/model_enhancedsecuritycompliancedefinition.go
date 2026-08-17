package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnhancedSecurityComplianceDefinition struct {
	AutomaticClusterUpdate     *AutomaticClusterUpdateDefinition     `json:"automaticClusterUpdate,omitempty"`
	ComplianceSecurityProfile  *ComplianceSecurityProfileDefinition  `json:"complianceSecurityProfile,omitempty"`
	EnhancedSecurityMonitoring *EnhancedSecurityMonitoringDefinition `json:"enhancedSecurityMonitoring,omitempty"`
}
