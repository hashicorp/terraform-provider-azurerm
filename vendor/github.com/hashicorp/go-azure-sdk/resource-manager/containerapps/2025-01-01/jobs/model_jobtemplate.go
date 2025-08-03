package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobTemplate struct {
	Containers     *[]Container     `json:"containers,omitempty"`
	InitContainers *[]BaseContainer `json:"initContainers,omitempty"`
	Volumes        *[]Volume        `json:"volumes,omitempty"`
}
