package raipolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SafetyProviderConfig struct {
	Blocking           *bool                   `json:"blocking,omitempty"`
	SafetyProviderName *string                 `json:"safetyProviderName,omitempty"`
	Source             *RaiPolicyContentSource `json:"source,omitempty"`
}
