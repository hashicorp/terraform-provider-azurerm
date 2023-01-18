package liveevents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveEventInputTrackSelection struct {
	Operation *string `json:"operation,omitempty"`
	Property  *string `json:"property,omitempty"`
	Value     *string `json:"value,omitempty"`
}
