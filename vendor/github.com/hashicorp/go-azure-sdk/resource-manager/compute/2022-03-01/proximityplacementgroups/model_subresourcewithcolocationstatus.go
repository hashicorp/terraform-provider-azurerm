package proximityplacementgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubResourceWithColocationStatus struct {
	ColocationStatus *InstanceViewStatus `json:"colocationStatus,omitempty"`
	Id               *string             `json:"id,omitempty"`
}
