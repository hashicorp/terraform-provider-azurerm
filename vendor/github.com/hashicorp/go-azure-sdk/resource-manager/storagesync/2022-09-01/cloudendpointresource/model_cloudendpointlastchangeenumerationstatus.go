package cloudendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudEndpointLastChangeEnumerationStatus struct {
	CompletedTimestamp        *string `json:"completedTimestamp,omitempty"`
	NamespaceDirectoriesCount *int64  `json:"namespaceDirectoriesCount,omitempty"`
	NamespaceFilesCount       *int64  `json:"namespaceFilesCount,omitempty"`
	NamespaceSizeBytes        *int64  `json:"namespaceSizeBytes,omitempty"`
	NextRunTimestamp          *string `json:"nextRunTimestamp,omitempty"`
	StartedTimestamp          *string `json:"startedTimestamp,omitempty"`
}

func (o *CloudEndpointLastChangeEnumerationStatus) GetCompletedTimestampAsTime() (*time.Time, error) {
	if o.CompletedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CompletedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudEndpointLastChangeEnumerationStatus) SetCompletedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CompletedTimestamp = &formatted
}

func (o *CloudEndpointLastChangeEnumerationStatus) GetNextRunTimestampAsTime() (*time.Time, error) {
	if o.NextRunTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextRunTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudEndpointLastChangeEnumerationStatus) SetNextRunTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextRunTimestamp = &formatted
}

func (o *CloudEndpointLastChangeEnumerationStatus) GetStartedTimestampAsTime() (*time.Time, error) {
	if o.StartedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudEndpointLastChangeEnumerationStatus) SetStartedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartedTimestamp = &formatted
}
