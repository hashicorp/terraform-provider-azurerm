package integrationaccountassemblies

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssemblyProperties struct {
	AssemblyCulture        *string      `json:"assemblyCulture,omitempty"`
	AssemblyName           string       `json:"assemblyName"`
	AssemblyPublicKeyToken *string      `json:"assemblyPublicKeyToken,omitempty"`
	AssemblyVersion        *string      `json:"assemblyVersion,omitempty"`
	ChangedTime            *string      `json:"changedTime,omitempty"`
	Content                *interface{} `json:"content,omitempty"`
	ContentLink            *ContentLink `json:"contentLink,omitempty"`
	ContentType            *string      `json:"contentType,omitempty"`
	CreatedTime            *string      `json:"createdTime,omitempty"`
	Metadata               *interface{} `json:"metadata,omitempty"`
}

func (o *AssemblyProperties) GetChangedTimeAsTime() (*time.Time, error) {
	if o.ChangedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ChangedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AssemblyProperties) SetChangedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ChangedTime = &formatted
}

func (o *AssemblyProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AssemblyProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}
