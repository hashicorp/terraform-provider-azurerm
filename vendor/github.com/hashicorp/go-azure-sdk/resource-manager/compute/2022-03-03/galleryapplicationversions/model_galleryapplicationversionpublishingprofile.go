package galleryapplicationversions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryApplicationVersionPublishingProfile struct {
	AdvancedSettings        *map[string]string                `json:"advancedSettings,omitempty"`
	CustomActions           *[]GalleryApplicationCustomAction `json:"customActions,omitempty"`
	EnableHealthCheck       *bool                             `json:"enableHealthCheck,omitempty"`
	EndOfLifeDate           *string                           `json:"endOfLifeDate,omitempty"`
	ExcludeFromLatest       *bool                             `json:"excludeFromLatest,omitempty"`
	ManageActions           *UserArtifactManage               `json:"manageActions,omitempty"`
	PublishedDate           *string                           `json:"publishedDate,omitempty"`
	ReplicaCount            *int64                            `json:"replicaCount,omitempty"`
	ReplicationMode         *ReplicationMode                  `json:"replicationMode,omitempty"`
	Settings                *UserArtifactSettings             `json:"settings,omitempty"`
	Source                  UserArtifactSource                `json:"source"`
	StorageAccountType      *StorageAccountType               `json:"storageAccountType,omitempty"`
	TargetExtendedLocations *[]GalleryTargetExtendedLocation  `json:"targetExtendedLocations,omitempty"`
	TargetRegions           *[]TargetRegion                   `json:"targetRegions,omitempty"`
}

func (o *GalleryApplicationVersionPublishingProfile) GetEndOfLifeDateAsTime() (*time.Time, error) {
	if o.EndOfLifeDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndOfLifeDate, "2006-01-02T15:04:05Z07:00")
}

func (o *GalleryApplicationVersionPublishingProfile) SetEndOfLifeDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndOfLifeDate = &formatted
}

func (o *GalleryApplicationVersionPublishingProfile) GetPublishedDateAsTime() (*time.Time, error) {
	if o.PublishedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PublishedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *GalleryApplicationVersionPublishingProfile) SetPublishedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PublishedDate = &formatted
}
