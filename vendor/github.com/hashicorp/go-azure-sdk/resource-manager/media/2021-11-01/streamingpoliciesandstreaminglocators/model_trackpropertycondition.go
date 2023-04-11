package streamingpoliciesandstreaminglocators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrackPropertyCondition struct {
	Operation TrackPropertyCompareOperation `json:"operation"`
	Property  TrackPropertyType             `json:"property"`
	Value     *string                       `json:"value,omitempty"`
}
