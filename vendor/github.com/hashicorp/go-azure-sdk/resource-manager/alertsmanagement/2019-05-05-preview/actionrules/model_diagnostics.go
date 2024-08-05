package actionrules

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ActionRuleProperties = Diagnostics{}

type Diagnostics struct {

	// Fields inherited from ActionRuleProperties
	Conditions     *Conditions       `json:"conditions,omitempty"`
	CreatedAt      *string           `json:"createdAt,omitempty"`
	CreatedBy      *string           `json:"createdBy,omitempty"`
	Description    *string           `json:"description,omitempty"`
	LastModifiedAt *string           `json:"lastModifiedAt,omitempty"`
	LastModifiedBy *string           `json:"lastModifiedBy,omitempty"`
	Scope          *Scope            `json:"scope,omitempty"`
	Status         *ActionRuleStatus `json:"status,omitempty"`
}

func (o *Diagnostics) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *Diagnostics) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *Diagnostics) GetLastModifiedAtAsTime() (*time.Time, error) {
	if o.LastModifiedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *Diagnostics) SetLastModifiedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedAt = &formatted
}

var _ json.Marshaler = Diagnostics{}

func (s Diagnostics) MarshalJSON() ([]byte, error) {
	type wrapper Diagnostics
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Diagnostics: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Diagnostics: %+v", err)
	}
	decoded["type"] = "Diagnostics"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Diagnostics: %+v", err)
	}

	return encoded, nil
}
