package configurationprofilehciassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationProfileAssignmentProperties struct {
	ConfigurationProfile *string `json:"configurationProfile,omitempty"`
	Status               *string `json:"status,omitempty"`
	TargetId             *string `json:"targetId,omitempty"`
}
