package bigdatapools

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BigDataPoolResourceProperties struct {
	AutoPause                   *AutoPauseProperties       `json:"autoPause,omitempty"`
	AutoScale                   *AutoScaleProperties       `json:"autoScale,omitempty"`
	CacheSize                   *int64                     `json:"cacheSize,omitempty"`
	CreationDate                *string                    `json:"creationDate,omitempty"`
	CustomLibraries             *[]LibraryInfo             `json:"customLibraries,omitempty"`
	DefaultSparkLogFolder       *string                    `json:"defaultSparkLogFolder,omitempty"`
	DynamicExecutorAllocation   *DynamicExecutorAllocation `json:"dynamicExecutorAllocation,omitempty"`
	IsAutotuneEnabled           *bool                      `json:"isAutotuneEnabled,omitempty"`
	IsComputeIsolationEnabled   *bool                      `json:"isComputeIsolationEnabled,omitempty"`
	LastSucceededTimestamp      *string                    `json:"lastSucceededTimestamp,omitempty"`
	LibraryRequirements         *LibraryRequirements       `json:"libraryRequirements,omitempty"`
	NodeCount                   *int64                     `json:"nodeCount,omitempty"`
	NodeSize                    *NodeSize                  `json:"nodeSize,omitempty"`
	NodeSizeFamily              *NodeSizeFamily            `json:"nodeSizeFamily,omitempty"`
	ProvisioningState           *string                    `json:"provisioningState,omitempty"`
	SessionLevelPackagesEnabled *bool                      `json:"sessionLevelPackagesEnabled,omitempty"`
	SparkConfigProperties       *SparkConfigProperties     `json:"sparkConfigProperties,omitempty"`
	SparkEventsFolder           *string                    `json:"sparkEventsFolder,omitempty"`
	SparkVersion                *string                    `json:"sparkVersion,omitempty"`
}

func (o *BigDataPoolResourceProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *BigDataPoolResourceProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}

func (o *BigDataPoolResourceProperties) GetLastSucceededTimestampAsTime() (*time.Time, error) {
	if o.LastSucceededTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSucceededTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *BigDataPoolResourceProperties) SetLastSucceededTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSucceededTimestamp = &formatted
}
