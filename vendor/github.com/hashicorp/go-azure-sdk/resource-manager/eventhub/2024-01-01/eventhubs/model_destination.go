package eventhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Destination struct {
	Identity   *CaptureIdentity       `json:"identity,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties *DestinationProperties `json:"properties,omitempty"`
}
