package liveoutputs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveOutputProperties struct {
	ArchiveWindowLength string                   `json:"archiveWindowLength"`
	AssetName           string                   `json:"assetName"`
	Created             *string                  `json:"created,omitempty"`
	Description         *string                  `json:"description,omitempty"`
	Hls                 *Hls                     `json:"hls,omitempty"`
	LastModified        *string                  `json:"lastModified,omitempty"`
	ManifestName        *string                  `json:"manifestName,omitempty"`
	OutputSnapTime      *int64                   `json:"outputSnapTime,omitempty"`
	ProvisioningState   *string                  `json:"provisioningState,omitempty"`
	ResourceState       *LiveOutputResourceState `json:"resourceState,omitempty"`
}

func (o *LiveOutputProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *LiveOutputProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}

func (o *LiveOutputProperties) GetLastModifiedAsTime() (*time.Time, error) {
	if o.LastModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o *LiveOutputProperties) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}
