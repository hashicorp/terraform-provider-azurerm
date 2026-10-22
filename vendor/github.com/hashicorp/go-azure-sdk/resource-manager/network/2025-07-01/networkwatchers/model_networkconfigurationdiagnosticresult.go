package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkConfigurationDiagnosticResult struct {
	NetworkSecurityGroupResult *NetworkSecurityGroupResult            `json:"networkSecurityGroupResult,omitempty"`
	Profile                    *NetworkConfigurationDiagnosticProfile `json:"profile,omitempty"`
}
