package sourcecontrol

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceControlProperties struct {
	AutoSync         *bool       `json:"autoSync,omitempty"`
	Branch           *string     `json:"branch,omitempty"`
	CreationTime     *string     `json:"creationTime,omitempty"`
	Description      *string     `json:"description,omitempty"`
	FolderPath       *string     `json:"folderPath,omitempty"`
	LastModifiedTime *string     `json:"lastModifiedTime,omitempty"`
	PublishRunbook   *bool       `json:"publishRunbook,omitempty"`
	RepoUrl          *string     `json:"repoUrl,omitempty"`
	SourceType       *SourceType `json:"sourceType,omitempty"`
}

func (o *SourceControlProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SourceControlProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *SourceControlProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SourceControlProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
