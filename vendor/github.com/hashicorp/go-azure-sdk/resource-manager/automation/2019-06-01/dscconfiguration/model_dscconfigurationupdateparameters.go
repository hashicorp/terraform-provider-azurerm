package dscconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscConfigurationUpdateParameters struct {
	Name       *string                                   `json:"name,omitempty"`
	Properties *DscConfigurationCreateOrUpdateProperties `json:"properties,omitempty"`
	Tags       *map[string]string                        `json:"tags,omitempty"`
}
