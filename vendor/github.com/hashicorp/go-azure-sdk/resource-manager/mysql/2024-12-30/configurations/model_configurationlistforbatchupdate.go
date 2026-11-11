package configurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationListForBatchUpdate struct {
	ResetAllToDefault *ResetAllToDefault             `json:"resetAllToDefault,omitempty"`
	Value             *[]ConfigurationForBatchUpdate `json:"value,omitempty"`
}
