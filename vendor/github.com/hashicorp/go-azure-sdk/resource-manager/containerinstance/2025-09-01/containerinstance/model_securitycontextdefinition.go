package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityContextDefinition struct {
	AllowPrivilegeEscalation *bool                                  `json:"allowPrivilegeEscalation,omitempty"`
	Capabilities             *SecurityContextCapabilitiesDefinition `json:"capabilities,omitempty"`
	Privileged               *bool                                  `json:"privileged,omitempty"`
	RunAsGroup               *int64                                 `json:"runAsGroup,omitempty"`
	RunAsUser                *int64                                 `json:"runAsUser,omitempty"`
	SeccompProfile           *string                                `json:"seccompProfile,omitempty"`
}
