package configurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationProperties struct {
	AllowedValues          *string                 `json:"allowedValues,omitempty"`
	CurrentValue           *string                 `json:"currentValue,omitempty"`
	DataType               *string                 `json:"dataType,omitempty"`
	DefaultValue           *string                 `json:"defaultValue,omitempty"`
	Description            *string                 `json:"description,omitempty"`
	DocumentationLink      *string                 `json:"documentationLink,omitempty"`
	IsConfigPendingRestart *IsConfigPendingRestart `json:"isConfigPendingRestart,omitempty"`
	IsDynamicConfig        *IsDynamicConfig        `json:"isDynamicConfig,omitempty"`
	IsReadOnly             *IsReadOnly             `json:"isReadOnly,omitempty"`
	Source                 *ConfigurationSource    `json:"source,omitempty"`
	Value                  *string                 `json:"value,omitempty"`
}
