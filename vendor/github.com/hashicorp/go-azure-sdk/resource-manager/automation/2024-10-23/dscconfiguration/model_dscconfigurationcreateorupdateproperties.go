package dscconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscConfigurationCreateOrUpdateProperties struct {
	Description *string                               `json:"description,omitempty"`
	LogProgress *bool                                 `json:"logProgress,omitempty"`
	LogVerbose  *bool                                 `json:"logVerbose,omitempty"`
	Parameters  *map[string]DscConfigurationParameter `json:"parameters,omitempty"`
	Source      ContentSource                         `json:"source"`
}
