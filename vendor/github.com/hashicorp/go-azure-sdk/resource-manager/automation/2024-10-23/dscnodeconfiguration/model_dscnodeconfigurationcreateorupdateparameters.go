package dscnodeconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscNodeConfigurationCreateOrUpdateParameters struct {
	Name       *string                                                 `json:"name,omitempty"`
	Properties *DscNodeConfigurationCreateOrUpdateParametersProperties `json:"properties,omitempty"`
	Tags       *map[string]string                                      `json:"tags,omitempty"`
}
