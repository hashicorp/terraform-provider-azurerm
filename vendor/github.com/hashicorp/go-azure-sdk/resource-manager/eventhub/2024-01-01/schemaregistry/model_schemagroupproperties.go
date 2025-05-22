package schemaregistry

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaGroupProperties struct {
	CreatedAtUtc        *string              `json:"createdAtUtc,omitempty"`
	ETag                *string              `json:"eTag,omitempty"`
	GroupProperties     *map[string]string   `json:"groupProperties,omitempty"`
	SchemaCompatibility *SchemaCompatibility `json:"schemaCompatibility,omitempty"`
	SchemaType          *SchemaType          `json:"schemaType,omitempty"`
	UpdatedAtUtc        *string              `json:"updatedAtUtc,omitempty"`
}

func (o *SchemaGroupProperties) GetCreatedAtUtcAsTime() (*time.Time, error) {
	if o.CreatedAtUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAtUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *SchemaGroupProperties) SetCreatedAtUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAtUtc = &formatted
}

func (o *SchemaGroupProperties) GetUpdatedAtUtcAsTime() (*time.Time, error) {
	if o.UpdatedAtUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedAtUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *SchemaGroupProperties) SetUpdatedAtUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAtUtc = &formatted
}
