package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficWeight struct {
	Label          *string `json:"label,omitempty"`
	LatestRevision *bool   `json:"latestRevision,omitempty"`
	RevisionName   *string `json:"revisionName,omitempty"`
	Weight         *int64  `json:"weight,omitempty"`
}
