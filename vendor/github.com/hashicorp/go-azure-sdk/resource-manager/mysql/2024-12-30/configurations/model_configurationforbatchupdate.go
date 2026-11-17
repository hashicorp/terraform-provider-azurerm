package configurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationForBatchUpdate struct {
	Name       *string                                `json:"name,omitempty"`
	Properties *ConfigurationForBatchUpdateProperties `json:"properties,omitempty"`
}
