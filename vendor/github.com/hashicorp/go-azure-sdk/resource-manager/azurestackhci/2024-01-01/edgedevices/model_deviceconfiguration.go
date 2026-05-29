package edgedevices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeviceConfiguration struct {
	DeviceMetadata *string      `json:"deviceMetadata,omitempty"`
	NicDetails     *[]NicDetail `json:"nicDetails,omitempty"`
}
