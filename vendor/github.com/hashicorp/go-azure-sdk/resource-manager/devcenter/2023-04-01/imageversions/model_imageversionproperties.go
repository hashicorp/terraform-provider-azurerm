package imageversions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageVersionProperties struct {
	ExcludeFromLatest   *bool              `json:"excludeFromLatest,omitempty"`
	Name                *string            `json:"name,omitempty"`
	OsDiskImageSizeInGb *int64             `json:"osDiskImageSizeInGb,omitempty"`
	ProvisioningState   *ProvisioningState `json:"provisioningState,omitempty"`
	PublishedDate       *string            `json:"publishedDate,omitempty"`
}

func (o *ImageVersionProperties) GetPublishedDateAsTime() (*time.Time, error) {
	if o.PublishedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PublishedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageVersionProperties) SetPublishedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PublishedDate = &formatted
}
