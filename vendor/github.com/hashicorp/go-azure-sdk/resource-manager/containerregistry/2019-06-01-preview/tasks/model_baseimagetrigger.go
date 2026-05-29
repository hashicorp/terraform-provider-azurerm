package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaseImageTrigger struct {
	BaseImageTriggerType     BaseImageTriggerType      `json:"baseImageTriggerType"`
	Name                     string                    `json:"name"`
	Status                   *TriggerStatus            `json:"status,omitempty"`
	UpdateTriggerEndpoint    *string                   `json:"updateTriggerEndpoint,omitempty"`
	UpdateTriggerPayloadType *UpdateTriggerPayloadType `json:"updateTriggerPayloadType,omitempty"`
}
