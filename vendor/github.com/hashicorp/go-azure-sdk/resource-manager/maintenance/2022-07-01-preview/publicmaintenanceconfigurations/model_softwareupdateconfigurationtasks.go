package publicmaintenanceconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationTasks struct {
	PostTasks *[]TaskProperties `json:"postTasks,omitempty"`
	PreTasks  *[]TaskProperties `json:"preTasks,omitempty"`
}
