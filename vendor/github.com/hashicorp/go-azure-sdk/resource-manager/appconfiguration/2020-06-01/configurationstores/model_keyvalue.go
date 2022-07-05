package configurationstores

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyValue struct {
	ContentType  *string            `json:"contentType,omitempty"`
	ETag         *string            `json:"eTag,omitempty"`
	Key          *string            `json:"key,omitempty"`
	Label        *string            `json:"label,omitempty"`
	LastModified *string            `json:"lastModified,omitempty"`
	Locked       *bool              `json:"locked,omitempty"`
	Tags         *map[string]string `json:"tags,omitempty"`
	Value        *string            `json:"value,omitempty"`
}

func (o *KeyValue) GetLastModifiedAsTime() (*time.Time, error) {
	if o.LastModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o *KeyValue) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}
