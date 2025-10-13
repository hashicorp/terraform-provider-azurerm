package jobsteps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStep struct {
	Id         *string            `json:"id,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties *JobStepProperties `json:"properties,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
