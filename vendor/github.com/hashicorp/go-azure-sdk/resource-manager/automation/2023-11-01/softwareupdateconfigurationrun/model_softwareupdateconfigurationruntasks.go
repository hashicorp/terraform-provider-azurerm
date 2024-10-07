package softwareupdateconfigurationrun

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationRunTasks struct {
	PostTask *SoftwareUpdateConfigurationRunTaskProperties `json:"postTask,omitempty"`
	PreTask  *SoftwareUpdateConfigurationRunTaskProperties `json:"preTask,omitempty"`
}
