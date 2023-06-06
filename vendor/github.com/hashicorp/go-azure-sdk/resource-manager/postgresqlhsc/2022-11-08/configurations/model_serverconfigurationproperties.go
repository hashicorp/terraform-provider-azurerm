package configurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerConfigurationProperties struct {
	AllowedValues     *string                `json:"allowedValues,omitempty"`
	DataType          *ConfigurationDataType `json:"dataType,omitempty"`
	DefaultValue      *string                `json:"defaultValue,omitempty"`
	Description       *string                `json:"description,omitempty"`
	ProvisioningState *ProvisioningState     `json:"provisioningState,omitempty"`
	RequiresRestart   *bool                  `json:"requiresRestart,omitempty"`
	Source            *string                `json:"source,omitempty"`
	Value             string                 `json:"value"`
}
