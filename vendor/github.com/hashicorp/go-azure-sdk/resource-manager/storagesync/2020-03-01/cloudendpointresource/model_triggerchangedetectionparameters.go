package cloudendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggerChangeDetectionParameters struct {
	ChangeDetectionMode *ChangeDetectionMode `json:"changeDetectionMode,omitempty"`
	DirectoryPath       *string              `json:"directoryPath,omitempty"`
	Paths               *[]string            `json:"paths,omitempty"`
}
