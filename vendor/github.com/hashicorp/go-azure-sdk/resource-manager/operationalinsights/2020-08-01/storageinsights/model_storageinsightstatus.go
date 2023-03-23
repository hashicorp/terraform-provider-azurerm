package storageinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageInsightStatus struct {
	Description *string             `json:"description,omitempty"`
	State       StorageInsightState `json:"state"`
}
