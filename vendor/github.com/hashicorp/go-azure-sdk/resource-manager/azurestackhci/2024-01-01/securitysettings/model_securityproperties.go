package securitysettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityProperties struct {
	ProvisioningState               *ProvisioningState        `json:"provisioningState,omitempty"`
	SecuredCoreComplianceAssignment *ComplianceAssignmentType `json:"securedCoreComplianceAssignment,omitempty"`
	SecurityComplianceStatus        *SecurityComplianceStatus `json:"securityComplianceStatus,omitempty"`
}
