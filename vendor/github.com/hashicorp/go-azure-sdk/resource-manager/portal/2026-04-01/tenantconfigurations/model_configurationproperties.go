package tenantconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationProperties struct {
	EnforcePrivateMarkdownStorage *bool                      `json:"enforcePrivateMarkdownStorage,omitempty"`
	ProvisioningState             *ResourceProvisioningState `json:"provisioningState,omitempty"`
}
