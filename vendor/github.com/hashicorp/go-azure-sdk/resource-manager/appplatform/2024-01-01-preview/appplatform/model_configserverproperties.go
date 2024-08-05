package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigServerProperties struct {
	ConfigServer      *ConfigServerSettings     `json:"configServer,omitempty"`
	EnabledState      *ConfigServerEnabledState `json:"enabledState,omitempty"`
	Error             *Error                    `json:"error,omitempty"`
	ProvisioningState *ConfigServerState        `json:"provisioningState,omitempty"`
}
