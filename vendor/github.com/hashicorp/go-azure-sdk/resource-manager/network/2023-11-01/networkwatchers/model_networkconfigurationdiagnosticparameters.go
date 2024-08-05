package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkConfigurationDiagnosticParameters struct {
	Profiles         []NetworkConfigurationDiagnosticProfile `json:"profiles"`
	TargetResourceId string                                  `json:"targetResourceId"`
	VerbosityLevel   *VerbosityLevel                         `json:"verbosityLevel,omitempty"`
}
