package configurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationProperties struct {
	AllowedValues          *string                `json:"allowedValues,omitempty"`
	DataType               *ConfigurationDataType `json:"dataType,omitempty"`
	DefaultValue           *string                `json:"defaultValue,omitempty"`
	Description            *string                `json:"description,omitempty"`
	DocumentationLink      *string                `json:"documentationLink,omitempty"`
	IsConfigPendingRestart *bool                  `json:"isConfigPendingRestart,omitempty"`
	IsDynamicConfig        *bool                  `json:"isDynamicConfig,omitempty"`
	IsReadOnly             *bool                  `json:"isReadOnly,omitempty"`
	Source                 *string                `json:"source,omitempty"`
	Unit                   *string                `json:"unit,omitempty"`
	Value                  *string                `json:"value,omitempty"`
}
