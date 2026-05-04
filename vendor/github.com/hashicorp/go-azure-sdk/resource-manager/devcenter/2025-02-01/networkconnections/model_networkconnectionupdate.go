package networkconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkConnectionUpdate struct {
	Location   *string                            `json:"location,omitempty"`
	Properties *NetworkConnectionUpdateProperties `json:"properties,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
}
