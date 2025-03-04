package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerExecutionStatus struct {
	AdditionalInformation *string `json:"additionalInformation,omitempty"`
	Code                  *int64  `json:"code,omitempty"`
	Name                  *string `json:"name,omitempty"`
	Status                *string `json:"status,omitempty"`
}
