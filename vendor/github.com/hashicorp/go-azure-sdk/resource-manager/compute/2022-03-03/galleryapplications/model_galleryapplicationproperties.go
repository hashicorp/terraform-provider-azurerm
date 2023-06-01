package galleryapplications

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryApplicationProperties struct {
	CustomActions       *[]GalleryApplicationCustomAction `json:"customActions,omitempty"`
	Description         *string                           `json:"description,omitempty"`
	EndOfLifeDate       *string                           `json:"endOfLifeDate,omitempty"`
	Eula                *string                           `json:"eula,omitempty"`
	PrivacyStatementUri *string                           `json:"privacyStatementUri,omitempty"`
	ReleaseNoteUri      *string                           `json:"releaseNoteUri,omitempty"`
	SupportedOSType     OperatingSystemTypes              `json:"supportedOSType"`
}

func (o *GalleryApplicationProperties) GetEndOfLifeDateAsTime() (*time.Time, error) {
	if o.EndOfLifeDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndOfLifeDate, "2006-01-02T15:04:05Z07:00")
}

func (o *GalleryApplicationProperties) SetEndOfLifeDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndOfLifeDate = &formatted
}
