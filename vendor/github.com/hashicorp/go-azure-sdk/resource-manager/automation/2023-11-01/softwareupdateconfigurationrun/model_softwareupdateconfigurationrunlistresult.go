package softwareupdateconfigurationrun

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationRunListResult struct {
	NextLink *string                           `json:"nextLink,omitempty"`
	Value    *[]SoftwareUpdateConfigurationRun `json:"value,omitempty"`
}
