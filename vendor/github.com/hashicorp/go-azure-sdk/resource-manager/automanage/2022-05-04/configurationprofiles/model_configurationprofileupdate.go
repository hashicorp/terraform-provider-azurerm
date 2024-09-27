package configurationprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationProfileUpdate struct {
	Properties *ConfigurationProfileProperties `json:"properties,omitempty"`
	Tags       *map[string]string              `json:"tags,omitempty"`
}
